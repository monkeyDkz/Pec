package streaming

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestHub_SubscribeReceivesPublishedChunks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-1")
	defer h.Close()

	_, ch, unsub, err := h.Subscribe()
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	defer unsub()

	if err := h.Publish([]byte("hello")); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case got := <-ch:
		if string(got) != "hello" {
			t.Fatalf("expected hello, got %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for chunk")
	}
}

func TestHub_MultipleListenersReceiveSameChunks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-1")
	defer h.Close()

	const n = 10
	channels := make([]<-chan []byte, n)
	for i := 0; i < n; i++ {
		_, ch, _, err := h.Subscribe()
		if err != nil {
			t.Fatalf("subscribe %d: %v", i, err)
		}
		channels[i] = ch
	}

	if err := h.Publish([]byte("audio")); err != nil {
		t.Fatalf("publish: %v", err)
	}

	for i, ch := range channels {
		select {
		case got := <-ch:
			if string(got) != "audio" {
				t.Fatalf("listener %d: expected audio, got %q", i, got)
			}
		case <-time.After(time.Second):
			t.Fatalf("listener %d: timeout", i)
		}
	}

	if h.ListenerCount() != n {
		t.Fatalf("expected %d listeners, got %d", n, h.ListenerCount())
	}
}

func TestHub_UnsubscribeDoesNotLeak(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-1")
	defer h.Close()

	_, ch, unsub, _ := h.Subscribe()

	unsub()

	// After unsubscribe, the channel must be closed.
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected closed channel, got value")
		}
	case <-time.After(time.Second):
		t.Fatal("channel was not closed after unsubscribe")
	}

	if h.ListenerCount() != 0 {
		t.Fatalf("expected 0 listeners after unsubscribe, got %d", h.ListenerCount())
	}
}

func TestHub_CloseClosesAllSubscribers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-1")
	_, ch1, _, _ := h.Subscribe()
	_, ch2, _, _ := h.Subscribe()

	h.Close()

	for i, ch := range []<-chan []byte{ch1, ch2} {
		select {
		case _, ok := <-ch:
			if ok {
				t.Fatalf("listener %d: expected closed channel", i)
			}
		case <-time.After(time.Second):
			t.Fatalf("listener %d: channel was not closed by hub.Close", i)
		}
	}

	if got := h.ListenerCount(); got != 0 {
		t.Fatalf("expected 0 listeners after Close, got %d", got)
	}
}

func TestHub_PublishAfterCloseReturnsError(t *testing.T) {
	h := NewHub(context.Background(), "stream-1")
	h.Close()

	if err := h.Publish([]byte("x")); err != ErrStreamClosed {
		t.Fatalf("expected ErrStreamClosed, got %v", err)
	}
}

func TestHub_SubscribeAfterCloseReturnsError(t *testing.T) {
	h := NewHub(context.Background(), "stream-1")
	h.Close()

	_, _, _, err := h.Subscribe()
	if err != ErrStreamClosed {
		t.Fatalf("expected ErrStreamClosed, got %v", err)
	}
}

func TestHub_SlowListenerIsDroppedNotBlocking(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-1")
	defer h.Close()

	// Use a tiny buffer to force overflow quickly.
	h.bufSize = 2
	_, _, _, _ = h.Subscribe() // never reads; will be detached

	// We expect 50 publishes to complete in under a second, proving the slow
	// listener does not block the broadcaster.
	done := make(chan struct{})
	go func() {
		for i := 0; i < 50; i++ {
			_ = h.Publish([]byte("chunk"))
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("slow listener blocked the broadcaster")
	}
}

func TestHub_ConcurrentPublishersAndSubscribers(t *testing.T) {
	// This test runs many concurrent operations to exercise the race detector.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-race")
	defer h.Close()

	var wg sync.WaitGroup
	var received atomic.Int64

	// 20 listeners that consume continuously
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ch, unsub, err := h.Subscribe()
			if err != nil {
				return
			}
			defer unsub()
			t := time.After(200 * time.Millisecond)
			for {
				select {
				case _, ok := <-ch:
					if !ok {
						return
					}
					received.Add(1)
				case <-t:
					return
				}
			}
		}()
	}

	// 5 publishers writing concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = h.Publish([]byte("x"))
			}
		}()
	}

	wg.Wait()

	if received.Load() == 0 {
		t.Fatal("no chunk received under concurrency")
	}
}

func TestRegistry_OpenGetClose(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	h := r.OpenStream(ctx, "s1")
	if !r.IsLive("s1") {
		t.Fatal("expected s1 to be live")
	}

	got, err := r.Get("s1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != h {
		t.Fatal("registry returned a different hub instance")
	}

	r.CloseStream("s1")
	if r.IsLive("s1") {
		t.Fatal("expected s1 to be closed")
	}
	if _, err := r.Get("s1"); err != ErrStreamNotFound {
		t.Fatalf("expected ErrStreamNotFound, got %v", err)
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()
	r.OpenStream(ctx, "a")
	r.OpenStream(ctx, "b")
	r.OpenStream(ctx, "c")

	snap := r.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 streams, got %d", len(snap))
	}
}
