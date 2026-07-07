// Package persistence wires GORM to PostgreSQL and exposes repository
// implementations. See docs/adr/0005-postgresql-database.md.
package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open returns a GORM DB connection configured for production-grade use:
// connection pooling, sane timeouts and reduced logging.
func Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Warn),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	// Instrument GORM with OpenTelemetry: every query becomes a child span
	// of the incoming HTTP request, propagated via context.
	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("streampulse"))); err != nil {
		return nil, fmt.Errorf("install otelgorm: %w", err)
	}

	return db, nil
}

// AutoMigrate creates or updates the schema for all known entities.
// In production, prefer versioned SQL migrations (see backend/migrations).
func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return fmt.Errorf("enable pgcrypto: %w", err)
	}
	return db.AutoMigrate(
		&entity.User{},
		&entity.Stream{},
		&entity.Playlist{},
		&entity.Track{},
		&entity.Feedback{},
	)
}
