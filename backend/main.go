package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/db"
	"crowdfunding-backend/handlers"
	"crowdfunding-backend/services"
)

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

	auth0Service, err := services.NewAuth0Service(gormDB, auth0Domain, auth0Audience)
	if err != nil {
		log.Fatalf("failed to load auth0 jwks: %v", err)
	}

	idMaskService, err := services.NewIDMaskService(scopeMaskSecret)
	if err != nil {
		log.Fatalf("failed to init id masker: %v", err)
	}

	r2AccountID := os.Getenv("R2_ACCOUNT_ID")
	r2AccessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2Bucket := os.Getenv("R2_BUCKET")
	r2PublicURL := os.Getenv("R2_PUBLIC_URL")

	r2Client, err := services.NewR2Client(context.Background(), r2AccountID, r2AccessKeyID, r2SecretAccessKey)
	if err != nil {
		log.Fatalf("failed to init r2 client: %v", err)
	}

	assetService := services.NewAssetService(gormDB, r2Client, r2Bucket, r2PublicURL)

	commentServiceAddr := os.Getenv("COMMENT_SERVICE_ADDR")
	commentServiceToken := os.Getenv("COMMENT_SERVICE_TOKEN")
	if commentServiceAddr == "" || commentServiceToken == "" {
		log.Fatal("COMMENT_SERVICE_ADDR and COMMENT_SERVICE_TOKEN must be set")
	}

	commentService, err := services.NewCommentService(commentServiceAddr, commentServiceToken)
	if err != nil {
		log.Fatalf("failed to connect to comment service: %v", err)
	}

	deps := &handlers.Dependencies{
		AuthService:           services.NewAuthService([]byte(jwtSecret)),
		Auth0Service:          auth0Service,
		ProfileService:        services.NewProfileService(gormDB),
		AssetService:          assetService,
		CampaignService:       services.NewCampaignService(gormDB, crowdFunding, idMaskService, assetService),
		PublicCampaignService: services.NewPublicCampaignService(gormDB, crowdFunding, idMaskService),
		WalletService:         services.NewWalletService(gormDB),
		CommentService:        commentService,
	}

	services.StartTransactionIndexer(gormDB, crowdFunding, client)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	handlers.RegisterRoutes(router, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
