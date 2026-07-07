package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

type UserUseCase struct {
	users     repository.UserRepository
	playlists repository.PlaylistRepository
	streams   repository.StreamRepository
}

func NewUserUseCase(
	users repository.UserRepository,
	playlists repository.PlaylistRepository,
	streams repository.StreamRepository,
) *UserUseCase {
	return &UserUseCase{users: users, playlists: playlists, streams: streams}
}

func (uc *UserUseCase) Me(ctx context.Context, id string) (*entity.User, error) {
	return uc.users.FindByID(ctx, id)
}

func (uc *UserUseCase) UpdateMe(ctx context.Context, id, email, username string) (*entity.User, error) {
	u, err := uc.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if email != "" {
		u.Email = email
	}
	if username != "" {
		u.Username = username
	}
	if err := uc.users.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

// DeleteMe implements GDPR right to erasure (cascade delete in the repo).
func (uc *UserUseCase) DeleteMe(ctx context.Context, id string) error {
	return uc.users.Delete(ctx, id)
}

// PersonalDataExport implements GDPR right of access.
type PersonalDataExport struct {
	User      *entity.User       `json:"user"`
	Playlists []entity.Playlist  `json:"playlists"`
	Streams   []entity.Stream    `json:"streams"`
}

func (uc *UserUseCase) ExportData(ctx context.Context, id string) (*PersonalDataExport, error) {
	u, err := uc.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Sanitize: never expose the password hash, even in export.
	u.Password = ""

	playlists, err := uc.playlists.ListByOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("playlists: %w", err)
	}
	streams, err := uc.streams.ListByBroadcaster(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("streams: %w", err)
	}
	return &PersonalDataExport{User: u, Playlists: playlists, Streams: streams}, nil
}

// AdminListUsers paginates users (admin only — authorization enforced at handler/middleware level).
func (uc *UserUseCase) AdminListUsers(ctx context.Context, offset, limit int) ([]entity.User, error) {
	return uc.users.List(ctx, offset, limit)
}

var ErrInvalidRole = errors.New("invalid role")

func (uc *UserUseCase) AdminUpdateRole(ctx context.Context, id, newRole string) (*entity.User, error) {
	switch entity.Role(newRole) {
	case entity.RoleUser, entity.RoleBroadcaster, entity.RoleAdmin:
		// ok
	default:
		return nil, ErrInvalidRole
	}
	u, err := uc.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.Role = entity.Role(newRole)
	if err := uc.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
