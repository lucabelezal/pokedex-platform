package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	TokenTTLmins int
}

func Load() Config {
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	return Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		JWTSecret:    jwtSecret,
		TokenTTLmins: 15,
	}
}
