package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/observability"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

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

// Register creates a new account and returns both the user and a fresh JWT.
// The caller may receive repository.ErrConflict if the email or username is already taken.
func (uc *AuthUseCase) Register(ctx context.Context, email, username, password string) (*entity.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)

	if existing, err := uc.userRepo.FindByEmail(ctx, email); err == nil && existing != nil {
		return nil, "", repository.ErrConflict
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, "", fmt.Errorf("lookup email: %w", err)
	}

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

// Refresh re-issues a JWT for the user authenticated by an already-valid
// token. Sliding-window strategy: as long as the user keeps using the app
// before expiry, they get a fresh window without re-authenticating.
func (uc *AuthUseCase) Refresh(ctx context.Context, userID string) (*entity.User, string, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	token, err := uc.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, "", fmt.Errorf("generate token: %w", err)
	}
	return user, token, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			observability.AuthLoginsTotal.WithLabelValues("invalid_credentials").Inc()
			return nil, "", ErrInvalidCredentials
		}
		observability.AuthLoginsTotal.WithLabelValues("error").Inc()
		return nil, "", fmt.Errorf("find user: %w", err)
	}

	if !uc.hasher.Verify(password, user.Password) {
		observability.AuthLoginsTotal.WithLabelValues("invalid_credentials").Inc()
		return nil, "", ErrInvalidCredentials
	}

	token, err := uc.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		observability.AuthLoginsTotal.WithLabelValues("error").Inc()
		return nil, "", fmt.Errorf("generate token: %w", err)
	}
	observability.AuthLoginsTotal.WithLabelValues("success").Inc()
	return user, token, nil
}
