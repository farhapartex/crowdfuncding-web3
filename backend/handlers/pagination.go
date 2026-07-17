package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

type PaginationParams struct {
	Offset uint64
	Limit  uint64
}

func parsePagination(c *gin.Context) (PaginationParams, error) {
	offset, err := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		return PaginationParams{}, errors.New("invalid offset")
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", strconv.Itoa(defaultPageLimit)), 10, 64)
	if err != nil {
		return PaginationParams{}, errors.New("invalid limit")
	}
	if limit == 0 || limit > maxPageLimit {
		limit = defaultPageLimit
	}

	return PaginationParams{Offset: offset, Limit: limit}, nil
}

type PaginatedResponse struct {
	Items  any    `json:"items"`
	Total  int64  `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
