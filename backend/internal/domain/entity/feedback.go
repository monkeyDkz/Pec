package entity

import "time"

// Feedback captures a user's rating + free-form comment about the product.
// It feeds the product roadmap (Ce3.3.2) and stays attached to its author
// for GDPR cascade deletion.
type Feedback struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Rating    int       `json:"rating" gorm:"not null"` // 1..5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
