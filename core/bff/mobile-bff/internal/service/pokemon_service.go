package service

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

type PokemonService struct {
	pokemonRepo  outbound.PokemonRepository
	favoriteRepo outbound.FavoriteRepository
}

func NewPokemonService(
	pokemonRepo outbound.PokemonRepository,
	favoriteRepo outbound.FavoriteRepository,
) *PokemonService {
	return &PokemonService{
		pokemonRepo:  pokemonRepo,
		favoriteRepo: favoriteRepo,
	}
}

func (s *PokemonService) ListPokemons(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error) {
	page, pageSize = validatePagination(page, pageSize)
	return s.pokemonRepo.GetAll(ctx, page, pageSize)
}

func (s *PokemonService) GetPokemonDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonDetail, error) {
	pokemon, err := s.pokemonRepo.GetByID(ctx, pokemonID)
	if err != nil {
		return nil, err
	}

	detail := &domain.PokemonDetail{
		Number:      pokemon.Number,
		Name:        pokemon.Name,
		ImageURL:    pokemon.ImageURL,
		Height:      pokemon.Height,
		Weight:      pokemon.Weight,
		Description: pokemon.Description,
		Element: domain.Element{
			Color: pokemon.ElementColor,
			Type:  pokemon.ElementType,
		},
		Types: convertStringTypesToDomainTypes(pokemon.Types),
	}

	if userID != "" {
		isFav, err := s.favoriteRepo.IsFavorite(ctx, userID, pokemonID)
		if err == nil {
			detail.IsFavorite = isFav
		}
	}

	return detail, nil
}

func (s *PokemonService) GetPokemonScreenDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonScreenDetail, error) {
	detail, err := s.pokemonRepo.GetDetailByID(ctx, pokemonID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (s *PokemonService) SearchPokemons(ctx context.Context, query string, page, pageSize int, userID string) (*domain.PokemonPage, error) {
	page, pageSize = validatePagination(page, pageSize)
	return s.pokemonRepo.Search(ctx, query, page, pageSize)
}

func (s *PokemonService) FilterByType(ctx context.Context, typeFilter string, page, pageSize int, userID string) (*domain.PokemonPage, error) {
	page, pageSize = validatePagination(page, pageSize)
	return s.pokemonRepo.GetByType(ctx, typeFilter, page, pageSize)
}

func (s *PokemonService) GetHomeData(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error) {
	return s.ListPokemons(ctx, page, pageSize, userID)
}

func (s *PokemonService) ListTypes(ctx context.Context) ([]domain.Type, error) {
	return s.pokemonRepo.ListTypes(ctx)
}

func (s *PokemonService) ListRegions(ctx context.Context) ([]domain.Region, error) {
	return s.pokemonRepo.ListRegions(ctx)
}

func validatePagination(page, pageSize int) (int, int) {
	if page < 0 {
		page = 0
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func convertStringTypesToDomainTypes(types []string) []domain.Type {
	domainTypes := make([]domain.Type, len(types))
	for i, t := range types {
		domainTypes[i] = domain.Type{
			Name:  t,
			Color: getTypeColor(t),
		}
	}
	return domainTypes
}

func getTypeColor(typeStr string) string {
	return domain.TypeColor(typeStr)
}

var _ inbound.PokemonUseCase = (*PokemonService)(nil)
