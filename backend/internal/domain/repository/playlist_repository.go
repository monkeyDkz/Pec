package repository

import (
	"context"

	"github.com/streampulse/backend/internal/domain/entity"
)

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *entity.Playlist) error
	FindByID(ctx context.Context, id string) (*entity.Playlist, error)
	Update(ctx context.Context, playlist *entity.Playlist) error
	Delete(ctx context.Context, id string) error
	ListByOwner(ctx context.Context, ownerID string) ([]entity.Playlist, error)
	AddTrack(ctx context.Context, playlistID, trackID string) error
	RemoveTrack(ctx context.Context, playlistID, trackID string) error
	ReorderTracks(ctx context.Context, playlistID string, orderedTrackIDs []string) error
}

type TrackRepository interface {
	Create(ctx context.Context, track *entity.Track) error
	FindByID(ctx context.Context, id string) (*entity.Track, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]entity.Track, error)
}
