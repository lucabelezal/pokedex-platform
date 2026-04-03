package httphandler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http/dto"
	inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
)

// Handler agrupa os handlers HTTP da aplicação.
type Handler struct {
	pokemonUseCase  inbound.PokemonUseCase
	favoriteUseCase inbound.FavoriteUseCase
	authUseCase     inbound.AuthUseCase
	responseBuilder *ResponseBuilder
}

// NewHandler cria um novo Handler com os use cases fornecidos.
func NewHandler(
	pokemonUseCase inbound.PokemonUseCase,
	favoriteUseCase inbound.FavoriteUseCase,
	authUseCase inbound.AuthUseCase,
) *Handler {
	return &Handler{
		pokemonUseCase:  pokemonUseCase,
		favoriteUseCase: favoriteUseCase,
		authUseCase:     authUseCase,
		responseBuilder: NewResponseBuilder(),
	}
}

// RegisterRoutes registra todas as rotas HTTP no mux fornecido.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/health", h.Health)
	mux.HandleFunc("POST /api/v1/auth/signup", h.Signup)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", h.Logout)
	mux.HandleFunc("GET /api/v1/me", h.GetMe)
	mux.HandleFunc("GET /api/v1/me/favorites", h.GetUserFavorites)
	mux.HandleFunc("GET /api/v1/pokemons", h.ListPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/search", h.SearchPokemons)
	mux.HandleFunc("GET /api/v1/pokemons/{id}/details", h.GetPokemonDetails)
	mux.HandleFunc("GET /api/v1/home", h.GetHome)
	mux.HandleFunc("GET /api/v1/regions", h.GetRegions)
	mux.HandleFunc("POST /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.AddFavorite))
	mux.HandleFunc("DELETE /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.RemoveFavorite))
}

// Health retorna o status de saúde do serviço.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildHealthResponse())
}

// GetMe retorna o perfil do usuário autenticado (ou não autenticado).
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	userID := getUserIDFromContext(ctx)
	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildProfileResponse(userID != "", getUserEmailFromContext(ctx)))
}

// RequireAuth envolve um handler para exigir autenticação.
func (h *Handler) RequireAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if getUserIDFromContext(r.Context()) == "" {
			RespondError(w, http.StatusUnauthorized, "autenticacao obrigatoria", "UNAUTHORIZED")
			return
		}
		handler(w, r)
	}
}

func (h *Handler) enrichFavoriteFlags(ctx context.Context, userID string, response *dto.RichPokemonListResponse) {
	if userID == "" || response == nil {
		return
	}
	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
	if err != nil {
		return
	}
	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, id := range favorites {
		favoriteSet[normalizePokemonID(id)] = struct{}{}
	}
	for i := range response.Content {
		_, isFavorite := favoriteSet[normalizePokemonID(response.Content[i].Number)]
		response.Content[i].IsFavorite = isFavorite
	}
}

func (h *Handler) buildFavoriteSet(ctx context.Context, userID string) map[string]struct{} {
	if userID == "" {
		return nil
	}
	favorites, err := h.favoriteUseCase.GetUserFavorites(ctx, userID)
	if err != nil {
		return nil
	}
	favoriteSet := make(map[string]struct{}, len(favorites))
	for _, id := range favorites {
		favoriteSet[normalizePokemonID(id)] = struct{}{}
	}
	return favoriteSet
}

func normalizePokemonID(value string) string {
	normalized := strings.TrimLeft(strings.TrimSpace(value), "0")
	if normalized == "" {
		return "0"
	}
	return normalized
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

func getOptionalQueryParamInt(r *http.Request, key string, defaultVal int) (int, bool) {
	values, exists := r.URL.Query()[key]
	if !exists || len(values) == 0 {
		return defaultVal, false
	}
	trimmed := strings.TrimSpace(values[0])
	if trimmed == "" {
		return defaultVal, false
	}
	intVal, err := strconv.Atoi(trimmed)
	if err != nil {
		return defaultVal, true
	}
	return intVal, true
}
