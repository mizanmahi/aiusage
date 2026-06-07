package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LastPushedAt time.Time `json:"last_pushed_at"`
}

func Load() (*State, error) {
	path, err := statePath()
	if err != nil {
		return &State{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return &State{}, nil
	}

	return &state, nil
}

func Save(state *State) error {
	path, err := statePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func statePath() (string, error) {
	if path := os.Getenv("AIUSAGE_STATE_PATH"); path != "" {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".aiusage", "state.json"), nil
}
