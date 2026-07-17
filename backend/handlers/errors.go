package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"crowdfunding-backend/services"
)

func respondError(c *gin.Context, err error) {
	var validationErr *services.ValidationError
	var notFoundErr *services.NotFoundError
	var forbiddenErr *services.ForbiddenError
	var unauthorizedErr *services.UnauthorizedError
	var conflictErr *services.ConflictError
	var unavailableErr *services.UnavailableError

	switch {
	case errors.As(err, &validationErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.As(err, &notFoundErr):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.As(err, &forbiddenErr):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.As(err, &unauthorizedErr):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.As(err, &conflictErr):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.As(err, &unavailableErr):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
