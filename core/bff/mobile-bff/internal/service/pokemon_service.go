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
	typeColors := map[string]string{
		"Normal":    "#A8A878",
		"Fogo":      "#EE8130",
		"Fire":      "#F08030",
		"Água":      "#6390F0",
		"Water":     "#6890F0",
		"Elétrico":  "#F7D02C",
		"Electric":  "#F8D030",
		"Grama":     "#7AC74C",
		"Grass":     "#78C850",
		"Gelo":      "#96D9D6",
		"Ice":       "#98D8D8",
		"Lutador":   "#C22E28",
		"Fighting":  "#C03028",
		"Venenoso":  "#A33EA1",
		"Poison":    "#A040A0",
		"Terrestre": "#E2BF65",
		"Ground":    "#E0C068",
		"Voador":    "#A98FF3",
		"Flying":    "#A890F0",
		"Psíquico":  "#F95587",
		"Psychic":   "#F85888",
		"Inseto":    "#A6B91A",
		"Bug":       "#A8B820",
		"Pedra":     "#B6A136",
		"Rock":      "#B8A038",
		"Fantasma":  "#735797",
		"Ghost":     "#705898",
		"Dragão":    "#6F35FC",
		"Dragon":    "#7038F8",
		"Sombrio":   "#705746",
		"Dark":      "#705848",
		"Aço":       "#B7B7CE",
		"Steel":     "#B8B8D0",
		"Fada":      "#D685AD",
		"Fairy":     "#EE99AC",
	}

	if color, exists := typeColors[typeStr]; exists {
		return color
	}
	return "#A9AC86"
}

var _ inbound.PokemonUseCase = (*PokemonService)(nil)
