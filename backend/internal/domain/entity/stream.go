package entity

import "time"

type StreamStatus string

const (
	StreamStatusLive    StreamStatus = "live"
	StreamStatusOffline StreamStatus = "offline"
)

type Stream struct {
	ID            string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title         string       `json:"title" gorm:"not null"`
	Description   string       `json:"description"`
	BroadcasterID string       `json:"broadcaster_id" gorm:"type:uuid;not null"`
	Broadcaster   User         `json:"broadcaster" gorm:"foreignKey:BroadcasterID"`
	Status        StreamStatus `json:"status" gorm:"type:varchar(20);default:'offline'"`
	ListenerCount int          `json:"listener_count" gorm:"-"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
