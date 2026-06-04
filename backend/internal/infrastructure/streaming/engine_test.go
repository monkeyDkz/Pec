package streaming

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeStreamRepo is an in-memory StreamRepository for unit tests.
type fakeStreamRepo struct {
	mu      sync.Mutex
	streams map[string]*entity.Stream
}

func newFakeStreamRepo() *fakeStreamRepo {
	return &fakeStreamRepo{streams: make(map[string]*entity.Stream)}
}

func (r *fakeStreamRepo) Create(_ context.Context, s *entity.Stream) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	r.streams[s.ID] = s
	return nil
}

func (r *fakeStreamRepo) FindByID(_ context.Context, id string) (*entity.Stream, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.streams[id]; ok {
		return s, nil
	}
	return nil, errors.New("not found")
}

func (r *fakeStreamRepo) Update(_ context.Context, s *entity.Stream) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streams[s.ID] = s
	return nil
}

func (r *fakeStreamRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.streams, id)
	return nil
}

func (r *fakeStreamRepo) ListLive(_ context.Context) ([]entity.Stream, error) { return nil, nil }
func (r *fakeStreamRepo) ListByBroadcaster(_ context.Context, _ string) ([]entity.Stream, error) {
	return nil, nil
}

func TestEngine_FanOutToMultipleListeners(t *testing.T) {
	ctx := context.Background()
	e := NewEngine(newFakeStreamRepo())

	stream, err := e.StartStream(ctx, "broadcaster-1", "My Show", "live jazz")
	require.NoError(t, err)
	require.Equal(t, entity.StreamStatusLive, stream.Status)

	// Two listeners subscribe.
	r1, err := e.ListenStream(ctx, stream.ID)
	require.NoError(t, err)
	defer r1.Close()
	r2, err := e.ListenStream(ctx, stream.ID)
	require.NoError(t, err)
	defer r2.Close()

	count, _ := e.GetListenerCount(ctx, stream.ID)
	assert.Equal(t, 2, count)

	// Broadcaster publishes audio.
	go func() {
		_ = e.PublishAudio(ctx, stream.ID, bytes.NewReader([]byte("hello world")))
	}()

	// Both listeners receive the same bytes.
	for _, r := range []io.Reader{r1, r2} {
		got := make([]byte, len("hello world"))
		n, err := io.ReadFull(r, got)
		require.NoError(t, err)
		assert.Equal(t, "hello world", string(got[:n]))
	}
}

func TestEngine_ListenerCountDropsOnClose(t *testing.T) {
	ctx := context.Background()
	e := NewEngine(newFakeStreamRepo())
	stream, _ := e.StartStream(ctx, "b1", "t", "")

	r, err := e.ListenStream(ctx, stream.ID)
	require.NoError(t, err)
	count, _ := e.GetListenerCount(ctx, stream.ID)
	assert.Equal(t, 1, count)

	require.NoError(t, r.Close())
	count, _ = e.GetListenerCount(ctx, stream.ID)
	assert.Equal(t, 0, count)
}

func TestEngine_ListenNotLive(t *testing.T) {
	e := NewEngine(newFakeStreamRepo())
	_, err := e.ListenStream(context.Background(), "does-not-exist")
	assert.ErrorIs(t, err, service.ErrStreamNotLive)
}

func TestEngine_StopRequiresOwner(t *testing.T) {
	ctx := context.Background()
	e := NewEngine(newFakeStreamRepo())
	stream, _ := e.StartStream(ctx, "owner-1", "t", "")

	err := e.StopStream(ctx, stream.ID, "someone-else")
	assert.ErrorIs(t, err, service.ErrForbidden)

	require.NoError(t, e.StopStream(ctx, stream.ID, "owner-1"))

	// After stop the bus is gone: listening is no longer possible.
	_, err = e.ListenStream(ctx, stream.ID)
	assert.ErrorIs(t, err, service.ErrStreamNotLive)
}

func TestEngine_ContextCancellationReleasesListener(t *testing.T) {
	ctx := context.Background()
	e := NewEngine(newFakeStreamRepo())
	stream, _ := e.StartStream(ctx, "b1", "t", "")

	listenCtx, cancel := context.WithCancel(ctx)
	r, err := e.ListenStream(listenCtx, stream.ID)
	require.NoError(t, err)
	defer r.Close()

	cancel() // simulate client disconnect

	done := make(chan struct{})
	go func() {
		_, _ = r.Read(make([]byte, 16)) // should return promptly with EOF
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Read did not unblock on context cancellation (leak)")
	}
}
