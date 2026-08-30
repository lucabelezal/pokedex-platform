package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	repository "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
)

// newCatalogWriteServer sobe um httptest.Server que simula o catalog-service
// para as rotas de escrita de favoritos, delegando ao repository real do BFF.
func newCatalogWriteServer(t *testing.T, db *repository.Database) *httptest.Server {
	t.Helper()

	repo := repository.NewPostgresFavoriteRepository(db.Pool)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		pokemonID := strings.TrimPrefix(r.URL.Path, "/v1/pokemons/")
		pokemonID = strings.TrimSuffix(pokemonID, "/favorite")

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		switch r.Method {
		case http.MethodPost:
			if err := repo.AddFavorite(ctx, userID, pokemonID); err != nil {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			if err := repo.RemoveFavorite(ctx, userID, pokemonID); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}
