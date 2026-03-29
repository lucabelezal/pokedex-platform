package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"pokedex-platform/core/bff/mobile-bff/internal/ports"
	"pokedex-platform/core/bff/mobile-bff/internal/service"
	"pokedex-platform/core/bff/mobile-bff/tests/mocks"

	"github.com/stretchr/testify/assert"
)

type stubAuthUseCase struct {
	session   *ports.AuthSession
	err       error
	logoutErr error
}

func (s *stubAuthUseCase) Signup(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	return s.session, s.err
}

func (s *stubAuthUseCase) Login(ctx context.Context, email, password string) (*ports.AuthSession, error) {
	return s.session, s.err
}

func (s *stubAuthUseCase) Refresh(ctx context.Context, token string) (*ports.AuthSession, error) {
	return s.session, s.err
}

func (s *stubAuthUseCase) Logout(ctx context.Context, token string) error {
	if s.logoutErr != nil {
		return s.logoutErr
	}
	return s.err
}

func TestHealthHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/pokemons/25/details", nil)
	req.SetPathValue("id", "25")
	w := httptest.NewRecorder()

	handler.GetPokemonDetails(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PokemonDetailScreenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Pikachu", response.Name)
	assert.Equal(t, "Nº025", response.Number)
}

func TestGetHomeHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/home", nil)
	w := httptest.NewRecorder()

	handler.GetHome(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Pokédex", response.Title)
	assert.Equal(t, "Procurar Pokémon...", response.Search.Placeholder)
	assert.Equal(t, "Tipos", response.Filters.Types.Title)
	assert.Equal(t, "Todos os tipos", response.Filters.Types.Selected)
	assert.NotEmpty(t, response.Filters.Types.Items)
	assert.Equal(t, "Ordenação", response.Filters.Ordering.Title)
	assert.Equal(t, "Menor número", response.Filters.Ordering.Selected)
	assert.Len(t, response.Filters.Ordering.Items, 4)
	assert.Greater(t, len(response.Pokemons), 0)
	assert.Equal(t, "Nº001", response.Pokemons[0].Number)
	assert.NotEmpty(t, response.Pokemons[0].Sprites.URL)
	assert.NotEmpty(t, response.Pokemons[0].Sprites.BackgroundColor)
}

func TestGetHomeHandlerWithFilters(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/home?q=char&type=Fire&order=A-Z&region=kanto", nil)
	w := httptest.NewRecorder()

	handler.GetHome(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "char", response.Search.Value)
	assert.Equal(t, "Fire", response.Filters.Types.Selected)
	assert.Equal(t, "A-Z", response.Filters.Ordering.Selected)
	assert.Equal(t, "kanto", response.Filters.Region.Selected)
	assert.Len(t, response.Pokemons, 1)
	if len(response.Pokemons) > 0 {
		assert.Equal(t, "Charmander", response.Pokemons[0].Name)
	}
}

func TestGetHomeHandlerWithDescendingNumberOrdering(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/home?order=Maior+n%C3%BAmero", nil)
	w := httptest.NewRecorder()

	handler.GetHome(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Pokemons), 3)
	assert.Equal(t, "Nº025", response.Pokemons[0].Number)
	assert.Equal(t, "Nº004", response.Pokemons[1].Number)
	assert.Equal(t, "Nº001", response.Pokemons[2].Number)
}

func TestGetRegionsHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/regions", nil)
	w := httptest.NewRecorder()

	handler.GetRegions(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RegionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Regiões", response.Title)
	assert.GreaterOrEqual(t, len(response.Regions), 8)
	assert.Equal(t, "kanto", response.Regions[0].ID)
	assert.Equal(t, "Kanto", response.Regions[0].Name)
	assert.Equal(t, "1º Geração", response.Regions[0].Generation)
}

func TestGetMeHandler(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user123"))
	req = req.WithContext(httpadapter.SetUserEmail(req.Context(), "user123@example.com"))
	w := httptest.NewRecorder()

	handler.GetMe(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Authenticated)
	assert.NotNil(t, response.User)
	assert.Equal(t, "user123@example.com", response.User.Email)
	assert.NotEmpty(t, response.Sections)
}

func TestGetMeWithoutAuth(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Authenticated)
	assert.NotNil(t, response.Header)
	assert.Equal(t, "Entre ou Cadastre-se", response.Header.Title)
}

func TestGetFavoritesWithoutAuthReturnsScreenState(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/me/favorites", nil)
	w := httptest.NewRecorder()

	handler.GetUserFavorites(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.FavoritesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthenticated", response.State)
	assert.NotNil(t, response.Message)
}

func TestGetFavoritesWithAuthAndEmptyListReturnsEmptyState(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("GET", "/api/v1/me/favorites", nil)
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user-empty"))
	w := httptest.NewRecorder()

	handler.GetUserFavorites(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.FavoritesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "empty", response.State)
	assert.NotNil(t, response.Message)
	assert.Empty(t, response.Pokemons)
}

func TestGetFavoritesWithAuthAndDataReturnsHasDataState(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	err := favoriteRepo.AddFavorite(context.Background(), "user-data", "25")
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/me/favorites", nil)
	req = req.WithContext(httpadapter.SetUserID(req.Context(), "user-data"))
	w := httptest.NewRecorder()

	handler.GetUserFavorites(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.FavoritesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "has_data", response.State)
	assert.Len(t, response.Pokemons, 1)
	assert.Equal(t, "Pikachu", response.Pokemons[0].Name)
	assert.True(t, response.Pokemons[0].IsFavorite)
}

func TestAddFavoriteWithoutAuth(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

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

func TestSignupReturnsCreatedWhenAuthSucceeds(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)
	authUseCase := &stubAuthUseCase{
		session: &ports.AuthSession{
			AccessToken: "token-123",
			TokenType:   "Bearer",
			ExpiresIn:   900,
			UserID:      "user-1",
			Email:       "ash@kanto.dev",
		},
	}

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, authUseCase)

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(`{"email":"ash@kanto.dev","password":"pikachu123"}`))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Header().Get("Set-Cookie"), "auth_token=token-123")
}

func TestLoginReturnsUnauthorizedForInvalidCredentials(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, &stubAuthUseCase{err: domain.ErrInvalidCredentials})

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"ash@kanto.dev","password":"wrong"}`))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_ERROR", response.Error)
}

func TestRefreshReturnsUnauthorizedForInvalidToken(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, &stubAuthUseCase{err: domain.ErrInvalidToken})

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	w := httptest.NewRecorder()

	handler.Refresh(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_TOKEN", response.Error)
}

func TestSignupReturnsServiceUnavailableWhenAuthUseCaseIsMissing(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	pokemonSvc := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	handler := httpadapter.NewHandler(pokemonSvc, favoriteSvc, nil)

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(`{"email":"ash@kanto.dev","password":"pikachu123"}`))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_UNAVAILABLE", response.Error)
}
