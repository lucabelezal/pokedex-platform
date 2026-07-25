package service

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

// FavoriteCatalogProvider abstrai operações de favoritos via catalog-service.
type FavoriteCatalogProvider interface {
	GetFavoriteDetails(ctx context.Context, ids []string) ([]domain.Pokemon, error)
	AddFavorite(ctx context.Context, userID, pokemonID string) error
	RemoveFavorite(ctx context.Context, userID, pokemonID string) error
}

type FavoriteService struct {
	favoriteRepo    outbound.FavoriteRepository
	pokemonRepo     outbound.PokemonRepository
	favoriteCatalog FavoriteCatalogProvider
}

func NewFavoriteService(
	favoriteRepo outbound.FavoriteRepository,
	pokemonRepo outbound.PokemonRepository,
	favoriteCatalog FavoriteCatalogProvider,
) *FavoriteService {
	return &FavoriteService{
		favoriteRepo:    favoriteRepo,
		pokemonRepo:     pokemonRepo,
		favoriteCatalog: favoriteCatalog,
	}
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	if _, err := s.pokemonRepo.GetByID(ctx, pokemonID); err != nil {
		return err
	}

	if s.favoriteCatalog != nil {
		return s.favoriteCatalog.AddFavorite(ctx, userID, pokemonID)
	}
	return s.favoriteRepo.AddFavorite(ctx, userID, pokemonID)
}

func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	if s.favoriteCatalog != nil {
		return s.favoriteCatalog.RemoveFavorite(ctx, userID, pokemonID)
	}
	return s.favoriteRepo.RemoveFavorite(ctx, userID, pokemonID)
}

func (s *FavoriteService) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	return s.favoriteRepo.GetUserFavorites(ctx, userID)
}

func (s *FavoriteService) GetFavoriteDetails(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
	if s.favoriteCatalog == nil {
		return []domain.Pokemon{}, nil
	}
	return s.favoriteCatalog.GetFavoriteDetails(ctx, ids)
}

var _ inbound.FavoriteUseCase = (*FavoriteService)(nil)
