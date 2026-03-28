package mocks

import "pokedex-platform/core/bff/mobile-bff/internal/adapters/repository"

type MockPokemonRepository = repository.MockPokemonRepository
type MockFavoriteRepository = repository.MockFavoriteRepository

func NewMockPokemonRepository() *MockPokemonRepository {
	return repository.NewMockPokemonRepository()
}

func NewMockFavoriteRepository() *MockFavoriteRepository {
	return repository.NewMockFavoriteRepository()
}
