package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if !cfg.AutoAttach {
		t.Error("AutoAttach should default to true")
	}
	if cfg.ZellijPath != "" {
		t.Errorf("ZellijPath should default to empty, got %q", cfg.ZellijPath)
	}
}

func TestLoadReturnsConfig(t *testing.T) {
	// Load should always return a non-nil config, regardless of whether
	// a config file exists on disk.
	cfg := Load()
	if cfg == nil {
		t.Fatal("Load returned nil")
	}
}
