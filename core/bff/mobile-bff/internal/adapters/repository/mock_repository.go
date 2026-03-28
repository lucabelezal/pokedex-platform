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
