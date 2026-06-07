package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServerURL  string `toml:"server_url"`
	APIKey     string `toml:"api_key"`
	ClaudePath string `toml:"claude_path"`
	CodexPath  string `toml:"codex_path"`
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	applyDefaults(&config)
	return &config, nil
}

func Save(config *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(config)
}

func applyDefaults(config *Config) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	if config.ClaudePath == "" {
		config.ClaudePath = filepath.Join(home, ".claude")
	}
	if config.CodexPath == "" {
		config.CodexPath = filepath.Join(home, ".codex")
	}
}

func configPath() (string, error) {
	if path := os.Getenv("AIUSAGE_CONFIG_PATH"); path != "" {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".aiusage", "config.toml"), nil
}
