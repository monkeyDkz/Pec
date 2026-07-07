package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/repository"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Me(c *gin.Context) {
	id, _ := c.Get("user_id")
	u, err := h.uc.Me(c.Request.Context(), id.(string))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}
	c.JSON(http.StatusOK, dto.UserFrom(u))
}

type updateMeRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := c.Get("user_id")
	u, err := h.uc.UpdateMe(c.Request.Context(), id.(string), req.Email, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, dto.UserFrom(u))
}

// DeleteMe — GDPR right to erasure.
func (h *UserHandler) DeleteMe(c *gin.Context) {
	id, _ := c.Get("user_id")
	if err := h.uc.DeleteMe(c.Request.Context(), id.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ExportMyData — GDPR right of access.
func (h *UserHandler) ExportMyData(c *gin.Context) {
	id, _ := c.Get("user_id")
	export, err := h.uc.ExportData(c.Request.Context(), id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}
	c.Header("Content-Disposition", `attachment; filename="streampulse-export.json"`)
	c.JSON(http.StatusOK, export)
}
