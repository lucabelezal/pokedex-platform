package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/config"
	apphttp "pokedex-platform/core/app/pokemon-catalog-service/internal/http"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/observability"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	tp, err := observability.InitTracer(context.Background(), "pokemon-catalog-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		slog.Warn("falha ao inicializar tracing", "error", err)
	}
	defer func() {
		if tp != nil {
			_ = observability.ShutdownTracer(context.Background(), tp)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pokemonRepo repository.PokemonRepository
	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			slog.Error("falha ao conectar no postgres", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		dbPool = pool
		pokemonRepo = repository.NewPostgresPokemonRepository(pool)
	} else {
		slog.Error("DATABASE_URL e obrigatoria para o catalog-service")
		os.Exit(1)
	}

	mux := apphttp.NewMux(pokemonRepo)
	if dbPool != nil {
		mux.HandleFunc("GET /ready", apphttp.ReadyHandler(dbPool))
	}
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	var handler http.Handler = mux
	handler = observability.TracingMiddleware("pokemon-catalog-service")(handler)
	handler = observability.MetricsMiddleware(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	slog.Info("servidor iniciado", "addr", srv.Addr, "service", "pokemon-catalog-service")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("servidor encerrado com erro", "error", err)
		os.Exit(1)
	}
}
