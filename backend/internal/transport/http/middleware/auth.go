package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/infrastructure/auth"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	jwtManager := auth.NewJWTManager(jwtSecret, 0)

	return func(c *gin.Context) {
		var tokenStr string

		if header := c.GetHeader("Authorization"); header != "" {
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
				return
			}
			tokenStr = parts[1]
		} else if q := c.Query("token"); q != "" {
			// Fallback for browser media elements (audio/video) that cannot set
			// an Authorization header — e.g. the web build of the player.
			tokenStr = q
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		claims, err := jwtManager.Validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no role found"})
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
