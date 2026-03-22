package main

import (
	"log"
	"net/http"
	"time"

	"pokedex-platform/app/pokedex-service/internal/config"
	apphttp "pokedex-platform/app/pokedex-service/internal/http"
)

func main() {
	cfg := config.Load()
	mux := apphttp.NewMux()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("pokedex-service listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("pokedex-service server error: %v", err)
	}
}
