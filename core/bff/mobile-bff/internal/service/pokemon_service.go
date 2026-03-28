package service

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

type PokemonService struct {
	pokemonRepo  ports.PokemonRepository
	favoriteRepo ports.FavoriteRepository
}

func NewPokemonService(
	pokemonRepo ports.PokemonRepository,
	favoriteRepo ports.FavoriteRepository,
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

type FavoriteService struct {
	favoriteRepo ports.FavoriteRepository
	pokemonRepo  ports.PokemonRepository
}

func NewFavoriteService(
	favoriteRepo ports.FavoriteRepository,
	pokemonRepo ports.PokemonRepository,
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
		"Normal":   "#A8A878",
		"Fire":     "#F08030",
		"Water":    "#6890F0",
		"Electric": "#F8D030",
		"Grass":    "#78C850",
		"Ice":      "#98D8D8",
		"Fighting": "#C03028",
		"Poison":   "#A040A0",
		"Ground":   "#E0C068",
		"Flying":   "#A890F0",
		"Psychic":  "#F85888",
		"Bug":      "#A8B820",
		"Rock":     "#B8A038",
		"Ghost":    "#705898",
		"Dragon":   "#7038F8",
		"Dark":     "#705848",
		"Steel":    "#B8B8D0",
		"Fairy":    "#EE99AC",
	}

	if color, exists := typeColors[typeStr]; exists {
		return color
	}
	return "#A9AC86"
}
