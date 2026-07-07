package dto

import (
	"time"

	"github.com/streampulse/backend/internal/domain/entity"
)

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=80"`
	Description string `json:"description" binding:"max=500"`
}

type UpdatePlaylistRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=80"`
	Description string `json:"description" binding:"max=500"`
}

type PlaylistResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	OwnerID     string          `json:"owner_id"`
	Tracks      []TrackResponse `json:"tracks"`
	CreatedAt   time.Time       `json:"created_at"`
}

type TrackResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration int    `json:"duration_seconds"`
	FileURL  string `json:"file_url"`
}

func PlaylistFrom(p *entity.Playlist) PlaylistResponse {
	out := PlaylistResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
		Tracks:      make([]TrackResponse, 0, len(p.Tracks)),
	}
	for _, t := range p.Tracks {
		out.Tracks = append(out.Tracks, TrackFrom(&t))
	}
	return out
}

func TrackFrom(t *entity.Track) TrackResponse {
	return TrackResponse{
		ID:       t.ID,
		Title:    t.Title,
		Artist:   t.Artist,
		Duration: t.Duration,
		FileURL:  t.FileURL,
	}
}
