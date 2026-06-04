//go:build integration

// Integration tests run against a real PostgreSQL instance.
//
//	docker compose up -d postgres
//	go test -tags=integration ./internal/infrastructure/persistence/...
//
// DATABASE_URL overrides the connection (defaults to the local docker-compose DB).
package persistence

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://streampulse:streampulse@localhost:5432/streampulse?sslmode=disable"
	}
	db, err := NewDB(dsn)
	require.NoError(t, err)
	require.NoError(t, AutoMigrate(db))
	return db
}

func makeUser(t *testing.T, repo interface {
	Create(context.Context, *entity.User) error
	Delete(context.Context, string) error
}) *entity.User {
	t.Helper()
	ctx := context.Background()
	u := &entity.User{
		Email:    uuid.NewString() + "@test.com",
		Username: "u-" + uuid.NewString()[:8],
		Password: "hashed",
		Role:     entity.RoleBroadcaster,
	}
	require.NoError(t, repo.Create(ctx, u))
	t.Cleanup(func() { _ = repo.Delete(ctx, u.ID) })
	return u
}

func TestUserRepository_Integration(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB(t))
	u := makeUser(t, repo)

	got, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.Email, got.Email)

	byEmail, err := repo.FindByEmail(ctx, u.Email)
	require.NoError(t, err)
	assert.Equal(t, u.ID, byEmail.ID)

	u.Username = "renamed"
	require.NoError(t, repo.Update(ctx, u))

	list, err := repo.List(ctx, 0, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, list)

	// missing id (valid uuid that does not exist) -> ErrNotFound
	_, err = repo.FindByID(ctx, uuid.NewString())
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStreamRepository_Integration(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	owner := makeUser(t, NewUserRepository(db))
	repo := NewStreamRepository(db)

	s := &entity.Stream{Title: "Radio", Description: "live", BroadcasterID: owner.ID, Status: entity.StreamStatusLive}
	require.NoError(t, repo.Create(ctx, s))
	t.Cleanup(func() { _ = repo.Delete(ctx, s.ID) })

	got, err := repo.FindByID(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, owner.Username, got.Broadcaster.Username) // preloaded

	live, err := repo.ListLive(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, live)

	byB, err := repo.ListByBroadcaster(ctx, owner.ID)
	require.NoError(t, err)
	assert.Len(t, byB, 1)

	s.Status = entity.StreamStatusOffline
	require.NoError(t, repo.Update(ctx, s))
}

func TestTrackAndPlaylistRepository_Integration(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	owner := makeUser(t, NewUserRepository(db))

	trackRepo := NewTrackRepository(db)
	track := &entity.Track{Title: "Song", Artist: "Artist", Duration: 200, FileURL: "https://cdn/x.mp3", UploadBy: owner.ID}
	require.NoError(t, trackRepo.Create(ctx, track))

	playlistRepo := NewPlaylistRepository(db)
	pl := &entity.Playlist{Name: "Chill", OwnerID: owner.ID}
	require.NoError(t, playlistRepo.Create(ctx, pl))

	// cleanups (LIFO): playlist clears the join, then track, then owner via makeUser cleanup.
	t.Cleanup(func() { _ = trackRepo.Delete(ctx, track.ID) })
	t.Cleanup(func() { _ = playlistRepo.Delete(ctx, pl.ID) })

	require.NoError(t, playlistRepo.AddTrack(ctx, pl.ID, track.ID))

	got, err := playlistRepo.FindByID(ctx, pl.ID)
	require.NoError(t, err)
	require.Len(t, got.Tracks, 1)
	assert.Equal(t, "Song", got.Tracks[0].Title)

	gotTrack, err := trackRepo.FindByID(ctx, track.ID)
	require.NoError(t, err)
	assert.Equal(t, "Song", gotTrack.Title)

	pl.Name = "Updated"
	require.NoError(t, playlistRepo.Update(ctx, pl))

	byOwner, err := playlistRepo.ListByOwner(ctx, owner.ID)
	require.NoError(t, err)
	assert.Len(t, byOwner, 1)

	require.NoError(t, playlistRepo.RemoveTrack(ctx, pl.ID, track.ID))
	got, err = playlistRepo.FindByID(ctx, pl.ID)
	require.NoError(t, err)
	assert.Len(t, got.Tracks, 0)

	tracks, err := trackRepo.List(ctx, 0, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, tracks)
}

func TestStatsRepository_Integration(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	_ = makeUser(t, NewUserRepository(db))
	repo := NewStatsRepository(db)

	stats, err := repo.Gather(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TotalUsers, int64(1))
}
