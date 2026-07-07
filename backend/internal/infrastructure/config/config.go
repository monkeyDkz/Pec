// Package config loads the application's configuration from environment
// variables, following the 12-Factor App methodology.
// See docs/cahier-des-charges-fr.md § 5.2.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port              string        `mapstructure:"PORT"`
	DatabaseURL       string        `mapstructure:"DATABASE_URL"`
	JWTSecret         string        `mapstructure:"JWT_SECRET"`
	JWTExpiration     time.Duration `mapstructure:"JWT_EXPIRATION"`
	OTELEndpoint      string        `mapstructure:"OTEL_ENDPOINT"`
	OTELServiceName   string        `mapstructure:"OTEL_SERVICE_NAME"`
	LogLevel          string        `mapstructure:"LOG_LEVEL"`
	Environment       string        `mapstructure:"ENVIRONMENT"`
	CORSOrigins       string        `mapstructure:"CORS_ORIGINS"`
	StorageBackend    string        `mapstructure:"STORAGE_BACKEND"`
	StorageLocalPath  string        `mapstructure:"STORAGE_LOCAL_PATH"`
	MetricsRequireAuth bool         `mapstructure:"METRICS_REQUIRE_AUTH"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("OTEL_ENDPOINT", "localhost:4317")
	viper.SetDefault("OTEL_SERVICE_NAME", "streampulse-api")
	viper.SetDefault("CORS_ORIGINS", "*")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("STORAGE_BACKEND", "local")
	viper.SetDefault("STORAGE_LOCAL_PATH", "./uploads")
	viper.SetDefault("METRICS_REQUIRE_AUTH", false)

	// Keys without a default are not known to viper, so AutomaticEnv alone
	// won't surface them through Unmarshal. Bind them explicitly so they are
	// read from the environment even when no .env file is present (e.g. in
	// containers, where config comes purely from env vars).
	for _, key := range []string{
		"PORT", "DATABASE_URL", "JWT_SECRET", "JWT_EXPIRATION",
		"OTEL_ENDPOINT", "OTEL_SERVICE_NAME", "LOG_LEVEL", "ENVIRONMENT",
		"CORS_ORIGINS", "STORAGE_BACKEND", "STORAGE_LOCAL_PATH",
		"METRICS_REQUIRE_AUTH",
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
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if cfg.JWTExpiration <= 0 {
		cfg.JWTExpiration = 24 * time.Hour
	}

	return &cfg, nil
}
