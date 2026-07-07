package dto

import "github.com/streampulse/backend/internal/domain/entity"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func UserFrom(u *entity.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Role:     string(u.Role),
	}
}
