package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/auth"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
	jwt      *auth.JWTManager
	hasher   auth.PasswordHasher
}

func NewAuthUseCase(userRepo repository.UserRepository, jwt *auth.JWTManager, hasher auth.PasswordHasher) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		jwt:      jwt,
		hasher:   hasher,
	}
}

// Register creates a new user and returns it along with a freshly issued JWT so
// the client is authenticated immediately after sign-up.
func (uc *AuthUseCase) Register(ctx context.Context, email, username, password string) (*entity.User, string, error) {
	hashed, err := uc.hasher.Hash(password)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		Email:    email,
		Username: username,
		Password: hashed,
		Role:     entity.RoleUser,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}

	token, err := uc.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, "", fmt.Errorf("generate token: %w", err)
	}

	return user, token, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}

	if !uc.hasher.Verify(password, user.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := uc.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	claims, err := uc.jwt.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}
