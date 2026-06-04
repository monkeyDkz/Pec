package persistence

import (
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a GORM connection to PostgreSQL using the provided DSN/URL.
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	return db, nil
}

// AutoMigrate ensures the required extensions and schema are present.
func AutoMigrate(db *gorm.DB) error {
	// gen_random_uuid() lives in pgcrypto on older Postgres; harmless on PG16.
	db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`)

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Track{},
		&entity.Stream{},
		&entity.Playlist{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
