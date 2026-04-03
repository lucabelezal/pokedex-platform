package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/repository"
	"pokedex-platform/core/bff/mobile-bff/internal/config"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
	"pokedex-platform/core/bff/mobile-bff/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// Inicializar repositórios com fallback para mocks
	var pokemonRepo outbound.PokemonRepository
	var favoriteRepo outbound.FavoriteRepository
	var db *repository.Database

	if strings.TrimSpace(cfg.PokemonCatalogServiceURL) == "" {
		log.Fatal("POKEMON_CATALOG_SERVICE_URL obrigatoria para iniciar o mobile-bff")
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		log.Fatal("JWT_SECRET obrigatoria para iniciar o mobile-bff")
	}

	pokemonRepo = repository.NewPokemonCatalogServiceRepository(cfg.PokemonCatalogServiceURL)
	log.Printf("Using pokemon-catalog-service catalog from %s", cfg.PokemonCatalogServiceURL)

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
		}
	}

	if db != nil {
		defer db.Close()
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
	handler = httpadapter.AuthRateLimitMiddleware(handler)
	handler = httpadapter.AuthMiddleware(handler)

	// Iniciar servidor
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("mobile-bff listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
