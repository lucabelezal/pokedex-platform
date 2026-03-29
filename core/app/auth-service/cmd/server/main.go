package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"pokedex-platform/core/app/auth-service/internal/config"
	authhttp "pokedex-platform/core/app/auth-service/internal/http"
	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"
)

func main() {
	cfg := config.Load()
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		log.Fatal("DATABASE_URL nao configurada")
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		log.Fatal("JWT_SECRET nao configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("falha ao conectar no banco: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenTTLmins, cfg.RefreshTokenTTLHours)

	startAuthCleanupJob(userRepo, cfg.CleanupIntervalMins)

	mux := authhttp.NewMux(authService)
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("auth-service listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("auth-service server error: %v", err)
	}
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
			log.Printf("auth cleanup falhou: %v", err)
			return
		}

		log.Printf("auth cleanup executado com sucesso")
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
