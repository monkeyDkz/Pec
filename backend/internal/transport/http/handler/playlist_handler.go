package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/repository"
)

type PlaylistHandler struct {
	uc *usecase.PlaylistUseCase
}

func NewPlaylistHandler(uc *usecase.PlaylistUseCase) *PlaylistHandler {
	return &PlaylistHandler{uc: uc}
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	var req dto.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	p, err := h.uc.Create(c.Request.Context(), userID.(string), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, dto.PlaylistFrom(p))
}

func (h *PlaylistHandler) ListMine(c *gin.Context) {
	userID, _ := c.Get("user_id")
	playlists, err := h.uc.ListByOwner(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	out := make([]dto.PlaylistResponse, 0, len(playlists))
	for i := range playlists {
		out = append(out, dto.PlaylistFrom(&playlists[i]))
	}
	c.JSON(http.StatusOK, out)
}

func (h *PlaylistHandler) Get(c *gin.Context) {
	p, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}
	c.JSON(http.StatusOK, dto.PlaylistFrom(p))
}

func (h *PlaylistHandler) Update(c *gin.Context) {
	var req dto.UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	p, err := h.uc.Update(c.Request.Context(), c.Param("id"), userID.(string), req.Name, req.Description)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your playlist"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
	default:
		c.JSON(http.StatusOK, dto.PlaylistFrom(p))
	}
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	err := h.uc.Delete(c.Request.Context(), c.Param("id"), userID.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your playlist"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}

func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	userID, _ := c.Get("user_id")
	err := h.uc.AddTrack(c.Request.Context(), c.Param("id"), c.Param("trackId"), userID.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist or track not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your playlist"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "add track failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}

func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {
	userID, _ := c.Get("user_id")
	err := h.uc.RemoveTrack(c.Request.Context(), c.Param("id"), c.Param("trackId"), userID.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your playlist"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "remove failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}

type reorderTracksRequest struct {
	TrackIDs []string `json:"track_ids" binding:"required"`
}

// Reorder persists a new track order for the playlist. The request body
// must list every track id in the desired final order.
func (h *PlaylistHandler) Reorder(c *gin.Context) {
	var req reorderTracksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	err := h.uc.ReorderTracks(c.Request.Context(), c.Param("id"), req.TrackIDs, userID.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your playlist"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "reorder failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}
