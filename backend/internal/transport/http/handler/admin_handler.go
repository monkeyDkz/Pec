package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
)

type AdminHandler struct {
	uc *usecase.AdminUseCase
}

func NewAdminHandler(uc *usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{uc: uc}
}

// ListUsers GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	users, err := h.uc.ListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, dto.NewUserResponse(&users[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateRole PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.uc.UpdateRole(c.Request.Context(), c.Param("id"), req.Role)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// Stats GET /api/v1/admin/stats
func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := h.uc.Stats(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.AdminStatsResponse{
		TotalUsers:    stats.TotalUsers,
		TotalStreams:  stats.TotalStreams,
		LiveStreams:   stats.LiveStreams,
		TotalTracks:   stats.TotalTracks,
		TotalPlaylist: stats.TotalPlaylists,
	})
}
