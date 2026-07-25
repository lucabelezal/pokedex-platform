package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"pokedex-platform/core/app/auth-service/internal/config"
	authhttp "pokedex-platform/core/app/auth-service/internal/http"
	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"
)

func main() {
	setupLogger()
	cfg := config.Load()
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		slog.Error("DATABASE_URL nao configurada")
		os.Exit(1)
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		slog.Error("JWT_SECRET nao configurada")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("falha ao conectar no banco", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenTTLmins, cfg.RefreshTokenTTLHours)

	startAuthCleanupJob(userRepo, cfg.CleanupIntervalMins)

	mux := authhttp.NewMux(authService)
	mux.HandleFunc("GET /ready", authhttp.ReadyHandler(pool))
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	slog.Info("servidor iniciado", "addr", srv.Addr, "service", "auth-service")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("servidor encerrado com erro", "error", err)
		os.Exit(1)
	}
}

func setupLogger() {
	level := slog.LevelInfo
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		switch v {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	format := os.Getenv("LOG_FORMAT")
	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(handler))
}

func startAuthCleanupJob(userRepo *repository.UserRepository, intervalMins int) {
	if userRepo == nil {
		return
	}

	interval := time.Duration(intervalMins) * time.Minute
	if interval <= 0 {
		interval = 30 * time.Minute
	}

	runCleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := userRepo.CleanupExpiredAuthArtifacts(ctx); err != nil {
			slog.Warn("auth cleanup falhou", "error", err)
			return
		}

		slog.Debug("auth cleanup executado com sucesso")
	}

	runCleanup()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			runCleanup()
		}
	}()
}
