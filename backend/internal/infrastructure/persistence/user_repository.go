package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ repository.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u *entity.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *entity.User) error {
	if err := r.db.WithContext(ctx).Save(u).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete removes a user and cascades to their streams, playlists, uploaded
// tracks and feedbacks. GDPR right to erasure (see docs/RGPD.md).
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("broadcaster_id = ?", id).Delete(&entity.Stream{}).Error; err != nil {
			return fmt.Errorf("delete user streams: %w", err)
		}
		if err := tx.Where("owner_id = ?", id).Delete(&entity.Playlist{}).Error; err != nil {
			return fmt.Errorf("delete user playlists: %w", err)
		}
		if err := tx.Where("upload_by = ?", id).Delete(&entity.Track{}).Error; err != nil {
			return fmt.Errorf("delete user tracks: %w", err)
		}
		if err := tx.Where("user_id = ?", id).Delete(&entity.Feedback{}).Error; err != nil {
			return fmt.Errorf("delete user feedback: %w", err)
		}
		if err := tx.Delete(&entity.User{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("delete user: %w", err)
		}
		return nil
	})
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]entity.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var users []entity.User
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
