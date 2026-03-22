package ports

import (
	"context"

	"pokedex-platform/bff/mobile-bff/internal/domain"
)

// PokemonRepository defines the contract for accessing Pokémon data
type PokemonRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
	GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error)
}

// FavoriteRepository defines the contract for accessing Favorite data
type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error)
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}
