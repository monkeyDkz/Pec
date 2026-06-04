package persistence

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type statsRepository struct {
	db *gorm.DB
}

// NewStatsRepository returns a GORM-backed StatsRepository.
func NewStatsRepository(db *gorm.DB) repository.StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) Gather(ctx context.Context) (repository.Stats, error) {
	db := r.db.WithContext(ctx)
	var s repository.Stats

	if err := db.Model(&entity.User{}).Count(&s.TotalUsers).Error; err != nil {
		return s, fmt.Errorf("count users: %w", err)
	}
	if err := db.Model(&entity.Stream{}).Count(&s.TotalStreams).Error; err != nil {
		return s, fmt.Errorf("count streams: %w", err)
	}
	if err := db.Model(&entity.Stream{}).Where("status = ?", entity.StreamStatusLive).Count(&s.LiveStreams).Error; err != nil {
		return s, fmt.Errorf("count live streams: %w", err)
	}
	if err := db.Model(&entity.Track{}).Count(&s.TotalTracks).Error; err != nil {
		return s, fmt.Errorf("count tracks: %w", err)
	}
	if err := db.Model(&entity.Playlist{}).Count(&s.TotalPlaylists).Error; err != nil {
		return s, fmt.Errorf("count playlists: %w", err)
	}
	return s, nil
}
