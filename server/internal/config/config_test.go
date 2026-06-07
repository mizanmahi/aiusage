package config

import "testing"

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("MIN_CLI_VERSION", "0.2.0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("DatabaseURL = %q, want postgres://example", cfg.DatabaseURL)
	}
	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Fatalf("Env = %q, want production", cfg.Env)
	}
	if cfg.MinCLIVersion != "0.2.0" {
		t.Fatalf("MinCLIVersion = %q, want 0.2.0", cfg.MinCLIVersion)
	}
}

func TestLoadDefaultsOptionalValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Fatalf("Env = %q, want development", cfg.Env)
	}
	if cfg.MinCLIVersion != "" {
		t.Fatalf("MinCLIVersion = %q, want empty", cfg.MinCLIVersion)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want required database error")
	}
}
