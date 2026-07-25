package mocks

import (
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/memory"
)

// MockPokemonRepository é um alias de memory.PokemonRepository para compatibilidade com os testes existentes.
type MockPokemonRepository = memory.PokemonRepository

// MockFavoriteRepository é um alias de memory.FavoriteRepository para compatibilidade com os testes existentes.
type MockFavoriteRepository = memory.FavoriteRepository

// NewMockPokemonRepository cria um repositório de Pokémon em memória.
func NewMockPokemonRepository() *MockPokemonRepository {
	return memory.NewPokemonRepository()
}

// NewMockFavoriteRepository cria um repositório de favoritos em memória.
func NewMockFavoriteRepository() *MockFavoriteRepository {
	return memory.NewFavoriteRepository()
}
