package outbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// AuthProvider define as operacoes de autenticacao outbound necessarias pelo BFF.
type AuthProvider interface {
	Signup(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Login(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Refresh(ctx context.Context, token string) (*domain.AuthSession, error)
	Logout(ctx context.Context, token string) error
}
