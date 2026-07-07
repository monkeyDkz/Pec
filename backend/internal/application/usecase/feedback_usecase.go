package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

var ErrInvalidRating = errors.New("rating must be between 1 and 5")

type FeedbackUseCase struct {
	repo repository.FeedbackRepository
}

func NewFeedbackUseCase(repo repository.FeedbackRepository) *FeedbackUseCase {
	return &FeedbackUseCase{repo: repo}
}

func (uc *FeedbackUseCase) Submit(ctx context.Context, userID string, rating int, comment string) (*entity.Feedback, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}
	f := &entity.Feedback{UserID: userID, Rating: rating, Comment: comment}
	if err := uc.repo.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("submit feedback: %w", err)
	}
	return f, nil
}

func (uc *FeedbackUseCase) AdminList(ctx context.Context, offset, limit int) ([]entity.Feedback, error) {
	return uc.repo.List(ctx, offset, limit)
}
