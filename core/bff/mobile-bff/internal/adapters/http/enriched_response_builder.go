package http

import (
	"context"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
)

// EnrichedResponseBuilder constrói respostas enriquecidas com favoritos e informações adicionais
type EnrichedResponseBuilder struct {
	responseBuilder *ResponseBuilder
	favoriteRepo    ports.FavoriteRepository
	pokemonRepo     ports.PokemonRepository
}

// NewEnrichedResponseBuilder cria um novo construtor de respostas enriquecidas
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

// BuildEnrichedPokemonDetail constrói uma resposta de detalhe com informação de favoritos
func (b *EnrichedResponseBuilder) BuildEnrichedPokemonDetail(
	ctx context.Context,
	detail *domain.PokemonDetail,
	userID string,
) *dto.PokemonDetailDTO {
	response := b.responseBuilder.BuildPokemonDetailDTO(detail)

	if userID != "" && b.favoriteRepo != nil {
		isFav, err := b.favoriteRepo.IsFavorite(ctx, userID, detail.Number)
		if err == nil {
			response.IsFavorite = isFav
		}
	}

	return response
}

// BuildEnrichedListResponse constrói uma resposta de lista com contagem de favoritos
func (b *EnrichedResponseBuilder) BuildEnrichedListResponse(
	ctx context.Context,
	page *domain.PokemonPage,
	userID string,
) *dto.RichPokemonListResponse {
	response := b.responseBuilder.BuildRichPokemonListResponse(page)

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

// BuildHomePageResponse constrói página home com Pokémons em destaque e em alta
func (b *EnrichedResponseBuilder) BuildHomePageResponse(
	ctx context.Context,
	page *domain.PokemonPage,
	userID string,
) *dto.HomeResponse {
	favoriteSet := map[string]struct{}{}
	if userID != "" && b.favoriteRepo != nil {
		favorites, err := b.favoriteRepo.GetUserFavorites(ctx, userID)
		if err == nil {
			for _, id := range favorites {
				favoriteSet[normalizePokemonID(id)] = struct{}{}
			}
		}
	}

	types, err := b.pokemonRepo.ListTypes(ctx)
	if err != nil {
		types = nil
	}

	return b.responseBuilder.BuildHomePageResponseWithTypes(page, types, favoriteSet)
}

// BuildFavoritePokemonResponse constrói resposta para operações de favoritos
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
