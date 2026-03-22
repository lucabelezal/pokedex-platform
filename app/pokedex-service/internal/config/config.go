package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	port := os.Getenv("POKEDEX_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	return Config{Port: port}
}
