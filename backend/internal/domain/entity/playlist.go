package entity

import "time"

type Playlist struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id" gorm:"type:uuid;not null"`
	Owner       User      `json:"owner" gorm:"foreignKey:OwnerID"`
	Tracks      []Track   `json:"tracks" gorm:"many2many:playlist_tracks;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Track struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title     string    `json:"title" gorm:"not null"`
	Artist    string    `json:"artist"`
	Duration  int       `json:"duration"` // seconds
	FileURL   string    `json:"file_url" gorm:"not null"`
	UploadBy  string    `json:"upload_by" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
}
