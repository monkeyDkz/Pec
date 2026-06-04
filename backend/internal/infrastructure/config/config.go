package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string        `mapstructure:"PORT"`
	DatabaseURL   string        `mapstructure:"DATABASE_URL"`
	JWTSecret     string        `mapstructure:"JWT_SECRET"`
	JWTDuration   time.Duration `mapstructure:"JWT_DURATION"`
	OTELEndpoint  string        `mapstructure:"OTEL_ENDPOINT"`
	LogLevel      string        `mapstructure:"LOG_LEVEL"`
	Environment   string        `mapstructure:"ENVIRONMENT"`
	CORSOrigins   string        `mapstructure:"CORS_ORIGINS"`
	AdminEmail    string        `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string        `mapstructure:"ADMIN_PASSWORD"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("OTEL_ENDPOINT", "localhost:4317")
	viper.SetDefault("CORS_ORIGINS", "*")
	viper.SetDefault("JWT_DURATION", "24h")

	// AutomaticEnv alone does NOT feed Unmarshal for keys without a default,
	// so bind every env var explicitly (works with or without a .env file).
	for _, key := range []string{
		"PORT", "DATABASE_URL", "JWT_SECRET", "JWT_DURATION",
		"OTEL_ENDPOINT", "LOG_LEVEL", "ENVIRONMENT", "CORS_ORIGINS",
		"ADMIN_EMAIL", "ADMIN_PASSWORD",
	} {
		_ = viper.BindEnv(key)
	}

	// .env file is optional — env vars take precedence
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.JWTDuration == 0 {
		cfg.JWTDuration = 24 * time.Hour
	}

	return &cfg, nil
}
