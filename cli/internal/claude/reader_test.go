package claude

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mizanmahi/aiusage/types"
)

func TestParseSessionFile(t *testing.T) {
	path := copyFixtureToProjectPath(t, "claude-sess-001.jsonl", "-home-dev-go-myapp")

	session, err := parseSessionFile(path, "myapp", "/home/dev/go/myapp")
	if err != nil {
		t.Fatalf("parseSessionFile() error = %v", err)
	}
	if session == nil {
		t.Fatal("parseSessionFile() returned nil")
	}

	assertSession(t, session, "claude-sess-001", "myapp", "/home/dev/go/myapp", 7000, 450, 6300)
	if session.Model != "claude-sonnet-4-5" {
		t.Errorf("Model = %q, want claude-sonnet-4-5", session.Model)
	}
	if session.ReasoningTokens != 0 {
		t.Errorf("ReasoningTokens = %d, want 0", session.ReasoningTokens)
	}

	event := session.ToUsageEvent("user-123")
	if event.Tool != types.ToolClaude {
		t.Errorf("ToUsageEvent().Tool = %q, want %q", event.Tool, types.ToolClaude)
	}
	if event.UserID != "user-123" {
		t.Errorf("ToUsageEvent().UserID = %q, want user-123", event.UserID)
	}
	if event.CacheTokens != 6300 {
		t.Errorf("ToUsageEvent().CacheTokens = %d, want 6300", event.CacheTokens)
	}
}

func TestDecodePath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "normal encoded cwd", in: "-home-dev-go-myapp", want: "myapp"},
		{name: "nested path basename", in: "-home-dev-work-backend-api", want: "api"},
		{name: "root project", in: "-root-project", want: "project"},
		{name: "empty-ish value", in: "-", want: "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodePath(tt.in)
			if got != tt.want {
				t.Errorf("decodePath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEncodedToPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "normal encoded cwd", in: "-home-dev-go-myapp", want: "/home/dev/go/myapp"},
		{name: "root", in: "-", want: "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodedToPath(tt.in)
			if got != tt.want {
				t.Errorf("encodedToPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestReadSessionsSinceFilter(t *testing.T) {
	dir := t.TempDir()
	path := copyFixtureToProjectPathInDir(t, dir, "claude-sess-001.jsonl", "-home-dev-go-myapp")

	past := time.Now().Add(-24 * time.Hour)
	sessions, err := ReadSessions(dir, past)
	if err != nil {
		t.Fatalf("ReadSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("ReadSessions() got %d sessions, want 1", len(sessions))
	}
	assertSession(t, &sessions[0], "claude-sess-001", "myapp", "/home/dev/go/myapp", 7000, 450, 6300)

	future := time.Now().Add(24 * time.Hour)
	sessions, err = ReadSessions(dir, future)
	if err != nil {
		t.Fatalf("ReadSessions() future filter error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("ReadSessions() future filter got %d sessions, want 0", len(sessions))
	}

	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatalf("Chtimes() error = %v", err)
	}
	sessions, err = ReadSessions(dir, time.Now())
	if err != nil {
		t.Fatalf("ReadSessions() old file error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("ReadSessions() old file got %d sessions, want 0", len(sessions))
	}
}

func TestReadSessionsMissingDirectory(t *testing.T) {
	sessions, err := ReadSessions(t.TempDir(), time.Time{})
	if err != nil {
		t.Fatalf("ReadSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("ReadSessions() got %d sessions, want 0", len(sessions))
	}
}

func assertSession(t *testing.T, got *Session, wantID, wantProject, wantCwd string, wantInput, wantOutput, wantCache int64) {
	t.Helper()

	if got.ID != wantID {
		t.Errorf("ID = %q, want %q", got.ID, wantID)
	}
	if got.Project != wantProject {
		t.Errorf("Project = %q, want %q", got.Project, wantProject)
	}
	if got.Cwd != wantCwd {
		t.Errorf("Cwd = %q, want %q", got.Cwd, wantCwd)
	}
	if got.Date != "2026-06-01" {
		t.Errorf("Date = %q, want 2026-06-01", got.Date)
	}
	if got.InputTokens != wantInput {
		t.Errorf("InputTokens = %d, want %d", got.InputTokens, wantInput)
	}
	if got.OutputTokens != wantOutput {
		t.Errorf("OutputTokens = %d, want %d", got.OutputTokens, wantOutput)
	}
	if got.CacheTokens != wantCache {
		t.Errorf("CacheTokens = %d, want %d", got.CacheTokens, wantCache)
	}
}

func copyFixtureToProjectPath(t *testing.T, fixture, encodedProject string) string {
	t.Helper()
	return copyFixtureToProjectPathInDir(t, t.TempDir(), fixture, encodedProject)
}

func copyFixtureToProjectPathInDir(t *testing.T, dir, fixture, encodedProject string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "testdata", "claude", fixture))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixture, err)
	}

	projectDir := filepath.Join(dir, "projects", encodedProject)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	path := filepath.Join(projectDir, fixture)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return path
}
