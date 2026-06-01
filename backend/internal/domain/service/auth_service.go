package service

import (
	"context"

	"github.com/streampulse/backend/internal/domain/entity"
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}
