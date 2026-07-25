package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/config"
	apphttp "pokedex-platform/core/app/pokemon-catalog-service/internal/http"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pokemonRepo repository.PokemonRepository
	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("falha ao conectar no postgres: %v", err)
		}
		defer pool.Close()
		dbPool = pool
		pokemonRepo = repository.NewPostgresPokemonRepository(pool)
	} else {
		log.Fatal("DATABASE_URL e obrigatoria para o catalog-service")
	}

	mux := apphttp.NewMux(pokemonRepo)
	if dbPool != nil {
		mux.HandleFunc("GET /ready", apphttp.ReadyHandler(dbPool))
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("pokemon-catalog-service listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("pokemon-catalog-service server error: %v", err)
	}
}
