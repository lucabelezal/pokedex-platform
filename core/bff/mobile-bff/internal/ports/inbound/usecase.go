package inbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// PokemonUseCase define casos de uso para operacoes de Pokemon
type PokemonUseCase interface {
	ListPokemons(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error)
	GetPokemonDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonDetail, error)
	GetPokemonScreenDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonScreenDetail, error)
	SearchPokemons(ctx context.Context, query string, page, pageSize int, userID string) (*domain.PokemonPage, error)
	FilterByType(ctx context.Context, typeFilter string, page, pageSize int, userID string) (*domain.PokemonPage, error)
	GetHomeData(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error)
	ListTypes(ctx context.Context) ([]domain.Type, error)
	ListRegions(ctx context.Context) ([]domain.Region, error)
}

// FavoriteUseCase define casos de uso para operacoes de Favorito
type FavoriteUseCase interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}

// AuthUseCase define casos de uso para autenticacao e sessao.
type AuthUseCase interface {
	Signup(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Login(ctx context.Context, email, password string) (*domain.AuthSession, error)
	Refresh(ctx context.Context, token string) (*domain.AuthSession, error)
	Logout(ctx context.Context, token string) error
}
