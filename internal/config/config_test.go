package config_test

import (
	"testing"

	"next-video-golang/internal/config"
)

func TestLoadUsesDevelopmentDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := config.Load()

	if cfg.Port != "8080" {
		t.Fatalf("Port = %q, want 8080", cfg.Port)
	}

	wantDatabaseURL := "postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable"
	if cfg.DatabaseURL != wantDatabaseURL {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, wantDatabaseURL)
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/custom?sslmode=disable")

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want 9090", cfg.Port)
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/custom?sslmode=disable" {
		t.Fatalf("DatabaseURL = %q", cfg.DatabaseURL)
	}
}
