package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "pokedex-platform/bff/mobile-bff/internal/adapters/http"
	"pokedex-platform/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/bff/mobile-bff/internal/service"
	"pokedex-platform/bff/mobile-bff/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "mobile-bff", response.Service)
}

func TestListPokemonsHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/pokemons?page=0&size=10", nil)
	w := httptest.NewRecorder()

	handler.ListPokemons(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RichPokemonListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Search.Placeholder)
	assert.Greater(t, len(response.Content), 0)
}

func TestListPokemonsHandlerWithAuthenticatedFavorites(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	err := favoriteRepo.AddFavorite(context.Background(), "user123", "25")
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/pokemons?page=0&size=10", nil)
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	handler.ListPokemons(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RichPokemonListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var pikachuFound bool
	for _, pokemon := range response.Content {
		if pokemon.Number == "025" {
			pikachuFound = true
			assert.True(t, pokemon.IsFavorite)
		}
	}

	assert.True(t, pikachuFound)
}

func TestSearchPokemonsHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/pokemons/search?q=Pikachu", nil)
	w := httptest.NewRecorder()

	handler.SearchPokemons(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RichPokemonListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Content, 1)
	assert.Equal(t, "Pikachu", response.Content[0].Name)
}

func TestListPokemonsWithTypeFilterHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/pokemons?type=Electric&page=0&size=10", nil)
	w := httptest.NewRecorder()

	handler.ListPokemons(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RichPokemonListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Content)
	assert.Equal(t, "Pikachu", response.Content[0].Name)
}

func TestGetPokemonDetailsHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/pokemons/25/details", nil)
	req.SetPathValue("id", "25")
	w := httptest.NewRecorder()

	handler.GetPokemonDetails(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PokemonDetailDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Pikachu", response.Name)
	assert.Equal(t, "025", response.Number)
}

func TestGetHomeHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/home", nil)
	w := httptest.NewRecorder()

	handler.GetHome(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.SearchPlaceholder)
	assert.NotEmpty(t, response.Data)
	assert.Greater(t, len(response.Data.Content), 0)
}

func TestGetMeHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user123"))
	req = req.WithContext(httpadapter.SetUserEmail(req.Context(), "user123@example.com"))
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Authenticated bool   `json:"authenticated"`
		UserID        string `json:"user_id"`
		Email         string `json:"email"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Authenticated)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "user123@example.com", response.Email)
}

func TestGetMeWithoutAuth(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddFavoriteWithoutAuth(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("POST", "/api/v1/pokemons/25/favorite", nil)
	req.SetPathValue("id", "25")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRemoveFavoriteNotFound(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil, favoriteRepo)

	req := httptest.NewRequest("DELETE", "/api/v1/pokemons/25/favorite", nil)
	req.SetPathValue("id", "25")
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	handler.RemoveFavorite(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", response.Error)
}
