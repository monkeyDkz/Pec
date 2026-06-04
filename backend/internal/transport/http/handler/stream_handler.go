package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
)

type StreamHandler struct {
	uc *usecase.StreamUseCase
}

func NewStreamHandler(uc *usecase.StreamUseCase) *StreamHandler {
	return &StreamHandler{uc: uc}
}

// List GET /api/v1/streams — live streams.
func (h *StreamHandler) List(c *gin.Context) {
	streams, err := h.uc.ListLive(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	resp := make([]dto.StreamResponse, 0, len(streams))
	for i := range streams {
		resp = append(resp, dto.NewStreamResponse(&streams[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// Get GET /api/v1/streams/:id
func (h *StreamHandler) Get(c *gin.Context) {
	stream, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewStreamResponse(stream))
}

// Create POST /api/v1/streams — broadcaster starts a live stream.
func (h *StreamHandler) Create(c *gin.Context) {
	var req dto.CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stream, err := h.uc.Start(c.Request.Context(), currentUserID(c), req.Title, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewStreamResponse(stream))
}

// Stop POST /api/v1/streams/:id/stop
func (h *StreamHandler) Stop(c *gin.Context) {
	if err := h.uc.Stop(c.Request.Context(), c.Param("id"), currentUserID(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "offline"})
}

// Delete DELETE /api/v1/streams/:id
func (h *StreamHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id"), currentUserID(c), currentUserRole(c)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Listen GET /api/v1/streams/:id/listen — streams audio chunks to the client.
func (h *StreamHandler) Listen(c *gin.Context) {
	reader, err := h.uc.Listen(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", "audio/mpeg")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	buf := make([]byte, 4096)
	c.Stream(func(w io.Writer) bool {
		n, rerr := reader.Read(buf)
		if n > 0 {
			if _, werr := w.Write(buf[:n]); werr != nil {
				return false
			}
		}
		return rerr == nil
	})
}

// Publish POST /api/v1/streams/:id/publish — broadcaster pushes audio bytes.
func (h *StreamHandler) Publish(c *gin.Context) {
	err := h.uc.Publish(c.Request.Context(), c.Param("id"), currentUserID(c), c.Request.Body)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stream ended"})
}
