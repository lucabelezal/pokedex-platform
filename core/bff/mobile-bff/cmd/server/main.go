package main

import (
	"context"
	"log"
	"net/http"
	"time"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/repository"
	"pokedex-platform/core/bff/mobile-bff/internal/config"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
	"pokedex-platform/core/bff/mobile-bff/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// Inicializar repositórios com fallback para mocks
	var pokemonRepo ports.PokemonRepository
	var favoriteRepo ports.FavoriteRepository
	var db *repository.Database

	if cfg.PokemonCatalogServiceURL != "" {
		pokemonRepo = repository.NewPokemonCatalogServiceRepository(cfg.PokemonCatalogServiceURL)
		log.Printf("Using pokemon-catalog-service catalog from %s", cfg.PokemonCatalogServiceURL)
	}

	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		db, err = repository.NewDatabase(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: PostgreSQL unavailable for favorites, using mock favorites: %v", err)
			favoriteRepo = repository.NewMockFavoriteRepository()
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
			log.Println("No POKEMON_CATALOG_SERVICE_URL or DATABASE_URL set, using mock pokemons")
		}
		pokemonRepo = repository.NewMockPokemonRepository()
	}

	if favoriteRepo == nil {
		if cfg.DatabaseURL == "" {
			log.Println("No DATABASE_URL set, using mock favorites")
		}
		favoriteRepo = repository.NewMockFavoriteRepository()
	}

	// Configurar serviços
	pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	// Configurar cliente de auth-service
	authClient := repository.NewAuthServiceClient(cfg.AuthServiceURL)
	authService := service.NewAuthService(authClient)

	// Configurar handlers HTTP
	mux := http.NewServeMux()
	h := httpadapter.NewHandler(pokemonService, favoriteService, authService)
	h.RegisterRoutes(mux)

	// Aplicar middleware
	var handler http.Handler = mux
	handler = httpadapter.CORSMiddleware(handler)
	handler = httpadapter.AuthMiddleware(handler)

	// Iniciar servidor
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
