package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mizanmahi/aiusage/cli/internal/state"
	"github.com/mizanmahi/aiusage/types"
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

func TestRunPushSendsPendingSessionsAndUpdatesState(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")
	codexHome := filepath.Join(dir, "codex")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	lastPushedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	pushStartedAt := time.Date(2026, 6, 8, 12, 30, 0, 0, time.UTC)
	saveTestConfig(t, claudeHome, codexHome)
	saveTestState(t, lastPushedAt)

	copyFixture(t, filepath.Join("internal", "testdata", "claude", "claude-sess-001.jsonl"), filepath.Join(claudeHome, "projects", "-home-dev-go-myapp", "claude-sess-001.jsonl"))
	copyFixture(t, filepath.Join("internal", "testdata", "codex", "last-token-usage.jsonl"), filepath.Join(codexHome, "sessions", "2026", "06", "03", "last-token-usage.jsonl"))

	oldSender := sendUsageEvents
	oldClock := currentTime
	defer func() {
		sendUsageEvents = oldSender
		currentTime = oldClock
	}()

	currentTime = func() time.Time {
		return pushStartedAt
	}

	var gotServerURL string
	var gotAPIKey string
	var gotEvents []types.UsageEvent
	sendUsageEvents = func(serverURL, apiKey, clientVersion string, events []types.UsageEvent) (*types.PushResponse, error) {
		gotServerURL = serverURL
		gotAPIKey = apiKey
		gotEvents = events
		return &types.PushResponse{Accepted: 2, Skipped: 0, Message: "ok"}, nil
	}

	var out, errOut bytes.Buffer
	if err := runPush(&out, &errOut, false); err != nil {
		t.Fatalf("runPush() error = %v", err)
	}

	if gotServerURL != "http://localhost:8080" {
		t.Fatalf("serverURL = %q, want configured server URL", gotServerURL)
	}
	if gotAPIKey != "ak_secret_value" {
		t.Fatalf("apiKey = %q, want configured API key", gotAPIKey)
	}
	if len(gotEvents) != 2 {
		t.Fatalf("events len = %d, want 2", len(gotEvents))
	}
	output := out.String()
	assertContains(t, output, "Pushed sessions: 2 accepted, 0 skipped")
	assertContains(t, output, "ok")
	if strings.Contains(output, "ak_secret_value") {
		t.Fatal("runPush() printed the API key")
	}
	if errOut.Len() != 0 {
		t.Fatalf("runPush() stderr = %q, want empty", errOut.String())
	}

	gotState, err := state.Load()
	if err != nil {
		t.Fatalf("state.Load() error = %v", err)
	}
	if !gotState.LastPushedAt.Equal(pushStartedAt) {
		t.Fatalf("LastPushedAt = %s, want %s", gotState.LastPushedAt, pushStartedAt)
	}
}

func TestRunPushDoesNotUpdateStateWhenSendFails(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")
	claudeHome := filepath.Join(dir, "claude")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	lastPushedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	saveTestConfig(t, claudeHome, filepath.Join(dir, "missing-codex"))
	saveTestState(t, lastPushedAt)
	copyFixture(t, filepath.Join("internal", "testdata", "claude", "claude-sess-001.jsonl"), filepath.Join(claudeHome, "projects", "-home-dev-go-myapp", "claude-sess-001.jsonl"))

	oldSender := sendUsageEvents
	defer func() {
		sendUsageEvents = oldSender
	}()
	sendUsageEvents = func(serverURL, apiKey, clientVersion string, events []types.UsageEvent) (*types.PushResponse, error) {
		return nil, errors.New("server unavailable")
	}

	var out, errOut bytes.Buffer
	err := runPush(&out, &errOut, false)
	if err == nil {
		t.Fatal("runPush() error = nil, want send error")
	}
	if !strings.Contains(err.Error(), "push failed: server unavailable") {
		t.Fatalf("runPush() error = %q, want wrapped send error", err.Error())
	}

	gotState, err := state.Load()
	if err != nil {
		t.Fatalf("state.Load() error = %v", err)
	}
	if !gotState.LastPushedAt.Equal(lastPushedAt) {
		t.Fatalf("LastPushedAt = %s, want unchanged %s", gotState.LastPushedAt, lastPushedAt)
	}
}

func TestRunPushSkipsSendWhenNoPendingSessions(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.json")

	t.Setenv("AIUSAGE_CONFIG_PATH", configPath)
	t.Setenv("AIUSAGE_STATE_PATH", statePath)

	saveTestConfig(t, filepath.Join(dir, "missing-claude"), filepath.Join(dir, "missing-codex"))

	oldSender := sendUsageEvents
	defer func() {
		sendUsageEvents = oldSender
	}()
	sendUsageEvents = func(serverURL, apiKey, clientVersion string, events []types.UsageEvent) (*types.PushResponse, error) {
		t.Fatal("sendUsageEvents() called with no pending sessions")
		return nil, nil
	}

	var out, errOut bytes.Buffer
	if err := runPush(&out, &errOut, false); err != nil {
		t.Fatalf("runPush() error = %v", err)
	}

	assertContains(t, out.String(), "No pending sessions.")
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Fatalf("state file exists after no-op push, Stat error = %v", err)
	}
}

func assertContains(t *testing.T, text, want string) {
	t.Helper()

	if !strings.Contains(text, want) {
		t.Fatalf("output missing %q:\n%s", want, text)
	}
}
