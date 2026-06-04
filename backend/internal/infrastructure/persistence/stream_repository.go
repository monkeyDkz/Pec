package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type streamRepository struct {
	db *gorm.DB
}

// NewStreamRepository returns a GORM-backed StreamRepository.
func NewStreamRepository(db *gorm.DB) repository.StreamRepository {
	return &streamRepository{db: db}
}

func (r *streamRepository) Create(ctx context.Context, stream *entity.Stream) error {
	if err := r.db.WithContext(ctx).Create(stream).Error; err != nil {
		return fmt.Errorf("create stream: %w", err)
	}
	return nil
}

func (r *streamRepository) FindByID(ctx context.Context, id string) (*entity.Stream, error) {
	var stream entity.Stream
	err := r.db.WithContext(ctx).Preload("Broadcaster").First(&stream, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find stream by id: %w", err)
	}
	return &stream, nil
}

func (r *streamRepository) Update(ctx context.Context, stream *entity.Stream) error {
	if err := r.db.WithContext(ctx).Save(stream).Error; err != nil {
		return fmt.Errorf("update stream: %w", err)
	}
	return nil
}

func (r *streamRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Stream{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete stream: %w", err)
	}
	return nil
}

func (r *streamRepository) ListLive(ctx context.Context) ([]entity.Stream, error) {
	var streams []entity.Stream
	err := r.db.WithContext(ctx).
		Preload("Broadcaster").
		Where("status = ?", entity.StreamStatusLive).
		Order("created_at desc").
		Find(&streams).Error
	if err != nil {
		return nil, fmt.Errorf("list live streams: %w", err)
	}
	return streams, nil
}

func (r *streamRepository) ListByBroadcaster(ctx context.Context, broadcasterID string) ([]entity.Stream, error) {
	var streams []entity.Stream
	err := r.db.WithContext(ctx).
		Preload("Broadcaster").
		Where("broadcaster_id = ?", broadcasterID).
		Order("created_at desc").
		Find(&streams).Error
	if err != nil {
		return nil, fmt.Errorf("list streams by broadcaster: %w", err)
	}
	return streams, nil
}
