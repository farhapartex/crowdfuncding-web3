package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/auth/nonce", func(c *gin.Context) {
		nonce, message, err := deps.AuthService.IssueNonce(c.Query("address"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"nonce": nonce, "message": message})
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

		token, address, err := deps.AuthService.VerifyAndIssueSession(req.Address, req.Signature)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "address": address})
	})

	api.GET("/me", authMiddleware(deps.AuthService), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"address": c.GetString("address")})
	})

	api.GET("/me/profile", authMiddleware(deps.AuthService), func(c *gin.Context) {
		profile, err := deps.ProfileService.GetProfile(c.GetString("address"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, profile)
	})

	api.PUT("/me/profile", authMiddleware(deps.AuthService), func(c *gin.Context) {
		var req struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		profile, err := deps.ProfileService.UpsertProfile(c.GetString("address"), req.DisplayName, req.Email)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, profile)
	})
}
