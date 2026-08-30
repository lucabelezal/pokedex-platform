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
	FavoritesViaCatalog      bool
}

func LoadConfig() *Config {
	port := strings.TrimSpace(os.Getenv("MOBILE_BFF_PORT"))
	if port == "" {
		port = "8080"
	}

	// FAVORITES_VIA_CATALOG controla se o BFF acessa favoritos via catalog-service
	// (REST) ou via PostgreSQL direto. Default: true (arquitetura hexagonal).
	favoritesViaCatalog := true
	if raw := strings.TrimSpace(os.Getenv("FAVORITES_VIA_CATALOG")); raw != "" {
		switch strings.ToLower(raw) {
		case "false", "0", "no", "off":
			favoritesViaCatalog = false
		}
	}

	return &Config{
		Port:                     port,
		DatabaseURL:              strings.TrimSpace(os.Getenv("DATABASE_URL")),
		PokemonCatalogServiceURL: strings.TrimSpace(os.Getenv("POKEMON_CATALOG_SERVICE_URL")),
		AuthServiceURL:           strings.TrimSpace(os.Getenv("AUTH_SERVICE_URL")),
		RedisURL:                 strings.TrimSpace(os.Getenv("REDIS_URL")),
		JWTSecret:                strings.TrimSpace(os.Getenv("JWT_SECRET")),
		FavoritesViaCatalog:      favoritesViaCatalog,
	}
}
