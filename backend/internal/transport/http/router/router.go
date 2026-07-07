// Package router wires HTTP routes to their handlers.
package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/transport/http/handler"
	"github.com/streampulse/backend/internal/transport/http/middleware"
	"github.com/streampulse/backend/internal/transport/http/openapi"
)

// Handlers groups every HTTP handler. Wired in main.go.
type Handlers struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Stream   *handler.StreamHandler
	Playlist *handler.PlaylistHandler
	Track    *handler.TrackHandler
	Feedback *handler.FeedbackHandler
	Admin    *handler.AdminHandler
}

// New builds the Gin router with all middlewares and routes bound.
func New(cfg *config.Config, jwt *auth.JWTManager, h Handlers) http.Handler {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Tracing(cfg.OTELServiceName)) // creates the root span first
	r.Use(middleware.RequestLogger())              // logs carry trace_id
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.SecurityHeaders())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     splitCSV(cfg.CORSOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// OpenAPI spec + Swagger UI
	r.GET("/openapi.yaml", openapi.SpecHandler())
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	r.GET("/swagger/", openapi.UIHandler())

	// Serve uploaded audio files (audio for just_audio listeners).
	// The path is the same one returned by LocalStorage in the FileURL.
	if cfg.StorageLocalPath != "" {
		r.Static("/uploads", cfg.StorageLocalPath)
	}

	v1 := r.Group("/api/v1")
	{
		// --- Public auth endpoints (rate-limited)
		authGroup := v1.Group("/auth")
		authGroup.Use(middleware.RateLimiter(10, time.Minute))
		{
			authGroup.POST("/register", h.Auth.Register)
			authGroup.POST("/login", h.Auth.Login)
		}

		// --- Public read-only stream endpoints (anonymous listening allowed)
		v1.GET("/streams", h.Stream.ListLive)
		v1.GET("/streams/:id", h.Stream.Get)
		v1.GET("/streams/:id/listen", h.Stream.Listen)

		// --- WebSocket publish (browsers can't stream an HTTP body nor set an
		// Authorization header, so auth is via ?token= and RequireRole guards it).
		wsPublish := v1.Group("")
		wsPublish.Use(middleware.AuthMiddlewareQuery(jwt))
		wsPublish.Use(middleware.RequireRole("broadcaster", "admin"))
		{
			wsPublish.GET("/streams/:id/publish/ws", h.Stream.PublishWS)
		}

		// --- Authenticated endpoints
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwt))
		{
			// Refresh JWT (must be a valid current JWT)
			protected.POST("/auth/refresh", h.Auth.Refresh)

			users := protected.Group("/users")
			{
				users.GET("/me", h.User.Me)
				users.PUT("/me", h.User.UpdateMe)
				users.DELETE("/me", h.User.DeleteMe)       // GDPR: right to erasure
				users.GET("/me/data", h.User.ExportMyData) // GDPR: right of access
			}

			playlists := protected.Group("/playlists")
			{
				playlists.GET("", h.Playlist.ListMine)
				playlists.POST("", h.Playlist.Create)
				playlists.GET("/:id", h.Playlist.Get)
				playlists.PUT("/:id", h.Playlist.Update)
				playlists.DELETE("/:id", h.Playlist.Delete)
				playlists.POST("/:id/tracks/:trackId", h.Playlist.AddTrack)
				playlists.DELETE("/:id/tracks/:trackId", h.Playlist.RemoveTrack)
				playlists.PUT("/:id/tracks/reorder", h.Playlist.Reorder)
			}

			// Tracks (read for all authenticated users)
			protected.GET("/tracks", h.Track.List)
			protected.GET("/tracks/:id", h.Track.Get)

			// Feedback (any authenticated user)
			protected.POST("/feedback", h.Feedback.Submit)

			// --- Broadcaster + admin endpoints
			broadcaster := protected.Group("")
			broadcaster.Use(middleware.RequireRole("broadcaster", "admin"))
			{
				broadcaster.POST("/streams", h.Stream.Create)
				broadcaster.DELETE("/streams/:id", h.Stream.Delete)
				broadcaster.POST("/streams/:id/publish", h.Stream.Publish)
				broadcaster.POST("/tracks", h.Track.Upload)
				broadcaster.DELETE("/tracks/:id", h.Track.Delete)
			}

			// --- Admin only
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", h.Admin.ListUsers)
				admin.PUT("/users/:id/role", h.Admin.UpdateRole)
				admin.GET("/stats", h.Admin.Stats)
				admin.GET("/feedback", h.Feedback.AdminList)
			}
		}
	}

	return r
}

func splitCSV(s string) []string {
	if s == "" || s == "*" {
		return []string{"*"}
	}
	out := []string{}
	cur := ""
	for _, r := range s {
		if r == ',' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
