package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/infrastructure/observability"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/streampulse/backend/internal/infrastructure/streaming"
	"github.com/streampulse/backend/internal/transport/http/handler"
	"github.com/streampulse/backend/internal/transport/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	observability.InitLogger(cfg.LogLevel)

	shutdown, err := observability.InitTracer(cfg.OTELEndpoint, "streampulse-api")
	if err != nil {
		slog.Error("failed to init tracer", "err", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// --- Persistence ---
	db, err := persistence.NewDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect database", "err", err)
		os.Exit(1)
	}
	if err := persistence.AutoMigrate(db); err != nil {
		slog.Error("failed to migrate database", "err", err)
		os.Exit(1)
	}

	userRepo := persistence.NewUserRepository(db)
	streamRepo := persistence.NewStreamRepository(db)
	playlistRepo := persistence.NewPlaylistRepository(db)
	trackRepo := persistence.NewTrackRepository(db)
	statsRepo := persistence.NewStatsRepository(db)

	// --- Infrastructure services ---
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)
	hasher := auth.NewBcryptHasher()
	engine := streaming.NewEngine(streamRepo)

	// --- Use cases ---
	authUC := usecase.NewAuthUseCase(userRepo, jwtManager, hasher)
	userUC := usecase.NewUserUseCase(userRepo)
	streamUC := usecase.NewStreamUseCase(engine, streamRepo)
	playlistUC := usecase.NewPlaylistUseCase(playlistRepo, trackRepo)
	trackUC := usecase.NewTrackUseCase(trackRepo)
	adminUC := usecase.NewAdminUseCase(userRepo, statsRepo)

	seedAdmin(cfg, userRepo, hasher)

	// --- HTTP ---
	handlers := &router.Handlers{
		Auth:     handler.NewAuthHandler(authUC),
		User:     handler.NewUserHandler(userUC),
		Stream:   handler.NewStreamHandler(streamUC),
		Playlist: handler.NewPlaylistHandler(playlistUC),
		Track:    handler.NewTrackHandler(trackUC),
		Admin:    handler.NewAdminHandler(adminUC),
	}

	r := router.New(cfg, handlers)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // streaming endpoints are long-lived; no write deadline
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
}

// seedAdmin creates an admin account from ADMIN_EMAIL/ADMIN_PASSWORD if it does
// not yet exist, so the platform is usable out of the box.
func seedAdmin(cfg *config.Config, userRepo interface {
	FindByEmail(context.Context, string) (*entity.User, error)
	Create(context.Context, *entity.User) error
}, hasher auth.PasswordHasher) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		return
	}
	ctx := context.Background()
	if _, err := userRepo.FindByEmail(ctx, cfg.AdminEmail); err == nil {
		return // already exists
	}
	hashed, err := hasher.Hash(cfg.AdminPassword)
	if err != nil {
		slog.Error("seed admin: hash failed", "err", err)
		return
	}
	admin := &entity.User{
		Email:    cfg.AdminEmail,
		Username: "admin",
		Password: hashed,
		Role:     entity.RoleAdmin,
	}
	if err := userRepo.Create(ctx, admin); err != nil {
		slog.Error("seed admin: create failed", "err", err)
		return
	}
	slog.Info("seeded admin account", "email", cfg.AdminEmail)
}
