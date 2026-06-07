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

func TestRunPushDryRunPrintsPreviewWithoutAPIKey(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")
	codexHome := filepath.Join(dir, "codex")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	lastPushedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	if err := config.Save(&config.Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_secret_value",
		ClaudePath: claudeHome,
		CodexPath:  codexHome,
	}); err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}
	if err := state.Save(&state.State{LastPushedAt: lastPushedAt}); err != nil {
		t.Fatalf("state.Save() error = %v", err)
	}

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
	assertContains(t, output, "claude claude-sess-001 project=myapp date=2026-06-01 model=claude-sonnet-4-5 input=7000 output=450 cache=6300 reasoning=0")
	assertContains(t, output, "codex last-token-session project=myproject date=2026-06-03 model=gpt-5.5 input=1500 output=300 cache=1200 reasoning=70")

	if strings.Contains(output, "ak_secret_value") {
		t.Fatal("runPush() printed the API key")
	}
	if errOut.Len() != 0 {
		t.Fatalf("runPush() stderr = %q, want empty", errOut.String())
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

	if err := config.Save(&config.Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_secret_value",
		ClaudePath: filepath.Join(dir, "missing-claude"),
		CodexPath:  filepath.Join(dir, "missing-codex"),
	}); err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}

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

func TestRunPushRequiresDryRunForNow(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := runPush(&out, &errOut, false); err == nil {
		t.Fatal("runPush() error = nil, want error")
	}
}

func assertContains(t *testing.T, text, want string) {
	t.Helper()

	if !strings.Contains(text, want) {
		t.Fatalf("output missing %q:\n%s", want, text)
	}
}
