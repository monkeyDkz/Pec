package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Tracing returns the standard otelgin middleware. Each HTTP request gets a
// span named "<METHOD> <route>", child spans (DB, HTTP downstream) attach
// automatically thanks to W3C propagation.
func Tracing(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
