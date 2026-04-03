package inbound

import "context"

// FavoriteUseCase define casos de uso para operacoes de Favorito.
type FavoriteUseCase interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}
