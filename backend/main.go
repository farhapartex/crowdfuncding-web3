package main

import (
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crowdfunding-backend/contract"
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
	if rpcURL == "" || contractAddress == "" {
		log.Fatal("RPC_URL and CONTRACT_ADDRESS must be set")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to rpc: %v", err)
	}

	crowdFunding, err := contract.NewCrowdFunding(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatalf("failed to bind contract: %v", err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/campaigns/count", func(c *gin.Context) {
		count, err := crowdFunding.CampaignCount(nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count.String()})
	})

	router.GET("/campaigns", func(c *gin.Context) {
		campaigns, err := crowdFunding.GetCampaigns(nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]CampaignResponse, len(campaigns))
		for i, campaign := range campaigns {
			response[i] = toCampaignResponse(uint64(i), campaign)
		}

		c.JSON(http.StatusOK, response)
	})

	router.GET("/campaigns/:id", func(c *gin.Context) {
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
