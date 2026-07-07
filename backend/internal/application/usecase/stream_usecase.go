package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/streaming"
)

var (
	ErrForbidden = errors.New("forbidden")
	ErrNotLive   = errors.New("stream is not live")
)

// StreamUseCase orchestrates CRUD on streams and connects to the in-memory
// streaming registry for live publish/listen operations.
type StreamUseCase struct {
	repo     repository.StreamRepository
	registry *streaming.Registry
}

func NewStreamUseCase(repo repository.StreamRepository, registry *streaming.Registry) *StreamUseCase {
	return &StreamUseCase{repo: repo, registry: registry}
}

func (uc *StreamUseCase) Create(ctx context.Context, broadcasterID, title, description string) (*entity.Stream, error) {
	s := &entity.Stream{
		Title:         title,
		Description:   description,
		BroadcasterID: broadcasterID,
		Status:        entity.StreamStatusOffline,
	}
	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return s, nil
}

func (uc *StreamUseCase) Get(ctx context.Context, id string) (*entity.Stream, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if hub, err := uc.registry.Get(id); err == nil {
		s.ListenerCount = hub.ListenerCount()
	}
	return s, nil
}

func (uc *StreamUseCase) ListLive(ctx context.Context) ([]entity.Stream, error) {
	streams, err := uc.repo.ListLive(ctx)
	if err != nil {
		return nil, err
	}
	for i := range streams {
		if hub, err := uc.registry.Get(streams[i].ID); err == nil {
			streams[i].ListenerCount = hub.ListenerCount()
		}
	}
	return streams, nil
}

func (uc *StreamUseCase) Delete(ctx context.Context, streamID, callerID, callerRole string) error {
	s, err := uc.repo.FindByID(ctx, streamID)
	if err != nil {
		return err
	}
	if s.BroadcasterID != callerID && callerRole != string(entity.RoleAdmin) {
		return ErrForbidden
	}
	uc.registry.CloseStream(streamID)
	return uc.repo.Delete(ctx, streamID)
}

// StartLive marks a stream as live and opens the streaming hub. Only the
// stream owner (or an admin) may start broadcasting.
func (uc *StreamUseCase) StartLive(ctx context.Context, streamID, callerID, callerRole string) (*streaming.Hub, error) {
	s, err := uc.repo.FindByID(ctx, streamID)
	if err != nil {
		return nil, err
	}
	if s.BroadcasterID != callerID && callerRole != string(entity.RoleAdmin) {
		return nil, ErrForbidden
	}
	s.Status = entity.StreamStatusLive
	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}
	return uc.registry.OpenStream(ctx, streamID), nil
}

// StopLive marks a stream as offline and closes its hub.
func (uc *StreamUseCase) StopLive(ctx context.Context, streamID, callerID, callerRole string) error {
	s, err := uc.repo.FindByID(ctx, streamID)
	if err != nil {
		return err
	}
	if s.BroadcasterID != callerID && callerRole != string(entity.RoleAdmin) {
		return ErrForbidden
	}
	uc.registry.CloseStream(streamID)
	s.Status = entity.StreamStatusOffline
	return uc.repo.Update(ctx, s)
}

// LiveHub returns the hub for a currently-live stream, or ErrNotLive.
func (uc *StreamUseCase) LiveHub(streamID string) (*streaming.Hub, error) {
	hub, err := uc.registry.Get(streamID)
	if err != nil {
		return nil, ErrNotLive
	}
	return hub, nil
}
