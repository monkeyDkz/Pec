// Package main is the entry point of the StreamPulse API.
// It wires configuration, observability, persistence, the streaming registry
// and HTTP handlers, then runs the server with graceful shutdown.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/infrastructure/observability"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/streampulse/backend/internal/infrastructure/streaming"
	"github.com/streampulse/backend/internal/transport/http/handler"
	"github.com/streampulse/backend/internal/transport/http/router"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	observability.InitLogger(cfg.LogLevel)
	slog.Info("config loaded", "env", cfg.Environment, "port", cfg.Port)

	// Observability
	shutdownTracer, err := observability.InitTracer(cfg.OTELEndpoint, cfg.OTELServiceName)
	if err != nil {
		slog.Warn("tracer init failed; continuing without OTEL", "err", err)
	}
	defer func() {
		if shutdownTracer != nil {
			_ = shutdownTracer(context.Background())
		}
	}()

	// Persistence
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()
	db, err := persistence.Open(dbCtx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	if err := persistence.AutoMigrate(db); err != nil {
		return err
	}
	slog.Info("database connected and migrated")

	// Repositories
	userRepo := persistence.NewUserRepository(db)
	streamRepo := persistence.NewStreamRepository(db)
	playlistRepo := persistence.NewPlaylistRepository(db)
	trackRepo := persistence.NewTrackRepository(db)
	feedbackRepo := persistence.NewFeedbackRepository(db)

	// Auth + streaming registry
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	hasher := auth.NewBcryptHasher()
	registry := streaming.NewRegistry()

	// File storage
	publicURL := "http://localhost:" + cfg.Port + "/uploads"
	storage, err := usecase.NewLocalStorage(cfg.StorageLocalPath, publicURL)
	if err != nil {
		return err
	}

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, jwtManager, hasher)
	userUC := usecase.NewUserUseCase(userRepo, playlistRepo, streamRepo)
	streamUC := usecase.NewStreamUseCase(streamRepo, registry)
	playlistUC := usecase.NewPlaylistUseCase(playlistRepo, trackRepo)
	trackUC := usecase.NewTrackUseCase(trackRepo, storage)
	feedbackUC := usecase.NewFeedbackUseCase(feedbackRepo)
	statsUC := usecase.NewStatsUseCase(db)

	// HTTP handlers
	handlers := router.Handlers{
		Auth:     handler.NewAuthHandler(authUC),
		User:     handler.NewUserHandler(userUC),
		Stream:   handler.NewStreamHandler(streamUC),
		Playlist: handler.NewPlaylistHandler(playlistUC),
		Track:    handler.NewTrackHandler(trackUC),
		Feedback: handler.NewFeedbackHandler(feedbackUC),
		Admin:    handler.NewAdminHandler(userUC, statsUC),
	}

	r := router.New(cfg, jwtManager, handlers)

	srv := &http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     r,
		ReadTimeout: 0, // streaming endpoints need long-lived reads
		// WriteTimeout intentionally 0 for the same reason; per-handler
		// timeouts are enforced via context.
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		slog.Info("server listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
	}
	// Close any remaining live streams cleanly.
	for _, id := range registry.Snapshot() {
		registry.CloseStream(id)
	}
	return nil
}
