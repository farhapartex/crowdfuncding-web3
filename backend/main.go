package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crowdfunding-backend/contract"
)

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
