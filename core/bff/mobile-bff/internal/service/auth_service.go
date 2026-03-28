package service

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

type AuthService struct {
	authProvider ports.AuthProvider
}

func NewAuthService(authProvider ports.AuthProvider) *AuthService {
	return &AuthService{authProvider: authProvider}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	if s.authProvider == nil {
		return nil, domain.ErrAuthUnavailable
	}
	return s.authProvider.Signup(ctx, email, password)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	if s.authProvider == nil {
		return nil, domain.ErrAuthUnavailable
	}
	return s.authProvider.Login(ctx, email, password)
}

func (s *AuthService) Refresh(ctx context.Context, token string) (*ports.AuthSession, error) {
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

var _ ports.AuthUseCase = (*AuthService)(nil)
