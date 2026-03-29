package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                     string
	DatabaseURL              string
	PokemonCatalogServiceURL string
	AuthServiceURL           string
	RedisURL                 string
	JWTSecret                string
}

func LoadConfig() *Config {
	port := strings.TrimSpace(os.Getenv("MOBILE_BFF_PORT"))
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:                     port,
		DatabaseURL:              strings.TrimSpace(os.Getenv("DATABASE_URL")),
		PokemonCatalogServiceURL: strings.TrimSpace(os.Getenv("POKEMON_CATALOG_SERVICE_URL")),
		AuthServiceURL:           strings.TrimSpace(os.Getenv("AUTH_SERVICE_URL")),
		RedisURL:                 strings.TrimSpace(os.Getenv("REDIS_URL")),
		JWTSecret:                strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}
}
