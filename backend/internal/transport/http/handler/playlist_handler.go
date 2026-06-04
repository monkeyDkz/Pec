package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
)

type PlaylistHandler struct {
	uc *usecase.PlaylistUseCase
}

func NewPlaylistHandler(uc *usecase.PlaylistUseCase) *PlaylistHandler {
	return &PlaylistHandler{uc: uc}
}

// List GET /api/v1/playlists — current user's playlists.
func (h *PlaylistHandler) List(c *gin.Context) {
	playlists, err := h.uc.ListByOwner(c.Request.Context(), currentUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	resp := make([]dto.PlaylistResponse, 0, len(playlists))
	for i := range playlists {
		resp = append(resp, dto.NewPlaylistResponse(&playlists[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// Get GET /api/v1/playlists/:id
func (h *PlaylistHandler) Get(c *gin.Context) {
	playlist, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPlaylistResponse(playlist))
}

// Create POST /api/v1/playlists
func (h *PlaylistHandler) Create(c *gin.Context) {
	var req dto.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	playlist, err := h.uc.Create(c.Request.Context(), currentUserID(c), req.Name, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewPlaylistResponse(playlist))
}

// Update PUT /api/v1/playlists/:id
func (h *PlaylistHandler) Update(c *gin.Context) {
	var req dto.UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	playlist, err := h.uc.Update(c.Request.Context(), c.Param("id"), currentUserID(c), req.Name, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPlaylistResponse(playlist))
}

// Delete DELETE /api/v1/playlists/:id
func (h *PlaylistHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id"), currentUserID(c)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// AddTrack POST /api/v1/playlists/:id/tracks
func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	var req dto.AddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.AddTrack(c.Request.Context(), c.Param("id"), currentUserID(c), req.TrackID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "track added"})
}

// RemoveTrack DELETE /api/v1/playlists/:id/tracks/:trackId
func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {
	if err := h.uc.RemoveTrack(c.Request.Context(), c.Param("id"), currentUserID(c), c.Param("trackId")); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "track removed"})
}
