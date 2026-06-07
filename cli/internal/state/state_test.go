package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMissingState(t *testing.T) {
	t.Setenv("AIUSAGE_STATE_PATH", filepath.Join(t.TempDir(), "state.json"))

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.LastPushedAt.IsZero() != true {
		t.Fatalf("Load().LastPushedAt = %v, want zero", got.LastPushedAt)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".aiusage", "state.json")
	t.Setenv("AIUSAGE_STATE_PATH", path)

	wantTime := time.Date(2026, 6, 8, 12, 30, 0, 0, time.UTC)
	if err := Save(&State{LastPushedAt: wantTime}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if gotMode := info.Mode().Perm(); gotMode != 0600 {
		t.Fatalf("state file mode = %o, want 0600", gotMode)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !got.LastPushedAt.Equal(wantTime) {
		t.Fatalf("Load().LastPushedAt = %v, want %v", got.LastPushedAt, wantTime)
	}
}

func TestLoadCorruptStateReturnsEmptyState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	t.Setenv("AIUSAGE_STATE_PATH", path)

	if err := os.WriteFile(path, []byte("{not-json"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !got.LastPushedAt.IsZero() {
		t.Fatalf("Load().LastPushedAt = %v, want zero", got.LastPushedAt)
	}
}

func TestStatePathUsesOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom-state.json")
	t.Setenv("AIUSAGE_STATE_PATH", path)

	got, err := statePath()
	if err != nil {
		t.Fatalf("statePath() error = %v", err)
	}
	if got != path {
		t.Fatalf("statePath() = %q, want %q", got, path)
	}
}
