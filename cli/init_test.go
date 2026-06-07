package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mizanmahi/aiusage/cli/internal/config"
)

func TestRunInitSavesConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".aiusage", "config.toml")
	t.Setenv("AIUSAGE_CONFIG_PATH", path)

	var out bytes.Buffer
	input := strings.NewReader("http://localhost:8080\nak_secret_value\n")

	if err := runInit(input, &out); err != nil {
		t.Fatalf("runInit() error = %v", err)
	}

	if strings.Contains(out.String(), "ak_secret_value") {
		t.Fatal("runInit() printed the API key")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if gotMode := info.Mode().Perm(); gotMode != 0600 {
		t.Fatalf("config file mode = %o, want 0600", gotMode)
	}

	got, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	if got.ServerURL != "http://localhost:8080" {
		t.Errorf("ServerURL = %q, want http://localhost:8080", got.ServerURL)
	}
	if got.APIKey != "ak_secret_value" {
		t.Errorf("APIKey = %q, want saved API key", got.APIKey)
	}
}

func TestRunInitRequiresValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "missing server URL", input: "\nak_secret_value\n"},
		{name: "missing API key", input: "http://localhost:8080\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AIUSAGE_CONFIG_PATH", filepath.Join(t.TempDir(), "config.toml"))

			var out bytes.Buffer
			if err := runInit(strings.NewReader(tt.input), &out); err == nil {
				t.Fatal("runInit() error = nil, want error")
			}
		})
	}
}
