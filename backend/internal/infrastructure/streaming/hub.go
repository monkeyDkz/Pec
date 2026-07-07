// Package streaming implements the in-memory pub/sub hub that broadcasts
// audio chunks from a single broadcaster to N listeners.
//
// Design rationale: see docs/adr/0008-streaming-pubsub.md.
//
// Concurrency model:
//   - one goroutine per listener (the listener's HTTP handler reads its channel);
//   - the broadcaster writes to a single channel that the hub fans out to listeners;
//   - backpressure is handled per-listener: a slow listener loses chunks but never
//     blocks the broadcaster or its peers;
//   - context cancellation propagates from broadcaster to all listeners,
//     guaranteeing goroutine cleanup (no leaks).
package streaming

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/infrastructure/observability"
)

// DefaultListenerBuffer is the buffer size of a listener's chunk channel.
// 256 chunks ≈ a few seconds of audio at typical bitrates — enough to absorb
// transient network jitter without unbounded memory growth.
const DefaultListenerBuffer = 256

var (
	// ErrStreamNotFound is returned when a stream id has no active hub.
	ErrStreamNotFound = errors.New("stream not found")
	// ErrStreamClosed is returned when publishing to a closed hub.
	ErrStreamClosed = errors.New("stream closed")
)

// Hub multiplexes audio chunks from one publisher to many subscribers.
// It is safe for concurrent use.
type Hub struct {
	streamID string

	mu          sync.RWMutex
	subscribers map[string]chan []byte
	closed      bool

	bufSize int

	// listenerCount mirrors len(subscribers) atomically so reads do not
	// require the lock and metrics stay fast.
	listenerCount atomic.Int64

	ctx    context.Context
	cancel context.CancelFunc
}

// NewHub creates a hub for the given stream id.
// The provided context controls the lifetime of the hub; canceling it closes
// all subscribers and prevents further publishes.
func NewHub(parent context.Context, streamID string) *Hub {
	ctx, cancel := context.WithCancel(parent)
	return &Hub{
		streamID:    streamID,
		subscribers: make(map[string]chan []byte),
		bufSize:     DefaultListenerBuffer,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Subscribe registers a new listener and returns its receive-only channel and
// an unsubscribe function. The returned channel is closed when the listener
// is removed or when the hub is closed.
func (h *Hub) Subscribe() (id string, ch <-chan []byte, unsubscribe func(), err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.closed {
		return "", nil, nil, ErrStreamClosed
	}

	id = uuid.NewString()
	out := make(chan []byte, h.bufSize)
	h.subscribers[id] = out
	h.listenerCount.Add(1)
	observability.ActiveListeners.WithLabelValues(h.streamID).Set(float64(h.listenerCount.Load()))

	unsub := func() {
		h.unsubscribe(id, "client_close")
	}

	return id, out, unsub, nil
}

// Publish broadcasts a chunk to every subscriber. Slow subscribers see their
// chunk dropped (and a disconnection counter incremented) rather than slowing
// the broadcaster.
func (h *Hub) Publish(chunk []byte) error {
	h.mu.RLock()
	if h.closed {
		h.mu.RUnlock()
		return ErrStreamClosed
	}

	// Copy the chunk so a slow consumer cannot retain a reference into the
	// broadcaster's underlying buffer if it later reuses it.
	buf := make([]byte, len(chunk))
	copy(buf, chunk)

	for id, ch := range h.subscribers {
		select {
		case ch <- buf:
			// delivered
		default:
			// listener is too slow; drop and account for it as a disconnection.
			observability.StreamDisconnections.WithLabelValues(h.streamID, "slow_consumer").Inc()
			// Detach the slow listener to free resources. We must drop the read
			// lock to take the write lock used by unsubscribe.
			go h.unsubscribe(id, "slow_consumer")
		}
	}
	h.mu.RUnlock()

	observability.StreamBytesTotal.WithLabelValues(h.streamID).Add(float64(len(chunk)))
	return nil
}

// Close terminates the hub, closes every subscriber channel and cancels the
// context. Safe to call multiple times.
func (h *Hub) Close() {
	h.mu.Lock()
	if h.closed {
		h.mu.Unlock()
		return
	}
	h.closed = true
	for id, ch := range h.subscribers {
		close(ch)
		delete(h.subscribers, id)
	}
	h.listenerCount.Store(0)
	observability.ActiveListeners.DeleteLabelValues(h.streamID)
	h.mu.Unlock()
	h.cancel()
}

// ListenerCount returns the current number of subscribers.
func (h *Hub) ListenerCount() int {
	return int(h.listenerCount.Load())
}

// Done returns the hub's cancellation channel for downstream goroutines.
func (h *Hub) Done() <-chan struct{} {
	return h.ctx.Done()
}

func (h *Hub) unsubscribe(id, reason string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch, ok := h.subscribers[id]
	if !ok {
		return
	}
	close(ch)
	delete(h.subscribers, id)
	h.listenerCount.Add(-1)
	observability.ActiveListeners.WithLabelValues(h.streamID).Set(float64(h.listenerCount.Load()))
	if reason != "client_close" {
		observability.StreamDisconnections.WithLabelValues(h.streamID, reason).Inc()
	}
}

// Registry tracks all active hubs by stream id. It is the single entry point
// for handlers to publish or subscribe.
type Registry struct {
	mu   sync.RWMutex
	hubs map[string]*Hub
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{hubs: make(map[string]*Hub)}
}

// OpenStream creates a new hub for a stream id, replacing any pre-existing one.
func (r *Registry) OpenStream(ctx context.Context, streamID string) *Hub {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.hubs[streamID]; ok {
		// Replacing a live hub: close it and account for its removal so the
		// gauge stays balanced (otherwise re-publishing the same stream leaks +1).
		existing.Close()
		observability.ActiveStreams.Dec()
	}
	h := NewHub(ctx, streamID)
	r.hubs[streamID] = h
	observability.ActiveStreams.Inc()
	return h
}

// CloseStream removes and closes a stream's hub.
func (r *Registry) CloseStream(streamID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	h, ok := r.hubs[streamID]
	if !ok {
		return
	}
	h.Close()
	delete(r.hubs, streamID)
	observability.ActiveStreams.Dec()
}

// Get returns the hub for a stream id, or ErrStreamNotFound.
func (r *Registry) Get(streamID string) (*Hub, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, ok := r.hubs[streamID]
	if !ok {
		return nil, ErrStreamNotFound
	}
	return h, nil
}

// IsLive reports whether a stream has an active hub.
func (r *Registry) IsLive(streamID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.hubs[streamID]
	return ok
}

// Snapshot returns the list of currently live stream ids.
func (r *Registry) Snapshot() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.hubs))
	for id := range r.hubs {
		ids = append(ids, id)
	}
	return ids
}
