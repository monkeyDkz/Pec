package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/transport/http/handler"
	"github.com/streampulse/backend/internal/transport/http/middleware"
)

// Handlers bundles all HTTP handlers wired into the router.
type Handlers struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Stream   *handler.StreamHandler
	Playlist *handler.PlaylistHandler
	Track    *handler.TrackHandler
	Admin    *handler.AdminHandler
}

func New(cfg *config.Config, h *Handlers) http.Handler {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.MetricsMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health & metrics
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	const (
		roleBroadcaster = "broadcaster"
		roleAdmin       = "admin"
	)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/me", h.User.Me)
				users.PUT("/me", h.User.UpdateMe)
			}

			// Streams
			streams := protected.Group("/streams")
			{
				streams.GET("", h.Stream.List)
				streams.GET("/:id", h.Stream.Get)
				streams.GET("/:id/listen", h.Stream.Listen)

				broadcast := streams.Group("")
				broadcast.Use(middleware.RequireRole(roleBroadcaster, roleAdmin))
				{
					broadcast.POST("", h.Stream.Create)
					broadcast.POST("/:id/publish", h.Stream.Publish)
					broadcast.POST("/:id/stop", h.Stream.Stop)
					broadcast.DELETE("/:id", h.Stream.Delete)
				}
			}

			// Playlists (any authenticated user)
			playlists := protected.Group("/playlists")
			{
				playlists.GET("", h.Playlist.List)
				playlists.POST("", h.Playlist.Create)
				playlists.GET("/:id", h.Playlist.Get)
				playlists.PUT("/:id", h.Playlist.Update)
				playlists.DELETE("/:id", h.Playlist.Delete)
				playlists.POST("/:id/tracks", h.Playlist.AddTrack)
				playlists.DELETE("/:id/tracks/:trackId", h.Playlist.RemoveTrack)
			}

			// Tracks
			tracks := protected.Group("/tracks")
			{
				tracks.GET("", h.Track.List)
				tracks.GET("/:id", h.Track.Get)

				tracksWrite := tracks.Group("")
				tracksWrite.Use(middleware.RequireRole(roleBroadcaster, roleAdmin))
				{
					tracksWrite.POST("", h.Track.Create)
					tracksWrite.DELETE("/:id", h.Track.Delete)
				}
			}

			// Admin
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole(roleAdmin))
			{
				admin.GET("/users", h.Admin.ListUsers)
				admin.PUT("/users/:id/role", h.Admin.UpdateRole)
				admin.GET("/stats", h.Admin.Stats)
			}
		}
	}

	return r
}
