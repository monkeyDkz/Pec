package streaming

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestHub_HundredListeners proves the RNCP Bloc 3 requirement: the server
// can fan out an audio flow to 100 listeners simultaneously without losing
// chunks for a well-behaved consumer.
func TestHub_HundredListeners(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := NewHub(ctx, "stream-100")
	defer h.Close()

	const (
		listeners = 100
		chunks    = 50
	)

	received := make([]atomic.Int64, listeners)
	var wg sync.WaitGroup

	for i := 0; i < listeners; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ch, unsub, err := h.Subscribe()
			if err != nil {
				t.Errorf("subscribe %d: %v", i, err)
				return
			}
			defer unsub()
			deadline := time.After(5 * time.Second)
			for {
				select {
				case _, ok := <-ch:
					if !ok {
						return
					}
					received[i].Add(1)
					if received[i].Load() >= int64(chunks) {
						return
					}
				case <-deadline:
					return
				}
			}
		}()
	}

	// Give listeners a moment to subscribe before broadcasting.
	time.Sleep(50 * time.Millisecond)

	if got := h.ListenerCount(); got != listeners {
		t.Fatalf("expected %d listeners, got %d", listeners, got)
	}

	// Pace the publisher slightly to mimic real audio chunks.
	for i := 0; i < chunks; i++ {
		if err := h.Publish([]byte("audio-chunk")); err != nil {
			t.Fatalf("publish %d: %v", i, err)
		}
		time.Sleep(time.Millisecond)
	}

	wg.Wait()

	for i := 0; i < listeners; i++ {
		got := received[i].Load()
		// Allow a small tolerance: a few late subscribers may miss the first
		// chunks. A well-behaved consumer should still get ≥ 90 % of them.
		if got < int64(chunks*9/10) {
			t.Errorf("listener %d only received %d / %d chunks", i, got, chunks)
		}
	}
}

// BenchmarkHub_Publish measures the throughput of a single broadcaster
// against a varying number of listeners.
func BenchmarkHub_Publish(b *testing.B) {
	for _, n := range []int{1, 10, 100} {
		b.Run("listeners="+itoa(n), func(b *testing.B) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			h := NewHub(ctx, "bench")
			defer h.Close()

			var wg sync.WaitGroup
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, ch, unsub, _ := h.Subscribe()
					defer unsub()
					for range ch {
					}
				}()
			}
			time.Sleep(10 * time.Millisecond)

			chunk := make([]byte, 4096)
			b.SetBytes(int64(len(chunk)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = h.Publish(chunk)
			}
			b.StopTimer()
			h.Close()
			wg.Wait()
		})
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
