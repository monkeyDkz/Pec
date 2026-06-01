package entity

import "time"

type Role string

const (
	RoleAnonymous   Role = "anonymous"
	RoleUser        Role = "user"
	RoleBroadcaster Role = "broadcaster"
	RoleAdmin       Role = "admin"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      Role      `json:"role" gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
