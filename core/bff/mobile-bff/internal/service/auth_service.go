package service

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

type AuthService struct {
	authProvider outbound.AuthProvider
}

func NewAuthService(authProvider outbound.AuthProvider) *AuthService {
	return &AuthService{authProvider: authProvider}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*domain.AuthSession, error) {
	if s.authProvider == nil {
		return nil, domain.ErrAuthUnavailable
	}
	return s.authProvider.Signup(ctx, email, password)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.AuthSession, error) {
	if s.authProvider == nil {
		return nil, domain.ErrAuthUnavailable
	}
	return s.authProvider.Login(ctx, email, password)
}

func (s *AuthService) Refresh(ctx context.Context, token string) (*domain.AuthSession, error) {
	if s.authProvider == nil {
		return nil, domain.ErrAuthUnavailable
	}
	return s.authProvider.Refresh(ctx, token)
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if s.authProvider == nil {
		return domain.ErrAuthUnavailable
	}
	return s.authProvider.Logout(ctx, token)
}

var _ inbound.AuthUseCase = (*AuthService)(nil)
