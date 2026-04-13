// Package config loads the zellij-tui TOML configuration file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config holds all user-configurable settings.
type Config struct {
	// Path to the zellij binary. If empty, uses PATH lookup.
	// SECURITY: this path is passed directly to exec.Command/syscall.Exec
	// with full user privileges. Only set this if you control the target
	// binary and the config file's write permissions.
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
//  1. $XDG_CONFIG_HOME/zellij-tui/config.toml
//  2. ~/.config/zellij-tui/config.toml
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

	if err := toml.Unmarshal(data, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config parse error in %s: %v\n", path, err)
	}
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
