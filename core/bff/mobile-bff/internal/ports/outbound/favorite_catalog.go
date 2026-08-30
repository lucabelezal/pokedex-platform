package outbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// FavoriteCatalogProvider abstrai operações de favoritos via catalog-service.
type FavoriteCatalogProvider interface {
	GetFavoriteDetails(ctx context.Context, ids []string) ([]domain.Pokemon, error)
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
}
