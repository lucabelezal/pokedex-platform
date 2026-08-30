package httpclient_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpclient "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteCatalogClient_GetFavoriteDetails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/pokemons/favorites", r.URL.Path)
		assert.Equal(t, "1,25", r.URL.Query().Get("ids"))
		_ = json.NewEncoder(w).Encode([]domain.Pokemon{
			{ID: "1", Name: "Bulbasaur", Number: "001"},
			{ID: "25", Name: "Pikachu", Number: "025"},
		})
	}))
	defer srv.Close()

	client := httpclient.NewFavoriteCatalogClient(srv.URL)
	result, err := client.GetFavoriteDetails(context.Background(), []string{"1", "25"})
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Bulbasaur", result[0].Name)
}

func TestFavoriteCatalogClient_GetFavoriteDetails_Vazio(t *testing.T) {
	client := httpclient.NewFavoriteCatalogClient("http://localhost:1")
	result, err := client.GetFavoriteDetails(context.Background(), []string{})
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestFavoriteCatalogClient_GetFavoriteDetails_Erro4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := httpclient.NewFavoriteCatalogClient(srv.URL)
	_, err := client.GetFavoriteDetails(context.Background(), []string{"1"})
	require.Error(t, err)
}

func TestFavoriteCatalogClient_AddFavorite(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr error
	}{
		{"sucesso 200", http.StatusOK, nil},
		{"sucesso 201", http.StatusCreated, nil},
		{"conflito", http.StatusConflict, domain.ErrFavoriteAlreadyExists},
		{"erro 500", http.StatusInternalServerError, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/v1/pokemons/25/favorite", r.URL.Path)
				assert.Equal(t, "user-1", r.Header.Get("X-User-ID"))
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			client := httpclient.NewFavoriteCatalogClient(srv.URL)
			err := client.AddFavorite(context.Background(), "user-1", "25")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.status >= 400 {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFavoriteCatalogClient_RemoveFavorite(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr error
	}{
		{"sucesso 200", http.StatusOK, nil},
		{"nao encontrado", http.StatusNotFound, domain.ErrFavoriteNotFound},
		{"erro 500", http.StatusInternalServerError, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, "/v1/pokemons/25/favorite", r.URL.Path)
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			client := httpclient.NewFavoriteCatalogClient(srv.URL)
			err := client.RemoveFavorite(context.Background(), "user-1", "25")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.status >= 400 {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFavoriteCatalogClient_GetUserFavorites(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		status  int
		body    string
		wantLen int
		wantErr bool
	}{
		{"sem favoritos", "user-1", http.StatusOK, `[]`, 0, false},
		{"com favoritos", "user-1", http.StatusOK, `["1","25"]`, 2, false},
		{"erro 500", "user-1", http.StatusInternalServerError, ``, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/favorites", r.URL.Path)
				assert.Equal(t, "user-1", r.URL.Query().Get("user_id"))
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			client := httpclient.NewFavoriteCatalogClient(srv.URL)
			ids, err := client.GetUserFavorites(context.Background(), tt.userID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, ids, tt.wantLen)
		})
	}
}

func TestFavoriteCatalogClient_IsFavorite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`["1","25"]`))
	}))
	defer srv.Close()

	client := httpclient.NewFavoriteCatalogClient(srv.URL)

	isFav, err := client.IsFavorite(context.Background(), "user-1", "25")
	require.NoError(t, err)
	assert.True(t, isFav)

	isFav, err = client.IsFavorite(context.Background(), "user-1", "999")
	require.NoError(t, err)
	assert.False(t, isFav)
}
