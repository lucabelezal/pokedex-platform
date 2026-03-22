package http

import (
	"context"
	"pokedex-platform/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/bff/mobile-bff/internal/domain"
	"pokedex-platform/bff/mobile-bff/internal/ports"
)

// EnrichedResponseBuilder builds enriched responses with favorites and additional info
type EnrichedResponseBuilder struct {
	responseBuilder *ResponseBuilder
	favoriteRepo    ports.FavoriteRepository
	pokemonRepo     ports.PokemonRepository
}

// NewEnrichedResponseBuilder creates a new enriched response builder
func NewEnrichedResponseBuilder(
	favoriteRepo ports.FavoriteRepository,
	pokemonRepo ports.PokemonRepository,
) *EnrichedResponseBuilder {
	return &EnrichedResponseBuilder{
		responseBuilder: NewResponseBuilder(),
		favoriteRepo:    favoriteRepo,
		pokemonRepo:     pokemonRepo,
	}
}

// BuildEnrichedPokemonDetail builds a detail response with favorites info
func (b *EnrichedResponseBuilder) BuildEnrichedPokemonDetail(
	ctx context.Context,
	detail *domain.PokemonDetail,
	userID string,
) *dto.PokemonDetailDTO {
	response := b.responseBuilder.BuildPokemonDetailDTO(detail)

	// Check if favorite
	if userID != "" && b.favoriteRepo != nil {
		isFav, err := b.favoriteRepo.IsFavorite(ctx, userID, detail.Number)
		if err == nil {
			response.IsFavorite = isFav
		}
	}

	return response
}

// BuildEnrichedListResponse builds a list response with favorite counts
func (b *EnrichedResponseBuilder) BuildEnrichedListResponse(
	ctx context.Context,
	page *domain.PokemonPage,
	userID string,
) *dto.RichPokemonListResponse {
	response := b.responseBuilder.BuildRichPokemonListResponse(page)

	// Enrich each pokemon with favorite status if user is logged in
	if userID != "" && b.favoriteRepo != nil {
		for i := range response.Content {
			isFav, err := b.favoriteRepo.IsFavorite(ctx, userID, response.Content[i].Number)
			if err == nil {
				response.Content[i].IsFavorite = isFav
			}
		}
	}

	return response
}

// BuildHomePageResponse builds home page with featured and trending Pokemon
func (b *EnrichedResponseBuilder) BuildHomePageResponse(
	ctx context.Context,
	page *domain.PokemonPage,
	userID string,
) *dto.HomeResponse {
	richResponse := b.BuildEnrichedListResponse(ctx, page, userID)

	homeResponse := &dto.HomeResponse{
		Status:             "success",
		Data:               richResponse,
		SearchPlaceholder:  "Search Pokemon by name or ID",
		RecommendedFilters: []string{"Fire", "Water", "Grass", "Electric", "Flying"},
		Message:            "Welcome to Pokedex",
	}

	return homeResponse
}

// BuildFavoritePokemonResponse builds response for favorite operations
func (b *EnrichedResponseBuilder) BuildFavoritePokemonResponse(
	ctx context.Context,
	pokemonID string,
	isFavorite bool,
) *dto.FavoriteResponse {
	action := "added"
	if !isFavorite {
		action = "removed"
	}

	return &dto.FavoriteResponse{
		Message:    "Pokemon " + action + " to favorites",
		PokemonID:  pokemonID,
		IsFavorite: isFavorite,
	}
}
