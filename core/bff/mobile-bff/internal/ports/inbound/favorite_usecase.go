package inbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// FavoriteUseCase define casos de uso para operacoes de Favorito.
type FavoriteUseCase interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
	GetFavoriteDetails(ctx context.Context, ids []string) ([]domain.Pokemon, error)
}
