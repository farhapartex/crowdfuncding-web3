package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAssetRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.POST("/assets", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		asset, err := deps.AssetService.UploadAsset(c.Request.Context(), c.GetString("sub"), fileHeader)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusCreated, asset)
	})
}
