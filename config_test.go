package rlstudy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Port != 21320 {
		t.Errorf("expected port 21320, got %d", cfg.Port)
	}
	if cfg.Arms != 10 {
		t.Errorf("expected 10 arms, got %d", cfg.Arms)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(`{"port":9999,"arms":5}`), 0644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.Port)
	}
	if cfg.Arms != 5 {
		t.Errorf("expected 5 arms, got %d", cfg.Arms)
	}
}
