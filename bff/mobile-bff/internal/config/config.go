package config

import (
	"os"
)

type Config struct {
	Port              string
	DatabaseURL       string
	PokedexServiceURL string
	RedisURL          string
}

func LoadConfig() *Config {
	port := os.Getenv("MOBILE_BFF_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:              port,
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		PokedexServiceURL: os.Getenv("POKEDEX_SERVICE_URL"),
		RedisURL:          os.Getenv("REDIS_URL"),
	}
}
