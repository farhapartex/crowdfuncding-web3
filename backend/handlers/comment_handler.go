package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func resolveAuthorName(deps *Dependencies, sub string) string {
	if user, err := deps.Auth0Service.GetUser(sub); err == nil && user.DisplayName != "" {
		return user.DisplayName
	}
	return sub
}

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

	api.POST("/campaigns/:id/comments", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
			return
		}

		if err := deps.PublicCampaignService.EnsureCommentable(c.Param("id")); err != nil {
			respondError(c, err)
			return
		}

		sub := c.GetString("sub")
		authorName := resolveAuthorName(deps, sub)

		comment, err := deps.CommentService.PostComment(c.Request.Context(), c.Param("id"), sub, authorName, req.Text)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusCreated, comment)
	})

	api.POST("/comments/:commentId/replies", auth0Middleware(deps.Auth0Service), func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
			return
		}

		parent, err := deps.CommentService.GetComment(c.Request.Context(), c.Param("commentId"))
		if err != nil {
			respondError(c, err)
			return
		}

		if err := deps.PublicCampaignService.EnsureCommentable(parent.CampaignID); err != nil {
			respondError(c, err)
			return
		}

		sub := c.GetString("sub")
		authorName := resolveAuthorName(deps, sub)

		reply, err := deps.CommentService.ReplyToComment(c.Request.Context(), c.Param("commentId"), sub, authorName, req.Text)
		if err != nil {
			respondError(c, err)
			return
		}

		c.JSON(http.StatusCreated, reply)
	})
}
