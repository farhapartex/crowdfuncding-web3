package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerTokenRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/tokens", func(c *gin.Context) {
		c.JSON(http.StatusOK, deps.TokenService.List())
	})
}
