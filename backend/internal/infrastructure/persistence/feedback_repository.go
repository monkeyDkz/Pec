package persistence

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type FeedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

var _ repository.FeedbackRepository = (*FeedbackRepository)(nil)

func (r *FeedbackRepository) Create(ctx context.Context, f *entity.Feedback) error {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return fmt.Errorf("create feedback: %w", err)
	}
	return nil
}

func (r *FeedbackRepository) List(ctx context.Context, offset, limit int) ([]entity.Feedback, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	var out []entity.Feedback
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list feedback: %w", err)
	}
	return out, nil
}

func (r *FeedbackRepository) ListByUser(ctx context.Context, userID string) ([]entity.Feedback, error) {
	var out []entity.Feedback
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list feedback by user: %w", err)
	}
	return out, nil
}
