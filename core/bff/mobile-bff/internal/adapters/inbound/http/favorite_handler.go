package httphandler

import (
	"context"
	"net/http"
	"sort"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// AddFavorite adiciona um Pokémon aos favoritos do usuário autenticado.
func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	err := h.favoriteUseCase.AddFavorite(ctx, userID, pokemonID)
	if err == domain.ErrFavoriteAlreadyExists {
		RespondError(w, http.StatusConflict, "pokemon ja esta nos favoritos", "ALREADY_EXISTS")
		return
	}
	if err == domain.ErrPokemonNotFound {
		RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao adicionar favorito", "INTERNAL_ERROR")
		return
	}

	RespondJSON(w, http.StatusOK, dto.FavoriteResponse{
		Message:    "Pokemon adicionado aos favoritos",
		PokemonID:  pokemonID,
		IsFavorite: true,
	})
}

// RemoveFavorite remove um Pokémon dos favoritos do usuário autenticado.
func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	err := h.favoriteUseCase.RemoveFavorite(ctx, userID, pokemonID)
	if err == domain.ErrFavoriteNotFound {
		RespondError(w, http.StatusNotFound, "favorito nao encontrado", "NOT_FOUND")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao remover favorito", "INTERNAL_ERROR")
		return
	}

	RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "Pokemon removido dos favoritos"})
}

// GetUserFavorites retorna a lista de Pokémons favoritos do usuário autenticado.
func (h *Handler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := getUserIDFromContext(ctx)
	if userID == "" {
		RespondJSON(w, http.StatusOK, h.responseBuilder.BuildFavoritesResponse(nil, nil, false))
		return
	}

	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar favoritos", "INTERNAL_ERROR")
		return
	}

	items := make([]domain.Pokemon, 0, len(favorites))
	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, favoriteID := range favorites {
		favoriteSet[normalizePokemonID(favoriteID)] = struct{}{}
		pokemon, err := h.pokemonUseCase.GetPokemonScreenDetails(ctx, favoriteID, userID)
		if err != nil {
			continue
		}
		items = append(items, domain.Pokemon{
			ID:           pokemon.ID,
			Name:         pokemon.Name,
			Number:       pokemon.Number,
			Types:        mapScreenTypesToNames(pokemon.Types),
			ImageURL:     pokemon.ImageURL,
			ElementColor: pokemon.ElementColor,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Number < items[j].Number
	})

	page := &domain.PokemonPage{Content: items}
	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildFavoritesResponse(page, favoriteSet, true))
}

func mapScreenTypesToNames(types []domain.Type) []string {
	result := make([]string, len(types))
	for i, item := range types {
		result[i] = item.Name
	}
	return result
}
