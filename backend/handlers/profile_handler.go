package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerProfileRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/profiles/:address", func(c *gin.Context) {
		profile, err := deps.ProfileService.GetProfile(c.Param("address"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"address": profile.Address, "displayName": profile.DisplayName})
	})
}
