// Package streaming implements the real-time audio fan-out engine: a single
// broadcaster publishes audio chunks that are multiplexed to N listeners using
// goroutines + channels, with context-based cancellation to avoid leaks.
package streaming

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/domain/service"
	"github.com/streampulse/backend/internal/infrastructure/observability"
)

const (
	// subBufferSize bounds per-listener buffering; slow listeners drop chunks
	// rather than blocking the broadcaster (backpressure isolation).
	subBufferSize = 128
	// readChunkSize is how much audio we read from the publisher at a time.
	readChunkSize = 4096
)

// bus is the in-memory pub/sub channel set for a single live stream.
type bus struct {
	mu   sync.RWMutex
	subs map[string]chan []byte
}

func (b *bus) broadcast(chunk []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs {
		select {
		case ch <- chunk:
		default:
			// Listener is too slow: drop this chunk to protect the broadcaster.
		}
	}
}

// Engine implements service.StreamingService backed by a StreamRepository for
// persistence and an in-memory bus map for the live fan-out.
type Engine struct {
	streamRepo repository.StreamRepository
	mu         sync.RWMutex
	buses      map[string]*bus
}

// NewEngine builds a streaming Engine.
func NewEngine(streamRepo repository.StreamRepository) *Engine {
	return &Engine{
		streamRepo: streamRepo,
		buses:      make(map[string]*bus),
	}
}

// StartStream persists a new live stream and opens its bus.
func (e *Engine) StartStream(ctx context.Context, broadcasterID, title, description string) (*entity.Stream, error) {
	stream := &entity.Stream{
		Title:         title,
		Description:   description,
		BroadcasterID: broadcasterID,
		Status:        entity.StreamStatusLive,
	}
	if err := e.streamRepo.Create(ctx, stream); err != nil {
		return nil, fmt.Errorf("create stream: %w", err)
	}

	e.mu.Lock()
	e.buses[stream.ID] = &bus{subs: make(map[string]chan []byte)}
	e.mu.Unlock()

	observability.ActiveStreams.Inc()
	return stream, nil
}

// StopStream marks a stream offline and tears down its bus (closing listeners).
func (e *Engine) StopStream(ctx context.Context, streamID, broadcasterID string) error {
	stream, err := e.streamRepo.FindByID(ctx, streamID)
	if err != nil {
		return fmt.Errorf("find stream: %w", err)
	}
	if stream.BroadcasterID != broadcasterID {
		return service.ErrForbidden
	}

	stream.Status = entity.StreamStatusOffline
	if err := e.streamRepo.Update(ctx, stream); err != nil {
		return fmt.Errorf("update stream: %w", err)
	}

	e.closeBus(streamID)
	observability.ActiveStreams.Dec()
	return nil
}

func (e *Engine) closeBus(streamID string) {
	e.mu.Lock()
	b := e.buses[streamID]
	delete(e.buses, streamID)
	e.mu.Unlock()
	if b == nil {
		return
	}
	b.mu.Lock()
	for id, ch := range b.subs {
		close(ch)
		delete(b.subs, id)
	}
	b.mu.Unlock()
	observability.ActiveListeners.DeleteLabelValues(streamID)
}

// ListenStream subscribes a new listener and returns a reader over its channel.
// Closing the returned reader unsubscribes and releases resources.
func (e *Engine) ListenStream(ctx context.Context, streamID string) (io.ReadCloser, error) {
	e.mu.RLock()
	b := e.buses[streamID]
	e.mu.RUnlock()
	if b == nil {
		return nil, service.ErrStreamNotLive
	}

	subID := uuid.NewString()
	ch := make(chan []byte, subBufferSize)

	b.mu.Lock()
	b.subs[subID] = ch
	b.mu.Unlock()

	observability.ActiveListeners.WithLabelValues(streamID).Inc()

	return &streamReader{
		ch:  ch,
		ctx: ctx,
		closeFn: func() {
			b.mu.Lock()
			if _, ok := b.subs[subID]; ok {
				delete(b.subs, subID)
				close(ch)
			}
			b.mu.Unlock()
			observability.ActiveListeners.WithLabelValues(streamID).Dec()
			observability.StreamDisconnections.WithLabelValues(streamID, "listener_left").Inc()
		},
	}, nil
}

// PublishAudio reads from the broadcaster source and fans chunks out to listeners
// until EOF or context cancellation.
func (e *Engine) PublishAudio(ctx context.Context, streamID string, audio io.Reader) error {
	e.mu.RLock()
	b := e.buses[streamID]
	e.mu.RUnlock()
	if b == nil {
		return service.ErrStreamNotLive
	}

	buf := make([]byte, readChunkSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := audio.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			b.broadcast(chunk)
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read audio: %w", err)
		}
	}
}

// GetListenerCount returns the number of active listeners for a stream.
func (e *Engine) GetListenerCount(ctx context.Context, streamID string) (int, error) {
	e.mu.RLock()
	b := e.buses[streamID]
	e.mu.RUnlock()
	if b == nil {
		return 0, nil
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs), nil
}

// streamReader adapts a chunk channel into an io.ReadCloser. It also watches the
// request context so a disconnected listener releases its goroutine immediately
// (no leak) instead of blocking until the next audio chunk.
type streamReader struct {
	ch      chan []byte
	ctx     context.Context
	buf     []byte
	closeFn func()
	once    sync.Once
}

func (r *streamReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		select {
		case chunk, ok := <-r.ch:
			if !ok {
				return 0, io.EOF
			}
			r.buf = chunk
		case <-r.ctx.Done():
			return 0, io.EOF
		}
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

func (r *streamReader) Close() error {
	r.once.Do(r.closeFn)
	return nil
}
