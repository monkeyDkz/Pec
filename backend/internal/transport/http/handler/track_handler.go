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

type TrackHandler struct {
	uc *usecase.TrackUseCase
}

func NewTrackHandler(uc *usecase.TrackUseCase) *TrackHandler {
	return &TrackHandler{uc: uc}
}

// Upload accepts a multipart/form-data POST with fields:
//   - file:   the audio file (required, ≤ 50 MB)
//   - title:  the track title (required)
//   - artist: optional artist name
func (h *TrackHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, usecase.MaxUploadBytes+1024)

	if err := c.Request.ParseMultipartForm(usecase.MaxUploadBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload: " + err.Error()})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	artist := c.PostForm("artist")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	userID, _ := c.Get("user_id")
	t, err := h.uc.Upload(c.Request.Context(), userID.(string), title, artist, contentType, file)
	switch {
	case errors.Is(err, usecase.ErrUnsupportedMedia):
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "audio format not supported"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed: " + err.Error()})
	default:
		c.JSON(http.StatusCreated, dto.TrackFrom(t))
	}
}

func (h *TrackHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	tracks, err := h.uc.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	out := make([]dto.TrackResponse, 0, len(tracks))
	for i := range tracks {
		out = append(out, dto.TrackFrom(&tracks[i]))
	}
	c.JSON(http.StatusOK, out)
}

func (h *TrackHandler) Get(c *gin.Context) {
	t, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}
	c.JSON(http.StatusOK, dto.TrackFrom(t))
}

func (h *TrackHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	err := h.uc.Delete(c.Request.Context(), c.Param("id"), userID.(string), role.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your track"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}
