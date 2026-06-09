package codex

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mizanmahi/aiusage/types"
)

func TestParseSessionFile(t *testing.T) {
	tests := []struct {
		name           string
		fixture        string
		wantID         string
		wantProject    string
		wantInput      int64
		wantOutput     int64
		wantCache      int64
		wantReasoning  int64
		wantModel      string
		wantUsageEvent bool
	}{
		{
			name:           "sums last token usage deltas",
			fixture:        "last-token-usage.jsonl",
			wantID:         "last-token-session",
			wantProject:    "myproject",
			wantInput:      1500,
			wantOutput:     300,
			wantCache:      1200,
			wantReasoning:  70,
			wantModel:      "gpt-5.5",
			wantUsageEvent: true,
		},
		{
			name:          "keeps latest cumulative total token usage",
			fixture:       "total-token-usage.jsonl",
			wantID:        "total-token-session",
			wantProject:   "cumulative",
			wantInput:     1500,
			wantOutput:    300,
			wantCache:     1200,
			wantReasoning: 70,
		},
		{
			name:          "sums older direct token fields",
			fixture:       "direct-token-fields.jsonl",
			wantID:        "direct-token-session",
			wantProject:   "direct",
			wantInput:     1500,
			wantOutput:    300,
			wantCache:     1200,
			wantReasoning: 70,
			wantModel:     "gpt-5.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := copyFixtureToSessionPath(t, tt.fixture)

			session, err := parseSessionFile(path)
			if err != nil {
				t.Fatalf("parseSessionFile() error = %v", err)
			}
			if session == nil {
				t.Fatal("parseSessionFile() returned nil")
			}

			assertSession(t, session, tt.wantID, tt.wantProject, tt.wantInput, tt.wantOutput, tt.wantCache, tt.wantReasoning, tt.wantModel)

			if tt.wantUsageEvent {
				event := session.ToUsageEvent("user-123")
				if event.Tool != types.ToolCodex {
					t.Errorf("ToUsageEvent().Tool = %q, want %q", event.Tool, types.ToolCodex)
				}
				if event.UserID != "user-123" {
					t.Errorf("ToUsageEvent().UserID = %q, want user-123", event.UserID)
				}
				if event.Project != tt.wantProject {
					t.Errorf("ToUsageEvent().Project = %q, want %q", event.Project, tt.wantProject)
				}
			}
		})
	}
}

func TestReadSessionsSinceFilter(t *testing.T) {
	dir := t.TempDir()
	path := copyFixtureToSessionPathInDir(t, dir, "last-token-usage.jsonl")

	past := time.Now().Add(-24 * time.Hour)
	sessions, err := ReadSessions(dir, past)
	if err != nil {
		t.Fatalf("ReadSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("ReadSessions() got %d sessions, want 1", len(sessions))
	}
	if sessions[0].ID != "last-token-session" {
		t.Errorf("ReadSessions()[0].ID = %q, want last-token-session", sessions[0].ID)
	}

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

func assertSession(t *testing.T, got *Session, wantID, wantProject string, wantInput, wantOutput, wantCache, wantReasoning int64, wantModel string) {
	t.Helper()

	if got.ID != wantID {
		t.Errorf("ID = %q, want %q", got.ID, wantID)
	}
	if got.Project != wantProject {
		t.Errorf("Project = %q, want %q", got.Project, wantProject)
	}
	if got.Date != "2026-06-03" {
		t.Errorf("Date = %q, want 2026-06-03", got.Date)
	}
	if got.InputTokens != wantInput {
		t.Errorf("InputTokens = %d, want %d", got.InputTokens, wantInput)
	}
	if got.OutputTokens != wantOutput {
		t.Errorf("OutputTokens = %d, want %d", got.OutputTokens, wantOutput)
	}
	if got.CacheReadTokens != wantCache {
		t.Errorf("CacheReadTokens = %d, want %d", got.CacheReadTokens, wantCache)
	}
	if got.ReasoningTokens != wantReasoning {
		t.Errorf("ReasoningTokens = %d, want %d", got.ReasoningTokens, wantReasoning)
	}
	if got.Model != wantModel {
		t.Errorf("Model = %q, want %q", got.Model, wantModel)
	}
}

func copyFixtureToSessionPath(t *testing.T, fixture string) string {
	t.Helper()
	return copyFixtureToSessionPathInDir(t, t.TempDir(), fixture)
}

func copyFixtureToSessionPathInDir(t *testing.T, dir, fixture string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "testdata", "codex", fixture))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixture, err)
	}

	sessionDir := filepath.Join(dir, "sessions", "2026", "06", "03")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	path := filepath.Join(sessionDir, fixture)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return path
}
