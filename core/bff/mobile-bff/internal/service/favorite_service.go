package service

import (
	"context"

	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

type FavoriteService struct {
	favoriteRepo outbound.FavoriteRepository
	pokemonRepo  outbound.PokemonRepository
}

func NewFavoriteService(
	favoriteRepo outbound.FavoriteRepository,
	pokemonRepo outbound.PokemonRepository,
) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		pokemonRepo:  pokemonRepo,
	}
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	if _, err := s.pokemonRepo.GetByID(ctx, pokemonID); err != nil {
		return err
	}

	return s.favoriteRepo.AddFavorite(ctx, userID, pokemonID)
}

func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	return s.favoriteRepo.RemoveFavorite(ctx, userID, pokemonID)
}

func (s *FavoriteService) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	return s.favoriteRepo.GetUserFavorites(ctx, userID)
}

var _ inbound.FavoriteUseCase = (*FavoriteService)(nil)
