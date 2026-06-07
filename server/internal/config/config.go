package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	Port          string
	Env           string
	MinCLIVersion string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          getenv("PORT", "8080"),
		Env:           getenv("ENV", "development"),
		MinCLIVersion: os.Getenv("MIN_CLI_VERSION"),
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
