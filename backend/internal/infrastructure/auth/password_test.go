package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher(t *testing.T) {
	h := NewBcryptHasher()

	hash, err := h.Hash("s3cr3t-password")
	require.NoError(t, err)
	assert.NotEqual(t, "s3cr3t-password", hash, "hash must not equal plaintext")

	assert.True(t, h.Verify("s3cr3t-password", hash), "correct password should verify")
	assert.False(t, h.Verify("wrong-password", hash), "wrong password must fail")
}
