package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAuth0Routes(api *gin.RouterGroup, deps *Dependencies) {
	api.POST("/auth0/sync", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		user, err := deps.Auth0Service.SyncUser(c.GetString("sub"), c.GetString("auth0Token"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, user)
	})

	api.GET("/auth0/me", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		user, err := deps.Auth0Service.GetUser(c.GetString("sub"))
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusOK, user)
	})
}
