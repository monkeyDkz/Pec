// Package handler contains the Gin HTTP handlers that adapt transport concerns
// to application use cases.
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/domain/service"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
)

// currentUserID returns the authenticated user id set by the auth middleware.
func currentUserID(c *gin.Context) string {
	v, _ := c.Get("user_id")
	s, _ := v.(string)
	return s
}

// currentUserRole returns the authenticated user's role.
func currentUserRole(c *gin.Context) string {
	v, _ := c.Get("user_role")
	s, _ := v.(string)
	return s
}

// respondError maps domain/persistence errors to HTTP status codes.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, service.ErrStreamNotLive):
		c.JSON(http.StatusConflict, gin.H{"error": "stream is not live"})
	case errors.Is(err, persistence.ErrNotFound), errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
