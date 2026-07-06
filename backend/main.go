package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/db"
	"crowdfunding-backend/models"
)

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
	if rpcURL == "" || contractAddress == "" || jwtSecret == "" ||
		postgresUser == "" || postgresPassword == "" || postgresDB == "" || postgresPort == "" {
		log.Fatal("RPC_URL, CONTRACT_ADDRESS, JWT_SECRET, and POSTGRES_* variables must be set")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
