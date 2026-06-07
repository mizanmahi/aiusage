package main

import (
	"testing"
	"time"

	"github.com/mizanmahi/aiusage/cli/internal/config"
	"github.com/mizanmahi/aiusage/cli/internal/state"
)

func saveTestConfig(t *testing.T, claudeHome, codexHome string) {
	t.Helper()

	if err := config.Save(&config.Config{
		ServerURL:  "http://localhost:8080",
		APIKey:     "ak_secret_value",
		ClaudePath: claudeHome,
		CodexPath:  codexHome,
	}); err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}
}

func saveTestState(t *testing.T, lastPushedAt time.Time) {
	t.Helper()

	if err := state.Save(&state.State{LastPushedAt: lastPushedAt}); err != nil {
		t.Fatalf("state.Save() error = %v", err)
	}
}
