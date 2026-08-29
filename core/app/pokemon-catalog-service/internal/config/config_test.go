package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("POKEMON_CATALOG_SERVICE_PORT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()
	if cfg.Port != "8081" {
		t.Errorf("Port = %q, want 8081", cfg.Port)
	}
	if cfg.DatabaseURL != "" {
		t.Errorf("DatabaseURL = %q, want vazio", cfg.DatabaseURL)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("POKEMON_CATALOG_SERVICE_PORT", "9999")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg := Load()
	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want 9999", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("DatabaseURL inesperada: %q", cfg.DatabaseURL)
	}
}
