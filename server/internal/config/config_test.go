package config

import "testing"

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("MIN_CLI_VERSION", "0.2.0")
	t.Setenv("STATIC_DIR", "/var/www/aiusage")
	t.Setenv("CORS_ORIGINS", "http://localhost:5173")

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
	if cfg.StaticDir != "/var/www/aiusage" {
		t.Fatalf("StaticDir = %q, want /var/www/aiusage", cfg.StaticDir)
	}
	if cfg.CORSOrigins != "http://localhost:5173" {
		t.Fatalf("CORSOrigins = %q, want http://localhost:5173", cfg.CORSOrigins)
	}
}

func TestLoadDefaultsOptionalValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("PORT", "")
	t.Setenv("ENV", "")
	t.Setenv("MIN_CLI_VERSION", "")
	t.Setenv("STATIC_DIR", "")
	t.Setenv("CORS_ORIGINS", "")

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
	if cfg.StaticDir != "../ui/dist" {
		t.Fatalf("StaticDir = %q, want ../ui/dist", cfg.StaticDir)
	}
	if cfg.CORSOrigins != "*" {
		t.Fatalf("CORSOrigins = %q, want *", cfg.CORSOrigins)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want required database error")
	}
}
