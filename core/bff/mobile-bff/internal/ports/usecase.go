package ports

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// PokemonUseCase define casos de uso para operações de Pokémon
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

// FavoriteUseCase define casos de uso para operações de Favorito
type FavoriteUseCase interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}

// AuthUseCase define casos de uso para autenticação e sessão.
type AuthUseCase interface {
	Signup(ctx context.Context, email, password string) (*AuthSession, error)
	Login(ctx context.Context, email, password string) (*AuthSession, error)
	Refresh(ctx context.Context, token string) (*AuthSession, error)
	Logout(ctx context.Context, token string) error
}
