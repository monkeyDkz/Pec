package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// dailySalt rotates the IP hash daily — enough granularity to debug a session
// without persisting raw IPs (GDPR-friendly).
var dailySalt = func() string {
	return time.Now().UTC().Format("2006-01-02")
}

// RequestLogger emits one structured log line per HTTP request. The trace_id
// is automatically injected by the slog traceHandler when a span is active.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// IP is hashed before logging — GDPR / Ce3.1.4 compliance.
		h := sha256.Sum256([]byte(c.ClientIP() + dailySalt()))
		ipHash := hex.EncodeToString(h[:])[:16]

		slog.InfoContext(c.Request.Context(), "http_request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"ip_hash", ipHash,
			"user_agent", c.Request.UserAgent(),
		)
	}
}
