package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// ErrNotFound is returned when an entity does not exist.
var ErrNotFound = errors.New("not found")

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
