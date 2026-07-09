package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/db"
	"crowdfunding-backend/models"
)

const maxAssetUploadSize = 10 * 1024 * 1024

type CampaignResponse struct {
	ID           uint64 `json:"id"`
	Owner        string `json:"owner"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Goal         string `json:"goal"`
	Deadline     string `json:"deadline"`
	AmountRaised string `json:"amountRaised"`
	Withdrawn    bool   `json:"withdrawn"`
	Status       string `json:"status"`
}

func campaignStatus(campaign contract.Campaign) string {
	if campaign.AmountRaised.Cmp(campaign.Goal) >= 0 {
		return "Successful"
	}
	if time.Now().Unix() >= campaign.Deadline.Int64() {
		return "Failed"
	}
	return "Active"
}

func toCampaignResponse(id uint64, campaign contract.Campaign) CampaignResponse {
	return CampaignResponse{
		ID:           id,
		Owner:        campaign.Owner.Hex(),
		Title:        campaign.Title,
		Description:  campaign.Description,
		Goal:         campaign.Goal.String(),
		Deadline:     campaign.Deadline.String(),
		AmountRaised: campaign.AmountRaised.String(),
		Withdrawn:    campaign.Withdrawn,
		Status:       campaignStatus(campaign),
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on existing environment variables")
	}

	rpcURL := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	jwtSecret := os.Getenv("JWT_SECRET")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresPort := os.Getenv("POSTGRES_PORT")
	auth0Domain := os.Getenv("AUTH0_APP_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	if rpcURL == "" || contractAddress == "" || jwtSecret == "" ||
		postgresUser == "" || postgresPassword == "" || postgresDB == "" || postgresPort == "" ||
		auth0Domain == "" || auth0Audience == "" {
		log.Fatal("RPC_URL, CONTRACT_ADDRESS, JWT_SECRET, POSTGRES_*, and AUTH0_* variables must be set")
	}

	auth0KeyFunc, err := newAuth0KeyFunc(auth0Domain)
	if err != nil {
		log.Fatalf("failed to load auth0 jwks: %v", err)
	}

	r2AccountID := os.Getenv("R2_ACCOUNT_ID")
	r2AccessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2Bucket := os.Getenv("R2_BUCKET")
	r2PublicURL := os.Getenv("R2_PUBLIC_URL")

	r2Client, err := newR2Client(context.Background(), r2AccountID, r2AccessKeyID, r2SecretAccessKey)
	if err != nil {
		log.Fatalf("failed to init r2 client: %v", err)
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB,
	)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to rpc: %v", err)
	}

	crowdFunding, err := contract.NewCrowdFunding(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatalf("failed to bind contract: %v", err)
	}

	gormDB, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Migrate(gormDB); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	startContributionIndexer(gormDB, crowdFunding, client)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	secret := []byte(jwtSecret)
	nonces := newNonceStore()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	api.GET("/auth/nonce", func(c *gin.Context) {
		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
			return
		}

		nonce, err := nonces.Issue(address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue nonce"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"nonce": nonce, "message": buildSignInMessage(nonce)})
	})

	api.POST("/auth/verify", func(c *gin.Context) {
		var req struct {
			Address   string `json:"address" binding:"required"`
			Signature string `json:"signature" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address and signature are required"})
			return
		}

		nonce, ok := nonces.PeekAndDelete(req.Address)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "nonce not found or expired, request a new one"})
			return
		}

		expectedMessage := buildSignInMessage(nonce)
		recovered, err := recoverAddress(expectedMessage, req.Signature)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		if !strings.EqualFold(recovered.Hex(), req.Address) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "signature does not match address"})
			return
		}

		token, err := generateJWT(secret, recovered.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "address": recovered.Hex()})
	})

	api.GET("/me", authMiddleware(secret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"address": c.GetString("address")})
	})

	api.GET("/me/profile", authMiddleware(secret), func(c *gin.Context) {
		profile, err := models.GetProfile(gormDB, c.GetString("address"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profile)
	})

	api.PUT("/me/profile", authMiddleware(secret), func(c *gin.Context) {
		var req struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		profile, err := models.UpsertProfile(gormDB, c.GetString("address"), req.DisplayName, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profile)
	})

	api.POST("/auth0/sync", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		sub := c.GetString("sub")
		accessToken := c.GetString("auth0Token")

		email, name, err := fetchAuth0UserInfo(auth0Domain, accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user, err := models.UpsertUserFromAuth0(gormDB, sub, email, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	api.GET("/auth0/me", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		user, err := models.GetUser(gormDB, c.GetString("sub"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	api.POST("/assets", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		if r2Bucket == "" || r2PublicURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "asset storage is not configured yet"})
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		if fileHeader.Size > maxAssetUploadSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file exceeds 10MB limit"})
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only image uploads are allowed"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
			return
		}
		defer file.Close()

		sub := c.GetString("sub")
		ext := filepath.Ext(fileHeader.Filename)
		objectKey := fmt.Sprintf("uploads/campaign/covers/%s%s", uuid.NewString(), ext)

		if err := uploadObjectToR2(c.Request.Context(), r2Client, r2Bucket, objectKey, contentType, file, fileHeader.Size); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		assetURL := fmt.Sprintf("%s/%s", strings.TrimRight(r2PublicURL, "/"), objectKey)

		asset, err := models.CreateAsset(gormDB, sub, r2Bucket, objectKey, assetURL, contentType, fileHeader.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, asset)
	})

	api.GET("/my-campaigns", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		pagination, err := parsePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sub := c.GetString("sub")

		campaigns, total, err := models.ListCampaignsByOwner(gormDB, sub, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		campaignIDs := make([]uint64, len(campaigns))
		for i, campaign := range campaigns {
			campaignIDs[i] = campaign.ID
		}

		coverURLs, err := models.GetCoverAssetsForCampaigns(gormDB, campaignIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items := make([]gin.H, len(campaigns))
		for i, campaign := range campaigns {
			items[i] = gin.H{
				"id":             campaign.ID,
				"title":          campaign.Title,
				"description":    campaign.Description,
				"targetEth":      campaign.TargetEth,
				"country":        campaign.Country,
				"fundraisingFor": campaign.FundraisingFor,
				"status":         campaign.Status,
				"coverUrl":       coverURLs[campaign.ID],
				"createdAt":      campaign.CreatedAt,
			}
		}

		c.JSON(http.StatusOK, PaginatedResponse{
			Items:  items,
			Total:  total,
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
	})

	api.POST("/my-campaigns", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		var req struct {
			Country        string   `json:"country" binding:"required"`
			Title          string   `json:"title" binding:"required"`
			Description    string   `json:"description"`
			TargetEth      string   `json:"targetEth" binding:"required"`
			FundraisingFor string   `json:"fundraisingFor" binding:"required"`
			AssetIDs       []uint64 `json:"assetIds" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "at least one image is required"})
			return
		}

		sub := c.GetString("sub")

		assets, err := models.GetAssetsByIDs(gormDB, req.AssetIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(assets) != len(req.AssetIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more assets do not exist"})
			return
		}
		for _, asset := range assets {
			if asset.UploadedBy != sub {
				c.JSON(http.StatusForbidden, gin.H{"error": "one or more assets do not belong to you"})
				return
			}
		}

		campaign, err := models.CreateCampaign(gormDB, sub, req.Country, req.Title, req.Description, req.TargetEth, req.FundraisingFor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := models.AttachAssetsToCampaign(gormDB, campaign.ID, req.AssetIDs, req.AssetIDs[0]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, campaign)
	})

	api.GET("/my-campaigns/:id", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		campaignID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		campaign, err := models.GetCampaignByID(gormDB, campaignID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if campaign == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
			return
		}
		if campaign.OwnerSub != c.GetString("sub") {
			c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this campaign"})
			return
		}

		assets, err := models.GetCampaignAssets(gormDB, campaign.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                campaign.ID,
			"ownerSub":          campaign.OwnerSub,
			"country":           campaign.Country,
			"title":             campaign.Title,
			"description":       campaign.Description,
			"targetEth":         campaign.TargetEth,
			"fundraisingFor":    campaign.FundraisingFor,
			"status":            campaign.Status,
			"walletAddress":     campaign.WalletAddress,
			"onChainCampaignId": campaign.OnChainCampaignID,
			"publishedAt":       campaign.PublishedAt,
			"createdAt":         campaign.CreatedAt,
			"assets":            assets,
		})
	})

	api.GET("/profiles/:address", func(c *gin.Context) {
		profile, err := models.GetProfile(gormDB, c.Param("address"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"address": profile.Address, "displayName": profile.DisplayName})
	})

	api.GET("/campaigns/count", func(c *gin.Context) {
		count, err := crowdFunding.CampaignCount(nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count.String()})
	})

	api.GET("/campaigns", func(c *gin.Context) {
		offset, err := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}

		limit, err := strconv.ParseUint(c.DefaultQuery("limit", "20"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}

		campaigns, err := crowdFunding.GetCampaigns(nil, new(big.Int).SetUint64(offset), new(big.Int).SetUint64(limit))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]CampaignResponse, len(campaigns))
		for i, campaign := range campaigns {
			response[i] = toCampaignResponse(offset+uint64(i), campaign)
		}

		c.JSON(http.StatusOK, response)
	})

	api.GET("/campaigns/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		campaign, err := crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toCampaignResponse(id, campaign))
	})

	api.GET("/campaigns/:id/contributors", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		summaries, err := models.GetContributorsForCampaign(gormDB, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]gin.H, len(summaries))
		for i, summary := range summaries {
			profile, err := models.GetProfile(gormDB, summary.Contributor)
			displayName := ""
			if err == nil {
				displayName = profile.DisplayName
			}

			response[i] = gin.H{
				"address":     summary.Contributor,
				"displayName": displayName,
				"amount":      summary.TotalAmount,
			}
		}

		c.JSON(http.StatusOK, response)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
