package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config holds all user-configurable settings.
type Config struct {
	// Path to the zellij binary. If empty, uses PATH lookup.
	ZellijPath string `toml:"zellij_path"`

	// AutoAttach controls what happens after creating a new session.
	// true  = immediately attach to the new session (default)
	// false = create in background, show alert, stay in TUI
	AutoAttach bool `toml:"auto_attach"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		AutoAttach: true,
	}
}

// Load reads config from the default locations:
//  1. ~/.config/zellij-tui/config.toml
//  2. $XDG_CONFIG_HOME/zellij-tui/config.toml
//
// Returns default config if no file is found.
func Load() *Config {
	cfg := DefaultConfig()

	path := configPath()
	if path == "" {
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	_ = toml.Unmarshal(data, cfg)
	return cfg
}

func configPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		p := filepath.Join(xdg, "zellij-tui", "config.toml")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	p := filepath.Join(home, ".config", "zellij-tui", "config.toml")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}
