// Package dto holds request/response payloads exchanged over HTTP, decoupled
// from domain entities.
package dto

import (
	"time"

	"github.com/streampulse/backend/internal/domain/entity"
)

// ---- Auth ----

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ---- User ----

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
	}
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user broadcaster admin"`
}

// ---- Stream ----

type CreateStreamRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=120"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type StreamResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	BroadcasterID string    `json:"broadcaster_id"`
	Broadcaster   string    `json:"broadcaster"`
	Status        string    `json:"status"`
	ListenerCount int       `json:"listener_count"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewStreamResponse(s *entity.Stream) StreamResponse {
	return StreamResponse{
		ID:            s.ID,
		Title:         s.Title,
		Description:   s.Description,
		BroadcasterID: s.BroadcasterID,
		Broadcaster:   s.Broadcaster.Username,
		Status:        string(s.Status),
		ListenerCount: s.ListenerCount,
		CreatedAt:     s.CreatedAt,
	}
}

// ---- Track ----

type CreateTrackRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=200"`
	Artist   string `json:"artist" binding:"omitempty,max=200"`
	Duration int    `json:"duration" binding:"omitempty,min=0"`
	FileURL  string `json:"file_url" binding:"required,url"`
}

type TrackResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Artist    string    `json:"artist"`
	Duration  int       `json:"duration"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTrackResponse(t *entity.Track) TrackResponse {
	return TrackResponse{
		ID:        t.ID,
		Title:     t.Title,
		Artist:    t.Artist,
		Duration:  t.Duration,
		FileURL:   t.FileURL,
		CreatedAt: t.CreatedAt,
	}
}

// ---- Playlist ----

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=120"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdatePlaylistRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=120"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type AddTrackRequest struct {
	TrackID string `json:"track_id" binding:"required,uuid"`
}

type PlaylistResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	OwnerID     string          `json:"owner_id"`
	Tracks      []TrackResponse `json:"tracks"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewPlaylistResponse(p *entity.Playlist) PlaylistResponse {
	tracks := make([]TrackResponse, 0, len(p.Tracks))
	for i := range p.Tracks {
		tracks = append(tracks, NewTrackResponse(&p.Tracks[i]))
	}
	return PlaylistResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		Tracks:      tracks,
		CreatedAt:   p.CreatedAt,
	}
}

// ---- Admin ----

type AdminStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	TotalStreams  int64 `json:"total_streams"`
	LiveStreams   int64 `json:"live_streams"`
	TotalTracks   int64 `json:"total_tracks"`
	TotalPlaylist int64 `json:"total_playlists"`
}
