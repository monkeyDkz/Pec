package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type playlistRepository struct {
	db *gorm.DB
}

// NewPlaylistRepository returns a GORM-backed PlaylistRepository.
func NewPlaylistRepository(db *gorm.DB) repository.PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) Create(ctx context.Context, playlist *entity.Playlist) error {
	if err := r.db.WithContext(ctx).Create(playlist).Error; err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}
	return nil
}

func (r *playlistRepository) FindByID(ctx context.Context, id string) (*entity.Playlist, error) {
	var playlist entity.Playlist
	err := r.db.WithContext(ctx).Preload("Tracks").Preload("Owner").First(&playlist, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find playlist by id: %w", err)
	}
	return &playlist, nil
}

func (r *playlistRepository) Update(ctx context.Context, playlist *entity.Playlist) error {
	if err := r.db.WithContext(ctx).Save(playlist).Error; err != nil {
		return fmt.Errorf("update playlist: %w", err)
	}
	return nil
}

func (r *playlistRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Select("Tracks").Delete(&entity.Playlist{ID: id}).Error; err != nil {
		return fmt.Errorf("delete playlist: %w", err)
	}
	return nil
}

func (r *playlistRepository) ListByOwner(ctx context.Context, ownerID string) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	err := r.db.WithContext(ctx).
		Preload("Tracks").
		Where("owner_id = ?", ownerID).
		Order("created_at desc").
		Find(&playlists).Error
	if err != nil {
		return nil, fmt.Errorf("list playlists by owner: %w", err)
	}
	return playlists, nil
}

func (r *playlistRepository) AddTrack(ctx context.Context, playlistID, trackID string) error {
	// Insert directly into the join table: Association().Append would attempt to
	// upsert a partial Track (empty NOT NULL columns) and fail.
	err := r.db.WithContext(ctx).Exec(
		`INSERT INTO playlist_tracks (playlist_id, track_id) VALUES (?, ?) ON CONFLICT DO NOTHING`,
		playlistID, trackID,
	).Error
	if err != nil {
		return fmt.Errorf("add track to playlist: %w", err)
	}
	return nil
}

func (r *playlistRepository) RemoveTrack(ctx context.Context, playlistID, trackID string) error {
	err := r.db.WithContext(ctx).Exec(
		`DELETE FROM playlist_tracks WHERE playlist_id = ? AND track_id = ?`,
		playlistID, trackID,
	).Error
	if err != nil {
		return fmt.Errorf("remove track from playlist: %w", err)
	}
	return nil
}
