package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("PORT", "9999")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "postgres://x", cfg.DatabaseURL)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, "9999", cfg.Port)
	assert.Equal(t, 24*time.Hour, cfg.JWTDuration) // default
}

func TestLoad_MissingRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")
	_, err := Load()
	require.Error(t, err)
}
