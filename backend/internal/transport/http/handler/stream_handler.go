package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/repository"
)

const (
	publishChunkSize = 4 * 1024 // 4 KB per read from the broadcaster
)

type StreamHandler struct {
	uc *usecase.StreamUseCase
}

func NewStreamHandler(uc *usecase.StreamUseCase) *StreamHandler {
	return &StreamHandler{uc: uc}
}

// Create godoc
// @Summary      Create a new stream
// @Tags         streams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body     dto.CreateStreamRequest true "Stream metadata"
// @Success      201  {object} dto.StreamResponse
// @Router       /api/v1/streams [post]
func (h *StreamHandler) Create(c *gin.Context) {
	var req dto.CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	s, err := h.uc.Create(c.Request.Context(), userID.(string), req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, dto.StreamFrom(s))
}

// ListLive godoc
// @Summary      List live streams
// @Tags         streams
// @Produce      json
// @Success      200 {array} dto.StreamResponse
// @Router       /api/v1/streams [get]
func (h *StreamHandler) ListLive(c *gin.Context) {
	streams, err := h.uc.ListLive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	out := make([]dto.StreamResponse, 0, len(streams))
	for i := range streams {
		out = append(out, dto.StreamFrom(&streams[i]))
	}
	c.JSON(http.StatusOK, out)
}

func (h *StreamHandler) Get(c *gin.Context) {
	id := c.Param("id")
	s, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}
	c.JSON(http.StatusOK, dto.StreamFrom(s))
}

func (h *StreamHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	err := h.uc.Delete(c.Request.Context(), id, userID.(string), role.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your stream"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
	default:
		c.Status(http.StatusNoContent)
	}
}

// Publish opens a long-lived HTTP connection that reads audio chunks from the
// broadcaster and pushes them into the streaming hub.
//
// The client should send chunked octet-stream data. The stream switches to
// "live" the moment publishing starts and back to "offline" when this
// handler returns (normally or on error).
func (h *StreamHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	// Authorize and mark the stream live.
	hub, err := h.uc.StartLive(c.Request.Context(), id, userID.(string), role.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
		return
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your stream"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start"})
		return
	}

	// Always stop the stream when we exit.
	defer func() {
		stopCtx, cancel := newDetachedContext(5 * time.Second)
		defer cancel()
		if err := h.uc.StopLive(stopCtx, id, userID.(string), role.(string)); err != nil {
			slog.Error("stop live", "stream_id", id, "err", err)
		}
	}()

	c.Header("Cache-Control", "no-store")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	buf := make([]byte, publishChunkSize)
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-hub.Done():
			return
		default:
		}
		n, err := c.Request.Body.Read(buf)
		if n > 0 {
			if perr := hub.Publish(buf[:n]); perr != nil {
				slog.Error("publish", "stream_id", id, "err", perr)
				return
			}
		}
		if err == io.EOF {
			return
		}
		if err != nil {
			slog.Info("publish read error", "stream_id", id, "err", err)
			return
		}
	}
}

// PublishWS is the WebSocket equivalent of Publish, used by web clients
// (browsers cannot stream an HTTP request body). The broadcaster sends binary
// audio frames over the socket; each frame is pushed into the streaming hub.
//
// Auth is carried in the `token` query parameter (browsers cannot set an
// Authorization header on a WebSocket); the route uses AuthMiddlewareQuery.
func (h *StreamHandler) PublishWS(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	// Authorize and mark the stream live before upgrading the connection so we
	// can still return a clean HTTP error on failure.
	hub, err := h.uc.StartLive(c.Request.Context(), id, userID.(string), role.(string))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
		return
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "not your stream"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start"})
		return
	}

	// Always stop the stream when the socket closes.
	defer func() {
		stopCtx, cancel := newDetachedContext(5 * time.Second)
		defer cancel()
		if err := h.uc.StopLive(stopCtx, id, userID.(string), role.(string)); err != nil {
			slog.Error("stop live (ws)", "stream_id", id, "err", err)
		}
	}()

	// OriginPatterns allows the cross-origin dev setup (web app on :8090, API on
	// :8080). Same-origin connections are always accepted; in production the web
	// bundle is served from the API's own origin so no extra patterns are needed.
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*", "127.0.0.1:*"},
	})
	if err != nil {
		slog.Error("ws accept", "stream_id", id, "err", err)
		return
	}
	defer conn.CloseNow()

	ctx := c.Request.Context()
	for {
		select {
		case <-hub.Done():
			_ = conn.Close(websocket.StatusNormalClosure, "stream stopped")
			return
		default:
		}
		_, data, err := conn.Read(ctx)
		if err != nil {
			// Normal client disconnect or context cancellation.
			return
		}
		if len(data) > 0 {
			if perr := hub.Publish(data); perr != nil {
				slog.Error("publish (ws)", "stream_id", id, "err", perr)
				return
			}
		}
	}
}

// Listen subscribes to a live stream and streams chunks back to the HTTP client
// using chunked transfer-encoding. The connection closes when the client
// disconnects or the broadcaster stops.
func (h *StreamHandler) Listen(c *gin.Context) {
	id := c.Param("id")
	hub, err := h.uc.LiveHub(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stream is not live"})
		return
	}

	_, ch, unsub, err := hub.Subscribe()
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"error": "stream closed"})
		return
	}
	defer unsub()

	c.Header("Content-Type", "audio/mpeg")
	c.Header("Cache-Control", "no-store")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return
			}
			if _, err := c.Writer.Write(chunk); err != nil {
				return
			}
			c.Writer.Flush()
		case <-clientGone:
			return
		case <-hub.Done():
			return
		}
	}
}
