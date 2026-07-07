package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type PlaylistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

var _ repository.PlaylistRepository = (*PlaylistRepository)(nil)

func (r *PlaylistRepository) Create(ctx context.Context, p *entity.Playlist) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}
	return nil
}

func (r *PlaylistRepository) FindByID(ctx context.Context, id string) (*entity.Playlist, error) {
	var p entity.Playlist
	if err := r.db.WithContext(ctx).
		Preload("Tracks").
		Preload("Owner").
		First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find playlist: %w", err)
	}
	return &p, nil
}

func (r *PlaylistRepository) Update(ctx context.Context, p *entity.Playlist) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("update playlist: %w", err)
	}
	return nil
}

func (r *PlaylistRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Playlist{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete playlist: %w", err)
	}
	return nil
}

func (r *PlaylistRepository) ListByOwner(ctx context.Context, ownerID string) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&playlists).Error; err != nil {
		return nil, fmt.Errorf("list playlists: %w", err)
	}
	return playlists, nil
}

func (r *PlaylistRepository) AddTrack(ctx context.Context, playlistID, trackID string) error {
	return r.db.WithContext(ctx).
		Exec(`INSERT INTO playlist_tracks (playlist_id, track_id, position)
		      VALUES (?, ?, COALESCE((SELECT MAX(position) + 1 FROM playlist_tracks WHERE playlist_id = ?), 0))
		      ON CONFLICT DO NOTHING`,
			playlistID, trackID, playlistID).Error
}

func (r *PlaylistRepository) RemoveTrack(ctx context.Context, playlistID, trackID string) error {
	return r.db.WithContext(ctx).
		Exec(`DELETE FROM playlist_tracks WHERE playlist_id = ? AND track_id = ?`,
			playlistID, trackID).Error
}

// ReorderTracks atomically rewrites the position column for every track in
// the playlist. Any track id not in the list keeps its current position.
func (r *PlaylistRepository) ReorderTracks(ctx context.Context, playlistID string, orderedTrackIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, trackID := range orderedTrackIDs {
			if err := tx.Exec(
				`UPDATE playlist_tracks SET position = ? WHERE playlist_id = ? AND track_id = ?`,
				i, playlistID, trackID,
			).Error; err != nil {
				return fmt.Errorf("reorder position %d: %w", i, err)
			}
		}
		return nil
	})
}

type TrackRepository struct {
	db *gorm.DB
}

func NewTrackRepository(db *gorm.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

var _ repository.TrackRepository = (*TrackRepository)(nil)

func (r *TrackRepository) Create(ctx context.Context, t *entity.Track) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create track: %w", err)
	}
	return nil
}

func (r *TrackRepository) FindByID(ctx context.Context, id string) (*entity.Track, error) {
	var t entity.Track
	if err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find track: %w", err)
	}
	return &t, nil
}

func (r *TrackRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Track{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete track: %w", err)
	}
	return nil
}

func (r *TrackRepository) List(ctx context.Context, offset, limit int) ([]entity.Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var tracks []entity.Track
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("list tracks: %w", err)
	}
	return tracks, nil
}
