package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	m := NewJWTManager("super-secret", time.Hour)

	token, err := m.Generate("user-123", "admin")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := m.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "admin", claims.Role)
}

func TestJWTManager_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name:    "garbage token",
			setup:   func() string { return "not-a-jwt" },
			wantErr: true,
		},
		{
			name: "wrong secret",
			setup: func() string {
				other := NewJWTManager("different-secret", time.Hour)
				tok, _ := other.Generate("u", "user")
				return tok
			},
			wantErr: true,
		},
		{
			name: "expired token",
			setup: func() string {
				expired := NewJWTManager("super-secret", -time.Hour)
				tok, _ := expired.Generate("u", "user")
				return tok
			},
			wantErr: true,
		},
	}

	m := NewJWTManager("super-secret", time.Hour)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.Validate(tt.setup())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
