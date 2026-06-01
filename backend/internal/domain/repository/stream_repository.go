package repository

import (
	"context"

	"github.com/streampulse/backend/internal/domain/entity"
)

type StreamRepository interface {
	Create(ctx context.Context, stream *entity.Stream) error
	FindByID(ctx context.Context, id string) (*entity.Stream, error)
	Update(ctx context.Context, stream *entity.Stream) error
	Delete(ctx context.Context, id string) error
	ListLive(ctx context.Context) ([]entity.Stream, error)
	ListByBroadcaster(ctx context.Context, broadcasterID string) ([]entity.Stream, error)
}
