package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

// AdminUseCase exposes platform administration: user management and global stats.
type AdminUseCase struct {
	userRepo  repository.UserRepository
	statsRepo repository.StatsRepository
}

func NewAdminUseCase(userRepo repository.UserRepository, statsRepo repository.StatsRepository) *AdminUseCase {
	return &AdminUseCase{userRepo: userRepo, statsRepo: statsRepo}
}

func (uc *AdminUseCase) ListUsers(ctx context.Context, offset, limit int) ([]entity.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return uc.userRepo.List(ctx, offset, limit)
}

func (uc *AdminUseCase) UpdateRole(ctx context.Context, userID, role string) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	user.Role = entity.Role(role)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	return user, nil
}

func (uc *AdminUseCase) Stats(ctx context.Context) (repository.Stats, error) {
	return uc.statsRepo.Gather(ctx)
}
