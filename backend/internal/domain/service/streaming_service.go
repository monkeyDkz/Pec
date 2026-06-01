package service

import (
	"context"
	"io"

	"github.com/streampulse/backend/internal/domain/entity"
)

type StreamingService interface {
	StartStream(ctx context.Context, broadcasterID, title, description string) (*entity.Stream, error)
	StopStream(ctx context.Context, streamID, broadcasterID string) error
	ListenStream(ctx context.Context, streamID string) (io.ReadCloser, error)
	PublishAudio(ctx context.Context, streamID string, audio io.Reader) error
	GetListenerCount(ctx context.Context, streamID string) (int, error)
}
