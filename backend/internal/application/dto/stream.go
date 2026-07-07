package dto

import (
	"time"

	"github.com/streampulse/backend/internal/domain/entity"
)

type CreateStreamRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=120"`
	Description string `json:"description" binding:"max=1000"`
}

type StreamResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	BroadcasterID string    `json:"broadcaster_id"`
	Broadcaster   string    `json:"broadcaster_username,omitempty"`
	Status        string    `json:"status"`
	ListenerCount int       `json:"listener_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func StreamFrom(s *entity.Stream) StreamResponse {
	return StreamResponse{
		ID:            s.ID,
		Title:         s.Title,
		Description:   s.Description,
		BroadcasterID: s.BroadcasterID,
		Broadcaster:   s.Broadcaster.Username,
		Status:        string(s.Status),
		ListenerCount: s.ListenerCount,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
