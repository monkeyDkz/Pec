package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	OTELEndpoint string `mapstructure:"OTEL_ENDPOINT"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
	Environment  string `mapstructure:"ENVIRONMENT"`
	CORSOrigins  string `mapstructure:"CORS_ORIGINS"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("OTEL_ENDPOINT", "localhost:4317")
	viper.SetDefault("CORS_ORIGINS", "*")

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

	return &cfg, nil
}
