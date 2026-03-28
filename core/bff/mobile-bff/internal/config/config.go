package config

import (
	"os"
)

type Config struct {
	Port                     string
	DatabaseURL              string
	PokemonCatalogServiceURL string
	AuthServiceURL           string
	RedisURL                 string
}

func LoadConfig() *Config {
	port := os.Getenv("MOBILE_BFF_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:                     port,
		DatabaseURL:              os.Getenv("DATABASE_URL"),
		PokemonCatalogServiceURL: os.Getenv("POKEMON_CATALOG_SERVICE_URL"),
		AuthServiceURL:           os.Getenv("AUTH_SERVICE_URL"),
		RedisURL:                 os.Getenv("REDIS_URL"),
	}
}
