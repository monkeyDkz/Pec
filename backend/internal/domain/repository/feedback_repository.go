package repository

import (
	"context"

	"github.com/streampulse/backend/internal/domain/entity"
)

type FeedbackRepository interface {
	Create(ctx context.Context, f *entity.Feedback) error
	List(ctx context.Context, offset, limit int) ([]entity.Feedback, error)
	ListByUser(ctx context.Context, userID string) ([]entity.Feedback, error)
}
