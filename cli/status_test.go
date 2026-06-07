package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mizanmahi/aiusage/cli/internal/config"
	"github.com/mizanmahi/aiusage/cli/internal/state"
)

func TestRunStatusPrintsPendingSessionsWithoutAPIKey(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")
	codexHome := filepath.Join(dir, "codex")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	if err := config.Save(&config.Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_secret_value",
		ClaudePath: claudeHome,
		CodexPath:  codexHome,
	}); err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}
	if err := state.Save(&state.State{}); err != nil {
		t.Fatalf("state.Save() error = %v", err)
	}

	copyFixture(t, filepath.Join("internal", "testdata", "claude", "claude-sess-001.jsonl"), filepath.Join(claudeHome, "projects", "-home-dev-go-myapp", "claude-sess-001.jsonl"))
	copyFixture(t, filepath.Join("internal", "testdata", "codex", "last-token-usage.jsonl"), filepath.Join(codexHome, "sessions", "2026", "06", "03", "last-token-usage.jsonl"))

	var out, errOut bytes.Buffer
	if err := runStatus(&out, &errOut); err != nil {
		t.Fatalf("runStatus() error = %v", err)
	}

	output := out.String()
	if strings.Contains(output, "ak_secret_value") {
		t.Fatal("runStatus() printed the API key")
	}
	if !strings.Contains(output, "Server URL: http://localhost:8080") {
		t.Fatalf("runStatus() output missing server URL: %q", output)
	}
	if !strings.Contains(output, "Last push: never") {
		t.Fatalf("runStatus() output missing last push: %q", output)
	}
	if !strings.Contains(output, "Pending sessions: 2 (1 Claude, 1 Codex)") {
		t.Fatalf("runStatus() output missing pending count: %q", output)
	}
	if errOut.Len() != 0 {
		t.Fatalf("runStatus() stderr = %q, want empty", errOut.String())
	}
}

func TestRunStatusRespectsLastPushCursor(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")
	codexHome := filepath.Join(dir, "codex")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	if err := config.Save(&config.Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_secret_value",
		ClaudePath: claudeHome,
		CodexPath:  codexHome,
	}); err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}
	if err := state.Save(&state.State{LastPushedAt: time.Now().Add(24 * time.Hour)}); err != nil {
		t.Fatalf("state.Save() error = %v", err)
	}

	copyFixture(t, filepath.Join("internal", "testdata", "codex", "last-token-usage.jsonl"), filepath.Join(codexHome, "sessions", "2026", "06", "03", "last-token-usage.jsonl"))

	var out, errOut bytes.Buffer
	if err := runStatus(&out, &errOut); err != nil {
		t.Fatalf("runStatus() error = %v", err)
	}

	if !strings.Contains(out.String(), "Pending sessions: 0 (0 Claude, 0 Codex)") {
		t.Fatalf("runStatus() output missing pending count: %q", out.String())
	}
}

func TestRunStatusRequiresConfig(t *testing.T) {
	t.Setenv("AIUSAGE_CONFIG_PATH", filepath.Join(t.TempDir(), "missing.toml"))
	t.Setenv("AIUSAGE_STATE_PATH", filepath.Join(t.TempDir(), "state.json"))

	var out, errOut bytes.Buffer
	if err := runStatus(&out, &errOut); err == nil {
		t.Fatal("runStatus() error = nil, want error")
	}
}

func copyFixture(t *testing.T, src, dst string) {
	t.Helper()

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
