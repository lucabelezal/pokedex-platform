package repository

import (
	"context"
	"sync"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

type MockPokemonRepository struct {
	mu       sync.RWMutex
	pokemons map[string]*domain.Pokemon
}

func NewMockPokemonRepository() *MockPokemonRepository {
	repo := &MockPokemonRepository{
		pokemons: make(map[string]*domain.Pokemon),
	}
	repo.seedData()
	return repo
}

func (m *MockPokemonRepository) seedData() {
	m.pokemons["1"] = &domain.Pokemon{
		ID:           "1",
		Name:         "Bulbasaur",
		Number:       "001",
		Types:        []string{"Grass", "Poison"},
		Height:       0.71,
		Weight:       6.9,
		Description:  "There is a plant seed on its back right from the day this Pokemon is born.",
		ImageURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png",
		ElementColor: "#78C850",
		ElementType:  "Grass",
	}

	m.pokemons["25"] = &domain.Pokemon{
		ID:           "25",
		Name:         "Pikachu",
		Number:       "025",
		Types:        []string{"Electric"},
		Height:       0.41,
		Weight:       6.0,
		Description:  "Its body is covered in fur that stores electricity.",
		ImageURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
		ElementColor: "#F8D030",
		ElementType:  "Electric",
	}

	m.pokemons["4"] = &domain.Pokemon{
		ID:           "4",
		Name:         "Charmander",
		Number:       "004",
		Types:        []string{"Fire"},
		Height:       0.61,
		Weight:       8.5,
		Description:  "Prefers hot places. It is said that Charmander was born in volcanoes.",
		ImageURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/4.png",
		ElementColor: "#F08030",
		ElementType:  "Fire",
	}
}

func (m *MockPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pokemon, exists := m.pokemons[id]
	if !exists {
		return nil, domain.ErrPokemonNotFound
	}

	return pokemon, nil
}

func (m *MockPokemonRepository) GetDetailByID(ctx context.Context, id string) (*domain.PokemonScreenDetail, error) {
	pokemon, err := m.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category := "Seed"
	abilities := []string{"Overgrow"}
	weaknesses := []domain.Type{{Name: "Fogo", Color: "#EE8130"}}
	if pokemon.ID == "25" {
		category = "Mouse"
		abilities = []string{"Static"}
		weaknesses = []domain.Type{{Name: "Terrestre", Color: "#E2BF65"}}
	}
	if pokemon.ID == "4" {
		category = "Lizard"
		abilities = []string{"Blaze"}
		weaknesses = []domain.Type{{Name: "Água", Color: "#6390F0"}}
	}

	genderMale := 87.5
	genderFemale := 12.5

	return &domain.PokemonScreenDetail{
		ID:           pokemon.ID,
		Name:         pokemon.Name,
		Number:       pokemon.Number,
		Types:        convertMockTypes(pokemon.Types),
		Description:  pokemon.Description,
		ImageURL:     pokemon.ImageURL,
		ElementColor: pokemon.ElementColor,
		Height:       pokemon.Height,
		Weight:       pokemon.Weight,
		Category:     category,
		Abilities:    abilities,
		GenderMale:   &genderMale,
		GenderFemale: &genderFemale,
		Weaknesses:   weaknesses,
		Evolutions: []domain.Evolution{
			{ID: pokemon.ID, Number: pokemon.Number, Name: pokemon.Name, ImageURL: pokemon.ImageURL, Types: convertMockTypes(pokemon.Types)},
		},
		Region:     "Kanto",
		Generation: "1º Geração",
	}, nil
}

func (m *MockPokemonRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var content []domain.Pokemon
	for _, p := range m.pokemons {
		content = append(content, *p)
	}

	totalElements := int64(len(content))
	totalPages := (int(totalElements) + pageSize - 1) / pageSize
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       content,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

func (m *MockPokemonRepository) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var content []domain.Pokemon
	for _, p := range m.pokemons {
		if contains(p.Name, query) || contains(p.Number, query) {
			content = append(content, *p)
		}
	}

	totalElements := int64(len(content))
	totalPages := (int(totalElements) + pageSize - 1) / pageSize
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       content,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

func (m *MockPokemonRepository) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var content []domain.Pokemon
	for _, p := range m.pokemons {
		for _, t := range p.Types {
			if t == typeFilter {
				content = append(content, *p)
				break
			}
		}
	}

	totalElements := int64(len(content))
	totalPages := (int(totalElements) + pageSize - 1) / pageSize
	hasNext := page+1 < totalPages

	return &domain.PokemonPage{
		Content:       content,
		TotalElements: totalElements,
		CurrentPage:   page,
		TotalPages:    totalPages,
		HasNext:       hasNext,
	}, nil
}

func (m *MockPokemonRepository) GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error) {
	return []string{}, nil
}

