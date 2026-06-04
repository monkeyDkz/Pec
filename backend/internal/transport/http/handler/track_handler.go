package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
)

type TrackHandler struct {
	uc *usecase.TrackUseCase
}

func NewTrackHandler(uc *usecase.TrackUseCase) *TrackHandler {
	return &TrackHandler{uc: uc}
}

// Create POST /api/v1/tracks
func (h *TrackHandler) Create(c *gin.Context) {
	var req dto.CreateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	track, err := h.uc.Create(c.Request.Context(), currentUserID(c), req.Title, req.Artist, req.Duration, req.FileURL)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewTrackResponse(track))
}

// Get GET /api/v1/tracks/:id
func (h *TrackHandler) Get(c *gin.Context) {
	track, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewTrackResponse(track))
}

// List GET /api/v1/tracks?offset=&limit=
func (h *TrackHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	tracks, err := h.uc.List(c.Request.Context(), offset, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := make([]dto.TrackResponse, 0, len(tracks))
	for i := range tracks {
		resp = append(resp, dto.NewTrackResponse(&tracks[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// Delete DELETE /api/v1/tracks/:id
func (h *TrackHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
