package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
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
	ID                string `json:"id"`
	OnChainCampaignID uint64 `json:"onChainCampaignId"`
	Owner             string `json:"owner"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Goal              string `json:"goal"`
	Deadline          string `json:"deadline"`
	AmountRaised      string `json:"amountRaised"`
	Withdrawn         bool   `json:"withdrawn"`
	Status            string `json:"status"`
	Country           string `json:"country"`
	Category          string `json:"category"`
	CoverURL          string `json:"coverUrl"`
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

func toCampaignResponse(maskedID string, onChainID uint64, campaign contract.Campaign, dbCampaign *models.Campaign, coverURL string) CampaignResponse {
	description := campaign.Description
	country := ""
	category := ""
	if dbCampaign != nil {
		description = dbCampaign.Description
		country = dbCampaign.Country
		category = dbCampaign.Category
	}

	return CampaignResponse{
		ID:                maskedID,
		OnChainCampaignID: onChainID,
		Owner:             campaign.Owner.Hex(),
		Title:             campaign.Title,
		Description:       description,
		Goal:              campaign.Goal.String(),
		Deadline:          campaign.Deadline.String(),
		AmountRaised:      campaign.AmountRaised.String(),
		Withdrawn:         campaign.Withdrawn,
		Status:            campaignStatus(campaign),
		Country:           country,
		Category:          category,
		CoverURL:          coverURL,
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
	scopeMaskSecret := os.Getenv("SCOPE_MASK_SECRET")
	if rpcURL == "" || contractAddress == "" || jwtSecret == "" ||
		postgresUser == "" || postgresPassword == "" || postgresDB == "" || postgresPort == "" ||
		auth0Domain == "" || auth0Audience == "" || scopeMaskSecret == "" {
		log.Fatal("RPC_URL, CONTRACT_ADDRESS, JWT_SECRET, POSTGRES_*, AUTH0_*, and SCOPE_MASK_SECRET variables must be set")
	}

	auth0KeyFunc, err := newAuth0KeyFunc(auth0Domain)
	if err != nil {
		log.Fatalf("failed to load auth0 jwks: %v", err)
	}

	idMasker, err := newIDMasker(scopeMaskSecret)
	if err != nil {
		log.Fatalf("failed to init id masker: %v", err)
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
			maskedID, err := maskID(idMasker, campaign.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			items[i] = gin.H{
				"id":             maskedID,
				"title":          campaign.Title,
				"description":    campaign.Description,
				"targetEth":      campaign.TargetEth,
				"country":        campaign.Country,
				"category":       campaign.Category,
				"durationDays":   campaign.DurationDays,
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
			Category       string   `json:"category" binding:"required"`
			Title          string   `json:"title" binding:"required"`
			Description    string   `json:"description"`
			TargetEth      string   `json:"targetEth" binding:"required"`
			DurationDays   uint32   `json:"durationDays" binding:"required,min=1,max=365"`
			FundraisingFor string   `json:"fundraisingFor" binding:"required"`
			AssetIDs       []uint64 `json:"assetIds" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid campaign fields"})
			return
		}
		if !models.IsValidCampaignCategory(req.Category) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
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

		campaign, err := models.CreateCampaign(gormDB, sub, req.Country, req.Category, req.Title, req.Description, req.TargetEth, req.DurationDays, req.FundraisingFor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := models.AttachAssetsToCampaign(gormDB, campaign.ID, req.AssetIDs, req.AssetIDs[0]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		maskedID, err := maskID(idMasker, campaign.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":             maskedID,
			"ownerSub":       campaign.OwnerSub,
			"country":        campaign.Country,
			"category":       campaign.Category,
			"title":          campaign.Title,
			"description":    campaign.Description,
			"targetEth":      campaign.TargetEth,
			"durationDays":   campaign.DurationDays,
			"fundraisingFor": campaign.FundraisingFor,
			"status":         campaign.Status,
			"createdAt":      campaign.CreatedAt,
		})
	})

	api.GET("/my-campaigns/:id", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		campaignID, err := unmaskID(idMasker, c.Param("id"))
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

		maskedID, err := maskID(idMasker, campaign.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                maskedID,
			"ownerSub":          campaign.OwnerSub,
			"country":           campaign.Country,
			"category":          campaign.Category,
			"title":             campaign.Title,
			"description":       campaign.Description,
			"targetEth":         campaign.TargetEth,
			"durationDays":      campaign.DurationDays,
			"fundraisingFor":    campaign.FundraisingFor,
			"status":            campaign.Status,
			"walletAddress":     campaign.WalletAddress,
			"onChainCampaignId": campaign.OnChainCampaignID,
			"publishedAt":       campaign.PublishedAt,
			"createdAt":         campaign.CreatedAt,
			"assets":            assets,
		})
	})

	api.DELETE("/my-campaigns/:id", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		campaignID, err := unmaskID(idMasker, c.Param("id"))
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
		if campaign.Status != models.CampaignStatusDraft {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only draft campaigns can be deleted"})
			return
		}

		orphanAssets, err := models.GetOrphanableAssetsForCampaign(gormDB, campaign.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		orphanAssetIDs := make([]uint64, len(orphanAssets))
		for i, asset := range orphanAssets {
			orphanAssetIDs[i] = asset.ID
		}

		if err := models.DeleteCampaign(gormDB, campaign.ID, orphanAssetIDs); err != nil {
			if errors.Is(err, models.ErrCampaignNotDraft) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "only draft campaigns can be deleted"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, asset := range orphanAssets {
			if err := deleteObjectFromR2(c.Request.Context(), r2Client, asset.Bucket, asset.ObjectKey); err != nil {
				log.Printf("failed to delete r2 object %s/%s: %v", asset.Bucket, asset.ObjectKey, err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	api.POST("/my-campaigns/:id/publish", auth0Middleware(auth0KeyFunc, auth0Domain, auth0Audience), func(c *gin.Context) {
		campaignID, err := unmaskID(idMasker, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		var req struct {
			WalletAddress     string `json:"walletAddress" binding:"required"`
			OnChainCampaignID uint64 `json:"onChainCampaignId"`
			TxHash            string `json:"txHash"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "walletAddress and onChainCampaignId are required"})
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
		if campaign.Status != models.CampaignStatusDraft {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only draft campaigns can be published"})
			return
		}

		onChainCampaign, err := crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(req.OnChainCampaignID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not find this campaign on-chain"})
			return
		}
		if !strings.EqualFold(onChainCampaign.Owner.Hex(), req.WalletAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "on-chain campaign owner does not match the provided wallet address"})
			return
		}

		updated, err := models.PublishCampaign(gormDB, campaign.ID, req.WalletAddress, req.OnChainCampaignID)
		if err != nil {
			if errors.Is(err, models.ErrCampaignNotDraft) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "only draft campaigns can be published"})
				return
			}
			if errors.Is(err, models.ErrOnChainCampaignAlreadyLinked) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		maskedID, err := maskID(idMasker, updated.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("campaign %d published on-chain as %d (tx %s)", updated.ID, req.OnChainCampaignID, req.TxHash)

		c.JSON(http.StatusOK, gin.H{
			"id":                maskedID,
			"status":            updated.Status,
			"walletAddress":     updated.WalletAddress,
			"onChainCampaignId": updated.OnChainCampaignID,
			"publishedAt":       updated.PublishedAt,
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
		category := c.Query("category")
		if category != "" && !models.IsValidCampaignCategory(category) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
			return
		}

		total, err := models.CountPublishedCampaigns(gormDB, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": total})
	})

	api.GET("/campaigns", func(c *gin.Context) {
		pagination, err := parsePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category := c.Query("category")
		if category != "" && !models.IsValidCampaignCategory(category) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
			return
		}

		dbCampaigns, total, err := models.ListPublishedCampaigns(gormDB, category, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		campaignIDs := make([]uint64, len(dbCampaigns))
		for i, dbCampaign := range dbCampaigns {
			campaignIDs[i] = dbCampaign.ID
		}

		coverURLs, err := models.GetCoverAssetsForCampaigns(gormDB, campaignIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items := make([]CampaignResponse, 0, len(dbCampaigns))
		for _, dbCampaign := range dbCampaigns {
			if dbCampaign.OnChainCampaignID == nil {
				continue
			}

			onChainID := *dbCampaign.OnChainCampaignID
			chainCampaign, err := crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(onChainID))
			if err != nil {
				continue
			}

			maskedID, err := maskID(idMasker, dbCampaign.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			items = append(items, toCampaignResponse(maskedID, onChainID, chainCampaign, &dbCampaign, coverURLs[dbCampaign.ID]))
		}

		c.JSON(http.StatusOK, PaginatedResponse{
			Items:  items,
			Total:  total,
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
	})

	api.GET("/campaigns/:id", func(c *gin.Context) {
		dbID, err := unmaskID(idMasker, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		dbCampaign, err := models.GetCampaignByID(gormDB, dbID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if dbCampaign == nil || dbCampaign.Status != models.CampaignStatusPublished || dbCampaign.OnChainCampaignID == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
			return
		}

		onChainID := *dbCampaign.OnChainCampaignID
		chainCampaign, err := crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(onChainID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		coverURLs, err := models.GetCoverAssetsForCampaigns(gormDB, []uint64{dbCampaign.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toCampaignResponse(c.Param("id"), onChainID, chainCampaign, dbCampaign, coverURLs[dbCampaign.ID]))
	})

	api.GET("/campaigns/:id/contributors", func(c *gin.Context) {
		dbID, err := unmaskID(idMasker, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
			return
		}

		dbCampaign, err := models.GetCampaignByID(gormDB, dbID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if dbCampaign == nil || dbCampaign.OnChainCampaignID == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
			return
		}

		summaries, err := models.GetContributorsForCampaign(gormDB, *dbCampaign.OnChainCampaignID)
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
