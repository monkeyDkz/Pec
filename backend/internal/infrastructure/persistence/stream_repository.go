package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type StreamRepository struct {
	db *gorm.DB
}

func NewStreamRepository(db *gorm.DB) *StreamRepository {
	return &StreamRepository{db: db}
}

var _ repository.StreamRepository = (*StreamRepository)(nil)

func (r *StreamRepository) Create(ctx context.Context, s *entity.Stream) error {
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("create stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) FindByID(ctx context.Context, id string) (*entity.Stream, error) {
	var s entity.Stream
	if err := r.db.WithContext(ctx).Preload("Broadcaster").First(&s, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find stream: %w", err)
	}
	return &s, nil
}

func (r *StreamRepository) Update(ctx context.Context, s *entity.Stream) error {
	if err := r.db.WithContext(ctx).Save(s).Error; err != nil {
		return fmt.Errorf("update stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Stream{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) ListLive(ctx context.Context) ([]entity.Stream, error) {
	var streams []entity.Stream
	if err := r.db.WithContext(ctx).
		Preload("Broadcaster").
		Where("status = ?", entity.StreamStatusLive).
		Order("updated_at DESC").
		Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("list live streams: %w", err)
	}
	return streams, nil
}

func (r *StreamRepository) ListByBroadcaster(ctx context.Context, broadcasterID string) ([]entity.Stream, error) {
	var streams []entity.Stream
	if err := r.db.WithContext(ctx).
		Where("broadcaster_id = ?", broadcasterID).
		Order("created_at DESC").
		Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("list broadcaster streams: %w", err)
	}
	return streams, nil
}
