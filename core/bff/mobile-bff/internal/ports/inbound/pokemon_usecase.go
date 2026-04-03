package inbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// PokemonUseCase define casos de uso para operacoes de Pokemon.
type PokemonUseCase interface {
	ListPokemons(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error)
	GetPokemonDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonDetail, error)
	GetPokemonScreenDetails(ctx context.Context, pokemonID, userID string) (*domain.PokemonScreenDetail, error)
	SearchPokemons(ctx context.Context, query string, page, pageSize int, userID string) (*domain.PokemonPage, error)
	FilterByType(ctx context.Context, typeFilter string, page, pageSize int, userID string) (*domain.PokemonPage, error)
	GetHomeData(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error)
	ListTypes(ctx context.Context) ([]domain.Type, error)
	ListRegions(ctx context.Context) ([]domain.Region, error)
}
