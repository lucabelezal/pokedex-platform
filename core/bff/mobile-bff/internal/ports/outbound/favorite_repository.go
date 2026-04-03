package outbound

import "context"

// FavoriteRepository define o contrato para acessar dados de Favorito.
type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
	IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error)
	GetUserFavorites(ctx context.Context, userID string) ([]string, error)
}
