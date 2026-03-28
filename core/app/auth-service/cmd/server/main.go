package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"pokedex-platform/core/app/auth-service/internal/config"
	authhttp "pokedex-platform/core/app/auth-service/internal/http"
	"pokedex-platform/core/app/auth-service/internal/repository"
	"pokedex-platform/core/app/auth-service/internal/service"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL nao configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("falha ao conectar no banco: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenTTLmins)

	mux := authhttp.NewMux(authService)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("auth-service listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("auth-service server error: %v", err)
	}
}
