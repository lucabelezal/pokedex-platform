package ports

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// PokemonRepository define o contrato para acessar dados de Pokémon
type PokemonRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
	GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error)
}

// FavoriteRepository define o contrato para acessar dados de Favorito
type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error)
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}
