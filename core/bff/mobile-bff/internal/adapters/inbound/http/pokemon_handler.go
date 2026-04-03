package httphandler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// ListPokemons lista Pokémons com suporte a filtro por tipo.
func (h *Handler) ListPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)
	typeFilter := r.URL.Query().Get("type")

	var (
		pokemonPage *domain.PokemonPage
		err         error
	)

	if typeFilter != "" {
		pokemonPage, err = h.pokemonUseCase.FilterByType(ctx, typeFilter, page, pageSize, userID)
	} else {
		pokemonPage, err = h.pokemonUseCase.ListPokemons(ctx, page, pageSize, userID)
	}

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar pokemons", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildRichPokemonListResponse(pokemonPage)
	h.enrichFavoriteFlags(ctx, userID, response)
	RespondJSON(w, http.StatusOK, response)
}

// SearchPokemons busca Pokémons pelo nome ou número.
func (h *Handler) SearchPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query().Get("q")
	if query == "" {
		RespondError(w, http.StatusBadRequest, "termo de busca obrigatorio", "INVALID_REQUEST")
		return
	}

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)

	pokemonPage, err := h.pokemonUseCase.SearchPokemons(ctx, query, page, pageSize, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao buscar pokemons", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildRichPokemonListResponse(pokemonPage)
	h.enrichFavoriteFlags(ctx, userID, response)
	RespondJSON(w, http.StatusOK, response)
}

// GetPokemonDetails retorna os detalhes de tela de um Pokémon pelo ID.
func (h *Handler) GetPokemonDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	detail, err := h.pokemonUseCase.GetPokemonScreenDetails(ctx, pokemonID, userID)
	if err != nil {
		if err == domain.ErrPokemonNotFound {
			RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
			return
		}
		RespondError(w, http.StatusInternalServerError, "falha ao obter detalhes do pokemon", "INTERNAL_ERROR")
		return
	}

	isFavorite := false
	if userID != "" {
		favoriteSet := h.buildFavoriteSet(ctx, userID)
		_, isFavorite = favoriteSet[normalizePokemonID(detail.Number)]
	}

	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildPokemonDetailScreenResponse(detail, isFavorite))
}

func hasPokemonType(types []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, item := range types {
		if strings.EqualFold(strings.TrimSpace(item), target) {
			return true
		}
	}
	return false
}
