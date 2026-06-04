package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

// UserUseCase exposes profile operations for the authenticated user.
type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UserUseCase) UpdateUsername(ctx context.Context, id, username string) (*entity.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if username != "" {
		user.Username = username
	}
	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}
