package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	TokenTTLmins         int
	RefreshTokenTTLHours int
	CleanupIntervalMins  int
}

func Load() Config {
	port := strings.TrimSpace(os.Getenv("AUTH_SERVICE_PORT"))
	if port == "" {
		port = "8082"
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	refreshTokenTTLHours := 168
	if raw := strings.TrimSpace(os.Getenv("REFRESH_TOKEN_TTL_HOURS")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			refreshTokenTTLHours = parsed
		}
	}

	cleanupIntervalMins := 30
	if raw := strings.TrimSpace(os.Getenv("AUTH_CLEANUP_INTERVAL_MINS")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			cleanupIntervalMins = parsed
		}
	}

	return Config{
		Port:                 port,
		DatabaseURL:          databaseURL,
		JWTSecret:            jwtSecret,
		TokenTTLmins:         15,
		RefreshTokenTTLHours: refreshTokenTTLHours,
		CleanupIntervalMins:  cleanupIntervalMins,
	}
}
