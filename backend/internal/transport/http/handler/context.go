package handler

import (
	"context"
	"time"
)

// newDetachedContext returns a context that survives the request's context
// cancellation. Use it for cleanup steps (e.g. marking a stream offline) that
// must run even after the client disconnects.
func newDetachedContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
