package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/repository"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body     dto.RegisterRequest true "Registration payload"
// @Success      201  {object} dto.AuthResponse
// @Failure      400  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.uc.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	switch {
	case errors.Is(err, repository.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "email or username already taken"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{Token: token, User: dto.UserFrom(user)})
}

// Login godoc
// @Summary      Sign in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body     dto.LoginRequest true "Credentials"
// @Success      200  {object} dto.AuthResponse
// @Failure      401  {object} map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token, User: dto.UserFrom(user)})
}

// Refresh godoc
// @Summary      Refresh JWT
// @Description  Returns a fresh JWT for the authenticated user. Slides the
// @Description  expiration window without forcing the user to re-enter credentials.
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.AuthResponse
// @Failure      401 {object} map[string]string
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, token, err := h.uc.Refresh(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh failed"})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{Token: token, User: dto.UserFrom(user)})
}
