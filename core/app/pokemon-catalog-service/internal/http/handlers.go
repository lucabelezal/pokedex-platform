package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/domain"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/repository"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type readyResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Pinger abstrai a verificação de saúde do banco de dados.
type Pinger interface {
	Ping(ctx context.Context) error
}

func ReadyHandler(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			respondJSON(w, http.StatusServiceUnavailable, readyResponse{Status: "degraded", Service: "pokemon-catalog-service"})
			return
		}
		respondJSON(w, http.StatusOK, readyResponse{Status: "ready", Service: "pokemon-catalog-service"})
	}
}

type pingResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	pokemonRepo repository.PokemonRepository
}

func NewMux(pokemonRepo repository.PokemonRepository) *http.ServeMux {
	h := &Handler{pokemonRepo: pokemonRepo}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /v1/pokemon/ping", pingHandler)
	mux.HandleFunc("GET /v1/pokemons", h.listPokemons)
	mux.HandleFunc("GET /v1/pokemons/search", h.searchPokemons)
	mux.HandleFunc("GET /v1/pokemons/type/{type}", h.filterByType)
	mux.HandleFunc("GET /v1/types", h.listTypes)
	mux.HandleFunc("GET /v1/regions", h.listRegions)
	mux.HandleFunc("GET /v1/pokemon-details/{id}", h.getPokemonDetailByID)
	mux.HandleFunc("GET /v1/pokemons/{id}", h.getPokemonByID)
	mux.HandleFunc("GET /v1/pokemons/favorites", h.getFavoritesBatch)
	mux.HandleFunc("POST /v1/pokemons/{id}/favorite", h.addFavorite)
	mux.HandleFunc("DELETE /v1/pokemons/{id}/favorite", h.removeFavorite)
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "pokemon-catalog-service"})
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, pingResponse{Message: "service is alive"})
}

func (h *Handler) listPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := queryInt(r, "page", 0)
	pageSize := queryInt(r, "size", 20)

	data, err := h.pokemonRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao listar pokemons"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) searchPokemons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "termo de busca obrigatorio"})
		return
	}

	page := queryInt(r, "page", 0)
	pageSize := queryInt(r, "size", 20)

	data, err := h.pokemonRepo.Search(ctx, q, page, pageSize)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao buscar pokemons"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) filterByType(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	typeFilter := strings.TrimSpace(r.PathValue("type"))
	if typeFilter == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "tipo obrigatorio"})
		return
	}

	page := queryInt(r, "page", 0)
	pageSize := queryInt(r, "size", 20)

	data, err := h.pokemonRepo.GetByType(ctx, typeFilter, page, pageSize)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao filtrar pokemons"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) getPokemonByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "id obrigatorio"})
		return
	}

	data, err := h.pokemonRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrPokemonNotFound {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "pokemon nao encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao buscar pokemon"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) listTypes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.pokemonRepo.ListTypes(ctx)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao listar tipos"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) listRegions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.pokemonRepo.ListRegions(ctx)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao listar regioes"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) getPokemonDetailByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "id obrigatorio"})
		return
	}

	data, err := h.pokemonRepo.GetDetailByID(ctx, id)
	if err != nil {
		if err == repository.ErrPokemonNotFound {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "pokemon nao encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao buscar detalhes do pokemon"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) getFavoritesBatch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idsParam := strings.TrimSpace(r.URL.Query().Get("ids"))
	if idsParam == "" {
		respondJSON(w, http.StatusOK, []domain.Pokemon{})
		return
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) > 100 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "maximo de 100 IDs por requisicao"})
		return
	}

	trimmedIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		tid := strings.TrimSpace(id)
		if tid != "" {
			trimmedIDs = append(trimmedIDs, tid)
		}
	}

	data, err := h.pokemonRepo.GetByIDs(ctx, trimmedIDs)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao buscar pokemons"})
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) addFavorite(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "autenticacao obrigatoria"})
		return
	}

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "id do pokemon obrigatorio"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.pokemonRepo.AddFavorite(ctx, userID, pokemonID); err != nil {
		if err == repository.ErrFavoriteAlreadyExists {
			respondJSON(w, http.StatusConflict, map[string]string{"error": "favorito ja existe"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao adicionar favorito"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "favorito adicionado"})
}

func (h *Handler) removeFavorite(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "autenticacao obrigatoria"})
		return
	}

	pokemonID := r.PathValue("id")
	if pokemonID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "id do pokemon obrigatorio"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.pokemonRepo.RemoveFavorite(ctx, userID, pokemonID); err != nil {
		if err == repository.ErrFavoriteNotFound {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "favorito nao encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao remover favorito"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "favorito removido"})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func queryInt(r *http.Request, key string, defaultValue int) int {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return v
}
