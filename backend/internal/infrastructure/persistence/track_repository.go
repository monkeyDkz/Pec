package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type trackRepository struct {
	db *gorm.DB
}

// NewTrackRepository returns a GORM-backed TrackRepository.
func NewTrackRepository(db *gorm.DB) repository.TrackRepository {
	return &trackRepository{db: db}
}

func (r *trackRepository) Create(ctx context.Context, track *entity.Track) error {
	if err := r.db.WithContext(ctx).Create(track).Error; err != nil {
		return fmt.Errorf("create track: %w", err)
	}
	return nil
}

func (r *trackRepository) FindByID(ctx context.Context, id string) (*entity.Track, error) {
	var track entity.Track
	err := r.db.WithContext(ctx).First(&track, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find track by id: %w", err)
	}
	return &track, nil
}

func (r *trackRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Track{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete track: %w", err)
	}
	return nil
}

func (r *trackRepository) List(ctx context.Context, offset, limit int) ([]entity.Track, error) {
	var tracks []entity.Track
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("list tracks: %w", err)
	}
	return tracks, nil
}
