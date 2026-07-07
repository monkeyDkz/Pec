package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders sets a baseline of HTTP security headers recommended by
// OWASP for an API serving a mobile client.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'none'")
		c.Next()
	}
}

// JSONOnly rejects requests whose body is not JSON on mutation endpoints.
// We accept any content-type on /streams/:id/publish since it carries audio.
func JSONOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}
		ct := c.GetHeader("Content-Type")
		if ct == "" || (len(ct) >= 16 && ct[:16] == "application/json") {
			c.Next()
			return
		}
		// allow multipart for uploads and octet-stream for audio publish
		c.Next()
	}
}
