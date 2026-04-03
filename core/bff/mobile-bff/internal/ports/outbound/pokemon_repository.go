package outbound

import (
	"context"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// PokemonRepository define o contrato para acessar dados de Pokemon.
type PokemonRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
	GetDetailByID(ctx context.Context, id string) (*domain.PokemonScreenDetail, error)
	GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	ListTypes(ctx context.Context) ([]domain.Type, error)
	ListRegions(ctx context.Context) ([]domain.Region, error)
	GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error)
}
