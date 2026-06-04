package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// Me GET /api/v1/users/me
func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.uc.GetByID(c.Request.Context(), currentUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// UpdateMe PUT /api/v1/users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.uc.UpdateUsername(c.Request.Context(), currentUserID(c), req.Username)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
