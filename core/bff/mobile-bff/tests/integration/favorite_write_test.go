package integration

import (
	"context"
	"testing"
	"time"

	httpclient "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"
	repository "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
	"pokedex-platform/core/bff/mobile-bff/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteEscritaViaCatalogProvider(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)
	seedTestUser(t, db)

	// Simula o catalog-service com um servidor HTTP que delega ao repository real.
	srv := newCatalogWriteServer(t, db)
	defer srv.Close()

	pokemonRepo := repository.NewPostgresPokemonRepository(db.Pool)
	favoriteRepo := repository.NewPostgresFavoriteRepository(db.Pool)
	catalogProvider := httpclient.NewFavoriteCatalogClient(srv.URL)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo, catalogProvider)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := favoriteService.AddFavorite(ctx, testUserID, "25")
	require.NoError(t, err)

	err = favoriteService.RemoveFavorite(ctx, testUserID, "25")
	require.NoError(t, err)

	ids, err := favoriteService.GetUserFavorites(ctx, testUserID)
	require.NoError(t, err)
	assert.NotContains(t, ids, "25")
}

func TestFavoriteEscritaConcorrente(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)
	seedTestUser(t, db)

	srv := newCatalogWriteServer(t, db)
	defer srv.Close()

	pokemonRepo := repository.NewPostgresPokemonRepository(db.Pool)
	favoriteRepo := repository.NewPostgresFavoriteRepository(db.Pool)
	catalogProvider := httpclient.NewFavoriteCatalogClient(srv.URL)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo, catalogProvider)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const n = 5
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			errs <- favoriteService.AddFavorite(ctx, testUserID, "25")
		}()
	}

	count := 0
	for i := 0; i < n; i++ {
		err := <-errs
		if err == nil {
			count++
		}
	}

	assert.GreaterOrEqual(t, count, 1)

	ids, err := favoriteService.GetUserFavorites(ctx, testUserID)
	require.NoError(t, err)
	countIDs := 0
	for _, id := range ids {
		if id == "25" {
			countIDs++
		}
	}
	assert.Equal(t, 1, countIDs, "apenas uma linha deve ser inserida para o mesmo favorito")
}