func (m *MockPokemonRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	_ = ctx
	return []domain.Type{
		{Name: "Água", Color: "#6390F0"},
		{Name: "Fantasma", Color: "#735797"},
		{Name: "Grama", Color: "#7AC74C"},
		{Name: "Venenoso", Color: "#A33EA1"},
		{Name: "Fogo", Color: "#EE8130"},
		{Name: "Elétrico", Color: "#F7D02C"},
		{Name: "Voador", Color: "#A98FF3"},
	}, nil
}

func (m *MockPokemonRepository) ListRegions(ctx context.Context) ([]domain.Region, error) {
	_ = ctx
	return []domain.Region{
		{ID: "kanto", Name: "Kanto", Generation: "1º Geração"},
		{ID: "johto", Name: "Johto", Generation: "2º Geração"},
		{ID: "hoenn", Name: "Hoenn", Generation: "3º Geração"},
		{ID: "sinnoh", Name: "Sinnoh", Generation: "4º Geração"},
		{ID: "unova", Name: "Unova", Generation: "5º Geração"},
		{ID: "kalos", Name: "Kalos", Generation: "6º Geração"},
		{ID: "alola", Name: "Alola", Generation: "7º Geração"},
		{ID: "galar", Name: "Galar", Generation: "8º Geração"},
	}, nil
}

func convertMockTypes(types []string) []domain.Type {
	result := make([]domain.Type, len(types))
	for i, item := range types {
		result[i] = domain.Type{Name: item, Color: mockTypeColor(item)}
	}
	return result
}

func mockTypeColor(name string) string {
	switch name {
	case "Grass", "Grama":
		return "#7AC74C"
	case "Poison", "Venenoso":
		return "#A33EA1"
	case "Fire", "Fogo":
		return "#EE8130"
	case "Electric", "Elétrico":
		return "#F7D02C"
	case "Water", "Água":
		return "#6390F0"
	default:
		return "#A9AC86"
	}
}

type MockFavoriteRepository struct {
	mu        sync.RWMutex
	favorites map[string]map[string]bool
}

func NewMockFavoriteRepository() *MockFavoriteRepository {
	return &MockFavoriteRepository{
		favorites: make(map[string]map[string]bool),
	}
}

func (m *MockFavoriteRepository) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.favorites[userID]; !exists {
		m.favorites[userID] = make(map[string]bool)
	}

	if m.favorites[userID][pokemonID] {
		return domain.ErrFavoriteAlreadyExists
	}

	m.favorites[userID][pokemonID] = true
	return nil
}

func (m *MockFavoriteRepository) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.favorites[userID]; !exists {
		return domain.ErrFavoriteNotFound
	}

	if !m.favorites[userID][pokemonID] {
		return domain.ErrFavoriteNotFound
	}

	delete(m.favorites[userID], pokemonID)
	return nil
}

func (m *MockFavoriteRepository) IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.favorites[userID]; !exists {
		return false, nil
	}

	return m.favorites[userID][pokemonID], nil
}

func (m *MockFavoriteRepository) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.favorites[userID]; !exists {
		return []string{}, nil
	}

	var favorites []string
	for pokemonID := range m.favorites[userID] {
		favorites = append(favorites, pokemonID)
	}

	return favorites, nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

var _ ports.PokemonRepository = (*MockPokemonRepository)(nil)
var _ ports.FavoriteRepository = (*MockFavoriteRepository)(nil)
