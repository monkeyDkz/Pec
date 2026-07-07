package usecase

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/streaming"
)

// fakeStreamRepo is an in-memory implementation of the StreamRepository
// interface, sufficient for usecase-level testing.
type fakeStreamRepo struct {
	mu      sync.Mutex
	byID    map[string]*entity.Stream
	counter int
}

func newFakeStreamRepo() *fakeStreamRepo {
	return &fakeStreamRepo{byID: map[string]*entity.Stream{}}
}

func (r *fakeStreamRepo) Create(_ context.Context, s *entity.Stream) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	if s.ID == "" {
		s.ID = string(rune('a'-1+r.counter)) + "-id"
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	r.byID[s.ID] = s
	return nil
}
func (r *fakeStreamRepo) FindByID(_ context.Context, id string) (*entity.Stream, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.byID[id]; ok {
		return s, nil
	}
	return nil, repository.ErrNotFound
}
func (r *fakeStreamRepo) Update(_ context.Context, s *entity.Stream) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[s.ID] = s
	return nil
}
func (r *fakeStreamRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byID, id)
	return nil
}
func (r *fakeStreamRepo) ListLive(_ context.Context) ([]entity.Stream, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []entity.Stream{}
	for _, s := range r.byID {
		if s.Status == entity.StreamStatusLive {
			out = append(out, *s)
		}
	}
	return out, nil
}
func (r *fakeStreamRepo) ListByBroadcaster(_ context.Context, bid string) ([]entity.Stream, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []entity.Stream{}
	for _, s := range r.byID {
		if s.BroadcasterID == bid {
			out = append(out, *s)
		}
	}
	return out, nil
}

func TestStreamUseCase_CreateAndGet(t *testing.T) {
	repo := newFakeStreamRepo()
	uc := NewStreamUseCase(repo, streaming.NewRegistry())

	created, err := uc.Create(context.Background(), "bcaster-1", "My Stream", "desc")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" || created.Status != entity.StreamStatusOffline {
		t.Fatalf("unexpected created stream: %+v", created)
	}
	got, err := uc.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "My Stream" {
		t.Fatalf("unexpected title: %s", got.Title)
	}
}

func TestStreamUseCase_StartLiveOnlyOwner(t *testing.T) {
	repo := newFakeStreamRepo()
	registry := streaming.NewRegistry()
	uc := NewStreamUseCase(repo, registry)

	s, _ := uc.Create(context.Background(), "owner", "t", "")

	// Foreign user can't start the stream
	_, err := uc.StartLive(context.Background(), s.ID, "stranger", "user")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	hub, err := uc.StartLive(context.Background(), s.ID, "owner", "user")
	if err != nil {
		t.Fatalf("owner start: %v", err)
	}
	if hub == nil {
		t.Fatal("expected hub")
	}
	if !registry.IsLive(s.ID) {
		t.Fatal("expected stream live after StartLive")
	}
}

func TestStreamUseCase_AdminCanStartAnyStream(t *testing.T) {
	repo := newFakeStreamRepo()
	uc := NewStreamUseCase(repo, streaming.NewRegistry())
	s, _ := uc.Create(context.Background(), "owner", "t", "")
	if _, err := uc.StartLive(context.Background(), s.ID, "admin-id", "admin"); err != nil {
		t.Fatalf("admin start: %v", err)
	}
}

func TestStreamUseCase_StopLiveResetsState(t *testing.T) {
	repo := newFakeStreamRepo()
	registry := streaming.NewRegistry()
	uc := NewStreamUseCase(repo, registry)

	s, _ := uc.Create(context.Background(), "owner", "t", "")
	_, _ = uc.StartLive(context.Background(), s.ID, "owner", "user")
	if err := uc.StopLive(context.Background(), s.ID, "owner", "user"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if registry.IsLive(s.ID) {
		t.Fatal("expected stream offline after StopLive")
	}
	final, _ := uc.Get(context.Background(), s.ID)
	if final.Status != entity.StreamStatusOffline {
		t.Fatalf("expected offline status, got %s", final.Status)
	}
}

func TestStreamUseCase_LiveHubReturnsErrNotLiveWhenOffline(t *testing.T) {
	repo := newFakeStreamRepo()
	uc := NewStreamUseCase(repo, streaming.NewRegistry())
	if _, err := uc.LiveHub("nope"); err != ErrNotLive {
		t.Fatalf("expected ErrNotLive, got %v", err)
	}
}

func TestStreamUseCase_ListLiveAddsListenerCount(t *testing.T) {
	repo := newFakeStreamRepo()
	registry := streaming.NewRegistry()
	uc := NewStreamUseCase(repo, registry)

	s, _ := uc.Create(context.Background(), "owner", "title", "")
	hub, _ := uc.StartLive(context.Background(), s.ID, "owner", "user")
	_, _, _, _ = hub.Subscribe()
	_, _, _, _ = hub.Subscribe()

	live, err := uc.ListLive(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(live) != 1 || live[0].ListenerCount != 2 {
		t.Fatalf("unexpected list: %+v", live)
	}
}

func TestStreamUseCase_DeleteOnlyOwnerOrAdmin(t *testing.T) {
	repo := newFakeStreamRepo()
	uc := NewStreamUseCase(repo, streaming.NewRegistry())

	s, _ := uc.Create(context.Background(), "owner", "t", "")
	if err := uc.Delete(context.Background(), s.ID, "stranger", "user"); err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if err := uc.Delete(context.Background(), s.ID, "owner", "user"); err != nil {
		t.Fatalf("owner delete: %v", err)
	}
}
