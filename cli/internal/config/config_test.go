package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".aiusage", "config.toml")
	t.Setenv("AIUSAGE_CONFIG_PATH", path)

	want := &Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_test_local",
		ClaudePath: "/tmp/claude",
		CodexPath:  "/tmp/codex",
	}
	if err := Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if gotMode := info.Mode().Perm(); gotMode != 0600 {
		t.Fatalf("config file mode = %o, want 0600", gotMode)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.ServerURL != want.ServerURL {
		t.Errorf("ServerURL = %q, want %q", got.ServerURL, want.ServerURL)
	}
	if got.APIKey != want.APIKey {
		t.Errorf("APIKey = %q, want %q", got.APIKey, want.APIKey)
	}
	if got.ClaudePath != want.ClaudePath {
		t.Errorf("ClaudePath = %q, want %q", got.ClaudePath, want.ClaudePath)
	}
	if got.CodexPath != want.CodexPath {
		t.Errorf("CodexPath = %q, want %q", got.CodexPath, want.CodexPath)
	}
}

func TestLoadAppliesDefaultToolPaths(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	t.Setenv("AIUSAGE_CONFIG_PATH", path)
	t.Setenv("HOME", dir)

	data := []byte("server_url = \"http://localhost:8080\"\napi_key = \"ak_test_local\"\n")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if want := filepath.Join(dir, ".claude"); got.ClaudePath != want {
		t.Errorf("ClaudePath = %q, want %q", got.ClaudePath, want)
	}
	if want := filepath.Join(dir, ".codex"); got.CodexPath != want {
		t.Errorf("CodexPath = %q, want %q", got.CodexPath, want)
	}
}

func TestLoadMissingConfigReturnsError(t *testing.T) {
	t.Setenv("AIUSAGE_CONFIG_PATH", filepath.Join(t.TempDir(), "missing.toml"))

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}

func TestConfigPathUsesOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom-config.toml")
	t.Setenv("AIUSAGE_CONFIG_PATH", path)

	got, err := configPath()
	if err != nil {
		t.Fatalf("configPath() error = %v", err)
	}
	if got != path {
		t.Fatalf("configPath() = %q, want %q", got, path)
	}
}
