package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("AUTH_SERVICE_PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("REFRESH_TOKEN_TTL_HOURS", "")
	t.Setenv("AUTH_CLEANUP_INTERVAL_MINS", "")

	cfg := Load()
	if cfg.Port != "8082" {
		t.Errorf("Port = %q, want 8082", cfg.Port)
	}
	if cfg.TokenTTLmins != 15 {
		t.Errorf("TokenTTLmins = %d, want 15", cfg.TokenTTLmins)
	}
	if cfg.RefreshTokenTTLHours != 168 {
		t.Errorf("RefreshTokenTTLHours = %d, want 168", cfg.RefreshTokenTTLHours)
	}
	if cfg.CleanupIntervalMins != 30 {
		t.Errorf("CleanupIntervalMins = %d, want 30", cfg.CleanupIntervalMins)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("AUTH_SERVICE_PORT", "9999")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/auth")
	t.Setenv("JWT_SECRET", "segredo")
	t.Setenv("REFRESH_TOKEN_TTL_HOURS", "24")
	t.Setenv("AUTH_CLEANUP_INTERVAL_MINS", "5")

	cfg := Load()
	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want 9999", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/auth" {
		t.Errorf("DatabaseURL inesperada")
	}
	if cfg.JWTSecret != "segredo" {
		t.Errorf("JWTSecret inesperado")
	}
	if cfg.RefreshTokenTTLHours != 24 {
		t.Errorf("RefreshTokenTTLHours = %d, want 24", cfg.RefreshTokenTTLHours)
	}
	if cfg.CleanupIntervalMins != 5 {
		t.Errorf("CleanupIntervalMins = %d, want 5", cfg.CleanupIntervalMins)
	}
}

func TestLoadInvalidEnvFallsBack(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_TTL_HOURS", "abc")
	t.Setenv("AUTH_CLEANUP_INTERVAL_MINS", "-5")

	cfg := Load()
	if cfg.RefreshTokenTTLHours != 168 {
		t.Errorf("RefreshTokenTTLHours = %d, want 168 (fallback)", cfg.RefreshTokenTTLHours)
	}
	if cfg.CleanupIntervalMins != 30 {
		t.Errorf("CleanupIntervalMins = %d, want 30 (fallback)", cfg.CleanupIntervalMins)
	}
}
