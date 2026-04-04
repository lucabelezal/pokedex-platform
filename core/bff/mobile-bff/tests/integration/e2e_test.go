package integration

import (
	"context"
	"testing"
	"time"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http"
	repository "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
	"pokedex-platform/core/bff/mobile-bff/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EListaPokemonComBuilder(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)

	pokemonRepo := repository.NewPostgresPokemonRepository(db.Pool)
	favoriteRepo := repository.NewPostgresFavoriteRepository(db.Pool)
	pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)
	builder := httpadapter.NewResponseBuilder()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := pokemonService.ListPokemons(ctx, 0, 20, "")
	require.NoError(t, err)
	require.NotNil(t, page)

	resp := builder.BuildRichPokemonListResponse(page)
	require.NotNil(t, resp)
	assert.GreaterOrEqual(t, len(resp.Content), 1)
	assert.Equal(t, page.TotalElements, resp.TotalElements)
}

func TestE2EFavoritosNoServico(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)

	pokemonRepo := repository.NewPostgresPokemonRepository(db.Pool)
	favoriteRepo := repository.NewPostgresFavoriteRepository(db.Pool)
	pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := favoriteService.AddFavorite(ctx, "user-teste", "25")
	require.NoError(t, err)

	detail, err := pokemonService.GetPokemonDetails(ctx, "25", "user-teste")
	require.NoError(t, err)
	assert.True(t, detail.IsFavorite)
}
