package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerCommentRoutes(api *gin.RouterGroup, deps *Dependencies) {
	api.GET("/campaigns/:id/comments", func(c *gin.Context) {
		pagination, err := parsePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items, total, err := deps.CommentService.ListComments(c.Request.Context(), c.Param("id"), pagination.Offset, pagination.Limit)
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
