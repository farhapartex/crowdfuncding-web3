package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crowdfunding-backend/services"
)

func registerCampaignRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/my-campaigns", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		pagination, err := parsePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items, total, err := deps.CampaignService.ListMyCampaigns(c.GetString("sub"), pagination.Offset, pagination.Limit)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, PaginatedResponse{
			Items:  items,
			Total:  total,
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
	})

	api.POST("/my-campaigns", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		var req struct {
			Country        string   `json:"country" binding:"required"`
			Category       string   `json:"category" binding:"required"`
			Title          string   `json:"title" binding:"required"`
			Description    string   `json:"description"`
			CurrencyMode   string   `json:"currencyMode" binding:"required"`
			TargetEth      string   `json:"targetEth"`
			TokenAddress   string   `json:"tokenAddress"`
			GoalToken      string   `json:"goalToken"`
			DurationDays   uint32   `json:"durationDays" binding:"required,min=1,max=365"`
			FundraisingFor string   `json:"fundraisingFor" binding:"required"`
			AssetIDs       []uint64 `json:"assetIds" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid campaign fields"})
			return
		}

		result, err := deps.CampaignService.CreateCampaign(c.GetString("sub"), services.CreateCampaignInput{
			Country:        req.Country,
			Category:       req.Category,
			Title:          req.Title,
			Description:    req.Description,
			CurrencyMode:   req.CurrencyMode,
			TargetEth:      req.TargetEth,
			TokenAddress:   req.TokenAddress,
			GoalToken:      req.GoalToken,
			DurationDays:   req.DurationDays,
			FundraisingFor: req.FundraisingFor,
			AssetIDs:       req.AssetIDs,
		})
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusCreated, result)
	})

	api.GET("/my-campaigns/:id", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		result, err := deps.CampaignService.GetMyCampaign(c.GetString("sub"), c.Param("id"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})

	api.DELETE("/my-campaigns/:id", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		err := deps.CampaignService.DeleteCampaign(c.Request.Context(), c.GetString("sub"), c.Param("id"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	api.POST("/my-campaigns/:id/publish", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		var req struct {
			WalletAddress     string `json:"walletAddress" binding:"required"`
			OnChainCampaignID uint64 `json:"onChainCampaignId"`
			TxHash            string `json:"txHash"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "walletAddress and onChainCampaignId are required"})
			return
		}

		result, err := deps.CampaignService.PublishCampaign(c.GetString("sub"), c.Param("id"), services.PublishCampaignInput{
			WalletAddress:     req.WalletAddress,
			OnChainCampaignID: req.OnChainCampaignID,
			TxHash:            req.TxHash,
		})
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})

	api.POST("/my-campaigns/:id/archive", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		var req struct {
			Note string `json:"note" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a note is required when archiving a campaign"})
			return
		}

		result, err := deps.CampaignService.ArchiveCampaign(c.GetString("sub"), c.Param("id"), req.Note)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
