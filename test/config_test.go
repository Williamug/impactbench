package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/williamug/impactbench/internal/config"
)

func TestConfigLoad(t *testing.T) {
	// Create mock project config
	tmpDir, err := os.MkdirTemp("", "impactbench-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	dotImpact := filepath.Join(tmpDir, ".impactbench")
	os.Mkdir(dotImpact, 0755)

	configContent := `
base_url: http://localhost:8080
framework: django
thresholds:
  response_time: 15.0
`
	os.WriteFile(filepath.Join(dotImpact, "config.yaml"), []byte(configContent), 0644)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected BaseURL http://localhost:8080, got %s", cfg.BaseURL)
	}
	if cfg.Framework != "django" {
		t.Errorf("Expected framework django, got %s", cfg.Framework)
	}
	if cfg.Thresholds.ResponseTime != 15.0 {
		t.Errorf("Expected threshold 15.0, got %f", cfg.Thresholds.ResponseTime)
	}
}
