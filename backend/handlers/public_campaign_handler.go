package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerPublicCampaignRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/campaigns/count", func(c *gin.Context) {
		total, err := deps.PublicCampaignService.CountPublished(c.Query("category"))
		if err != nil {
			respondError(c, err)
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

		items, total, err := deps.PublicCampaignService.ListPublished(c.Query("category"), pagination.Offset, pagination.Limit)
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

	api.GET("/campaigns/:id", func(c *gin.Context) {
		result, err := deps.PublicCampaignService.GetPublished(c.Param("id"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})

	api.GET("/campaigns/:id/contributors", func(c *gin.Context) {
		result, err := deps.PublicCampaignService.GetContributors(c.Param("id"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})

	api.GET("/campaigns/:id/transactions", func(c *gin.Context) {
		pagination, err := parsePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items, total, err := deps.PublicCampaignService.GetCampaignTransactions(c.Param("id"), pagination.Offset, pagination.Limit)
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
}
