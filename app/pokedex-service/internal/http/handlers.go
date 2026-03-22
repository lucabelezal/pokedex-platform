package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/app/pokedex-service/internal/repository"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
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
	mux.HandleFunc("GET /v1/pokemons/{id}", h.getPokemonByID)
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "pokedex-service"})
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
