package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamUseCase_Listen(t *testing.T) {
	uc := NewStreamUseCase(&fakeEngine{}, newFakeStreamRepoUC())
	r, err := uc.Listen(context.Background(), "any-id")
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NoError(t, r.Close())
}

func TestTrackUseCase_GetAndDelete(t *testing.T) {
	ctx := context.Background()
	uc := NewTrackUseCase(newFakeTrackRepo())

	tr, err := uc.Create(ctx, "uploader", "Title", "Artist", 100, "https://cdn/x.mp3")
	require.NoError(t, err)

	got, err := uc.Get(ctx, tr.ID)
	require.NoError(t, err)
	assert.Equal(t, tr.ID, got.ID)

	require.NoError(t, uc.Delete(ctx, tr.ID))
}
