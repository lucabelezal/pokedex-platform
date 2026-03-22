package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"pokedex-platform/app/pokedex-service/internal/config"
	apphttp "pokedex-platform/app/pokedex-service/internal/http"
	"pokedex-platform/app/pokedex-service/internal/repository"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pokemonRepo repository.PokemonRepository
	if cfg.DatabaseURL != "" {
		pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("falha ao conectar no postgres, usando fallback em memoria: %v", err)
			pokemonRepo = repository.NewInMemoryPokemonRepository()
		} else {
			defer pool.Close()
			pokemonRepo = repository.NewPostgresPokemonRepository(pool)
		}
	} else {
		log.Println("DATABASE_URL ausente, usando fallback em memoria")
		pokemonRepo = repository.NewInMemoryPokemonRepository()
	}

	mux := apphttp.NewMux(pokemonRepo)

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
