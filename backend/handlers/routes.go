package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crowdfunding-backend/services"
)

type Dependencies struct {
	AuthService           *services.AuthService
	Auth0Service          *services.Auth0Service
	ProfileService        *services.ProfileService
	AssetService          *services.AssetService
	CampaignService       *services.CampaignService
	PublicCampaignService *services.PublicCampaignService
	WalletService         *services.WalletService
	CommentService        *services.CommentService
}

func RegisterRoutes(router *gin.Engine, deps *Dependencies) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	registerAuthRoutes(api, deps)
	registerAuth0Routes(api, deps)
	registerAssetRoutes(api, deps)
	registerCampaignRoutes(api, deps)
	registerPublicCampaignRoutes(api, deps)
	registerProfileRoutes(api, deps)
	registerWalletRoutes(api, deps)
	registerCommentRoutes(api, deps)
}
