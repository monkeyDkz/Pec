package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/domain/service"
)

// StreamUseCase orchestrates the live streaming engine and stream persistence.
type StreamUseCase struct {
	engine service.StreamingService
	repo   repository.StreamRepository
}

func NewStreamUseCase(engine service.StreamingService, repo repository.StreamRepository) *StreamUseCase {
	return &StreamUseCase{engine: engine, repo: repo}
}

func (uc *StreamUseCase) ListLive(ctx context.Context) ([]entity.Stream, error) {
	streams, err := uc.repo.ListLive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list live: %w", err)
	}
	for i := range streams {
		count, _ := uc.engine.GetListenerCount(ctx, streams[i].ID)
		streams[i].ListenerCount = count
	}
	return streams, nil
}

func (uc *StreamUseCase) Get(ctx context.Context, id string) (*entity.Stream, error) {
	stream, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get stream: %w", err)
	}
	count, _ := uc.engine.GetListenerCount(ctx, id)
	stream.ListenerCount = count
	return stream, nil
}

func (uc *StreamUseCase) Start(ctx context.Context, broadcasterID, title, description string) (*entity.Stream, error) {
	return uc.engine.StartStream(ctx, broadcasterID, title, description)
}

func (uc *StreamUseCase) Stop(ctx context.Context, id, broadcasterID string) error {
	return uc.engine.StopStream(ctx, id, broadcasterID)
}

func (uc *StreamUseCase) Listen(ctx context.Context, id string) (io.ReadCloser, error) {
	return uc.engine.ListenStream(ctx, id)
}

func (uc *StreamUseCase) Publish(ctx context.Context, id, broadcasterID string, audio io.Reader) error {
	stream, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find stream: %w", err)
	}
	if stream.BroadcasterID != broadcasterID {
		return service.ErrForbidden
	}
	return uc.engine.PublishAudio(ctx, id, audio)
}

func (uc *StreamUseCase) Delete(ctx context.Context, id, userID, role string) error {
	stream, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find stream: %w", err)
	}
	if stream.BroadcasterID != userID && role != string(entity.RoleAdmin) {
		return service.ErrForbidden
	}
	return uc.repo.Delete(ctx, id)
}
