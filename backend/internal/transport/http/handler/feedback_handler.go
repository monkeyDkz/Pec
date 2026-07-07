package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/usecase"
)

type FeedbackHandler struct {
	uc *usecase.FeedbackUseCase
}

func NewFeedbackHandler(uc *usecase.FeedbackUseCase) *FeedbackHandler {
	return &FeedbackHandler{uc: uc}
}

type submitFeedbackRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=2000"`
}

func (h *FeedbackHandler) Submit(c *gin.Context) {
	var req submitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	f, err := h.uc.Submit(c.Request.Context(), userID.(string), req.Rating, req.Comment)
	switch {
	case errors.Is(err, usecase.ErrInvalidRating):
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be 1..5"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "submit failed"})
	default:
		c.JSON(http.StatusCreated, f)
	}
}

// AdminList returns the paginated feedback list. Authorization is enforced by
// the RequireRole("admin") middleware on the route group.
func (h *FeedbackHandler) AdminList(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	list, err := h.uc.AdminList(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, list)
}
