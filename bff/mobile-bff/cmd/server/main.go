package main

import (
	"context"
	"log"
	"net/http"
	"time"

	httpadapter "pokedex-platform/bff/mobile-bff/internal/adapters/http"
	"pokedex-platform/bff/mobile-bff/internal/adapters/repository"
	"pokedex-platform/bff/mobile-bff/internal/config"
	"pokedex-platform/bff/mobile-bff/internal/ports"
	"pokedex-platform/bff/mobile-bff/internal/service"
	"pokedex-platform/bff/mobile-bff/tests/mocks"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize repositories with fallback to mocks
	var pokemonRepo ports.PokemonRepository
	var favoriteRepo ports.FavoriteRepository
	var db *repository.Database

	if cfg.PokedexServiceURL != "" {
		pokemonRepo = repository.NewPokedexServiceRepository(cfg.PokedexServiceURL)
		log.Printf("Using pokedex-service catalog from %s", cfg.PokedexServiceURL)
	}

	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		db, err = repository.NewDatabase(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: PostgreSQL unavailable for favorites, using mock favorites: %v", err)
			favoriteRepo = mocks.NewMockFavoriteRepository()
		} else {
			favoriteRepo = repository.NewPostgresFavoriteRepository(db.Pool)
			if pokemonRepo == nil {
				pokemonRepo = repository.NewPostgresPokemonRepository(db.Pool)
			}
		}
	}

	if db != nil {
		defer db.Close()
	}

	if pokemonRepo == nil {
		if cfg.DatabaseURL == "" {
			log.Println("No POKEDEX_SERVICE_URL or DATABASE_URL set, using mock pokemons")
		}
		pokemonRepo = mocks.NewMockPokemonRepository()
	}

	if favoriteRepo == nil {
		if cfg.DatabaseURL == "" {
			log.Println("No DATABASE_URL set, using mock favorites")
		}
		favoriteRepo = mocks.NewMockFavoriteRepository()
	}

	// Setup services
	pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	// Setup HTTP handlers
	mux := http.NewServeMux()
	h := httpadapter.NewHandler(pokemonService, favoriteService)
	h.RegisterRoutes(mux)

	// Apply middleware
	var handler http.Handler = mux
	handler = httpadapter.CORSMiddleware(handler)
	handler = httpadapter.AuthMiddleware(handler)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("mobile-bff listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
