package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunPushDryRunPrintsPreviewWithoutAPIKey(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")
	codexHome := filepath.Join(dir, "codex")
	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	lastPushedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	saveTestConfig(t, claudeHome, codexHome)
	saveTestState(t, lastPushedAt)
	copyFixture(t, filepath.Join("internal", "testdata", "claude", "claude-sess-001.jsonl"), filepath.Join(claudeHome, "projects", "-home-dev-go-myapp", "claude-sess-001.jsonl"))
	copyFixture(t, filepath.Join("internal", "testdata", "codex", "last-token-usage.jsonl"), filepath.Join(codexHome, "sessions", "2026", "06", "03", "last-token-usage.jsonl"))

	var out, errOut bytes.Buffer
	if err := runPush(&out, &errOut, true); err != nil {
		t.Fatalf("runPush() error = %v", err)
	}

	output := out.String()
	assertContains(t, output, "Dry run: no data sent.")
	assertContains(t, output, "Server URL: http://localhost:8080")
	assertContains(t, output, "Last push: 2026-06-01T00:00:00Z")
	assertContains(t, output, "Pending sessions: 2 (1 Claude, 1 Codex)")
	assertContains(t, output, "claude claude-sess-001 project=myapp date=2026-06-01 model=claude-sonnet-4-5 input=7000 output=450 cache_create=3000 cache_read=3300 reasoning=0")
	assertContains(t, output, "codex last-token-session project=myproject date=2026-06-03 model=gpt-5.5 input=1500 output=300 cache_create=0 cache_read=1200 reasoning=70")
	if strings.Contains(output, "ak_secret_value") || errOut.Len() != 0 {
		t.Fatal("dry run exposed the API key or wrote warnings")
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", statePath, err)
	}
	if !strings.Contains(string(data), `"last_pushed_at": "2026-06-01T00:00:00Z"`) {
		t.Fatalf("state file changed after dry run: %s", string(data))
	}
}

func TestRunPushDryRunDoesNotCreateStateFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)
	saveTestConfig(t, filepath.Join(dir, "missing-claude"), filepath.Join(dir, "missing-codex"))

	var out, errOut bytes.Buffer
	if err := runPush(&out, &errOut, true); err != nil {
		t.Fatalf("runPush() error = %v", err)
	}
	assertContains(t, out.String(), "Last push: never")
	assertContains(t, out.String(), "Pending sessions: 0 (0 Claude, 0 Codex)")
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Fatalf("state file exists after dry run, Stat error = %v", err)
	}
}
