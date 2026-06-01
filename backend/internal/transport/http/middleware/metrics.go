package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/infrastructure/observability"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		observability.HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
		observability.HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
	}
}
