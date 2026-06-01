package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/transport/http/middleware"
)

func New(cfg *config.Config) http.Handler {
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

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", placeholder)
			auth.POST("/login", placeholder)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/me", placeholder)
				users.PUT("/me", placeholder)
			}

			// Streams
			streams := protected.Group("/streams")
			{
				streams.GET("", placeholder)
				streams.POST("", placeholder)
				streams.GET("/:id", placeholder)
				streams.DELETE("/:id", placeholder)
				streams.GET("/:id/listen", placeholder)
				streams.POST("/:id/publish", placeholder)
			}

			// Playlists
			playlists := protected.Group("/playlists")
			{
				playlists.GET("", placeholder)
				playlists.POST("", placeholder)
				playlists.GET("/:id", placeholder)
				playlists.PUT("/:id", placeholder)
				playlists.DELETE("/:id", placeholder)
				playlists.POST("/:id/tracks", placeholder)
				playlists.DELETE("/:id/tracks/:trackId", placeholder)
			}

			// Tracks
			tracks := protected.Group("/tracks")
			{
				tracks.GET("", placeholder)
				tracks.POST("", placeholder)
				tracks.GET("/:id", placeholder)
				tracks.DELETE("/:id", placeholder)
			}

			// Admin
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", placeholder)
				admin.PUT("/users/:id/role", placeholder)
				admin.GET("/stats", placeholder)
			}
		}
	}

	return r
}

func placeholder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
