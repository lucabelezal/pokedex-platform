package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http"
	httpclient "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
	"pokedex-platform/core/bff/mobile-bff/internal/config"
	applogger "pokedex-platform/core/bff/mobile-bff/internal/infrastructure/logger"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
	"pokedex-platform/core/bff/mobile-bff/internal/service"
	"pokedex-platform/core/bff/mobile-bff/tests/mocks"
)

func main() {
	applogger.Setup("mobile-bff")

	cfg := config.LoadConfig()

	// Inicializar repositórios com fallback para mocks
	var pokemonRepo outbound.PokemonRepository
	var favoriteRepo outbound.FavoriteRepository
	var db *postgres.Database

	if strings.TrimSpace(cfg.PokemonCatalogServiceURL) == "" {
		slog.Error("configuracao invalida", "motivo", "POKEMON_CATALOG_SERVICE_URL obrigatoria")
		os.Exit(1)
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		slog.Error("configuracao invalida", "motivo", "JWT_SECRET obrigatoria")
		os.Exit(1)
	}

	pokemonRepo = httpclient.NewPokemonCatalogServiceRepository(cfg.PokemonCatalogServiceURL)
	slog.Info("pokemon catalog configurado", "url", cfg.PokemonCatalogServiceURL)

	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		db, err = postgres.NewDatabase(ctx, cfg.DatabaseURL)
		if err != nil {
			slog.Warn("postgres indisponivel, usando mock de favoritos", "error", err)
			favoriteRepo = mocks.NewMockFavoriteRepository()
		} else {
			favoriteRepo = postgres.NewPostgresFavoriteRepository(db.Pool)
		}
	}

	if db != nil {
		defer db.Close()
	}

	if favoriteRepo == nil {
		if cfg.DatabaseURL == "" {
			slog.Info("database_url nao configurada, usando mock de favoritos")
		}
		favoriteRepo = mocks.NewMockFavoriteRepository()
	}

	// Configurar serviços
	pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	// Configurar cliente de auth-service
	authClient := httpclient.NewAuthServiceClient(cfg.AuthServiceURL)
	authService := service.NewAuthService(authClient)

	// Configurar handlers HTTP
	mux := http.NewServeMux()
	h := httpadapter.NewHandler(pokemonService, favoriteService, authService)
	h.RegisterRoutes(mux)

	// Aplicar middleware
	var handler http.Handler = mux
	handler = httpadapter.CORSMiddleware(handler)
	handler = httpadapter.AuthRateLimitMiddleware(handler)
	handler = httpadapter.AuthMiddleware(authClient, handler)
	handler = httpadapter.RequestLoggerMiddleware(handler)

	// Iniciar servidor
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	slog.Info("servidor iniciado", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("servidor encerrado com erro", "error", err)
		os.Exit(1)
	}
}
