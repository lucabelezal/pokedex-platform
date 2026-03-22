package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"pokedex-platform/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/bff/mobile-bff/internal/domain"
	"pokedex-platform/bff/mobile-bff/internal/ports"
)

type Handler struct {
	pokemonUseCase  ports.PokemonUseCase
	favoriteUseCase ports.FavoriteUseCase
	responseBuilder *ResponseBuilder
}

func NewHandler(
	pokemonUseCase ports.PokemonUseCase,
	favoriteUseCase ports.FavoriteUseCase,
) *Handler {
	return &Handler{
		pokemonUseCase:  pokemonUseCase,
		favoriteUseCase: favoriteUseCase,
		responseBuilder: NewResponseBuilder(),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/pokemons", h.ListPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/search", h.SearchPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/{id}/details", h.GetPokemonDetails)
	mux.HandleFunc("GET /api/v1/home", h.GetHome)
	mux.HandleFunc("POST /api/v1/pokemons/{id}/favorite", h.withAuth(h.AddFavorite))
	mux.HandleFunc("DELETE /api/v1/pokemons/{id}/favorite", h.withAuth(h.RemoveFavorite))
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	health := h.responseBuilder.BuildHealthResponse()
	RespondJSON(w, http.StatusOK, health)
}

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
	RespondJSON(w, http.StatusOK, response)
}

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
	RespondJSON(w, http.StatusOK, response)
}
func (h *Handler) GetPokemonDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		RespondError(w, http.StatusBadRequest, "id do pokemon obrigatorio", "INVALID_REQUEST")
		return
	}

	userID := getUserIDFromContext(ctx)

	detail, err := h.pokemonUseCase.GetPokemonDetails(ctx, pokemonID, userID)
	if err != nil {
		if err == domain.ErrPokemonNotFound {
			RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
			return
		}
		RespondError(w, http.StatusInternalServerError, "falha ao obter detalhes do pokemon", "INTERNAL_ERROR")
		return
	}

	detailDTO := h.responseBuilder.BuildPokemonDetailDTO(detail)
	RespondJSON(w, http.StatusOK, detailDTO)
}

func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := getQueryParamInt(r, "page", 0)
	pageSize := getQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)

	pokemonPage, err := h.pokemonUseCase.GetHomeData(ctx, page, pageSize, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter dados da home", "INTERNAL_ERROR")
		return
	}

	response := h.responseBuilder.BuildHomePageResponse(pokemonPage)
	RespondJSON(w, http.StatusOK, response)
}

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
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao adicionar favorito", "INTERNAL_ERROR")
		return
	}

	response := dto.FavoriteResponse{
		Message:    "Pokemon adicionado aos favoritos",
		PokemonID:  pokemonID,
		IsFavorite: true,
	}
	RespondJSON(w, http.StatusOK, response)
}

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

	response := dto.MessageResponse{
		Message: "Pokemon removido dos favoritos",
	}
	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) withAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())
		if userID == "" {
			RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
			return
		}
		handler(w, r)
	}
}

func getQueryParamInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}
