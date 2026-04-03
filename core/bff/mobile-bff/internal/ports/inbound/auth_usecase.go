package inbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// AuthUseCase define casos de uso para autenticacao e sessao.
type AuthUseCase interface {
	Signup(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Login(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Refresh(ctx context.Context, token string) (*domain.AuthSession, error)
	Logout(ctx context.Context, token string) error
}
