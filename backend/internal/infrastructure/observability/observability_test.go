package observability

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	for _, lvl := range []string{"debug", "warn", "error", "info", "unknown"} {
		InitLogger(lvl)
	}
}

func TestInitTracer(t *testing.T) {
	shutdown, err := InitTracer("localhost:4317", "test-service")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	require.NoError(t, shutdown(ctx))
}
