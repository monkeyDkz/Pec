package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/repository"
)

type AdminHandler struct {
	uc    *usecase.UserUseCase
	stats *usecase.StatsUseCase
}

func NewAdminHandler(uc *usecase.UserUseCase, stats *usecase.StatsUseCase) *AdminHandler {
	return &AdminHandler{uc: uc, stats: stats}
}

func (h *AdminHandler) Stats(c *gin.Context) {
	s, err := h.stats.Snapshot(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "stats failed"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	users, err := h.uc.AdminListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	out := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		out = append(out, dto.UserFrom(&users[i]))
	}
	c.JSON(http.StatusOK, out)
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user broadcaster admin"`
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.uc.AdminUpdateRole(c.Request.Context(), c.Param("id"), req.Role)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, usecase.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
	default:
		c.JSON(http.StatusOK, dto.UserFrom(u))
	}
}
