package unit

import (
	"context"
	"testing"

	"pokedex-platform/bff/mobile-bff/internal/domain"
	"pokedex-platform/bff/mobile-bff/internal/service"
	"pokedex-platform/bff/mobile-bff/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestPokemonServiceListPokemons(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewPokemonService(pokemonRepo, favoriteRepo)

	ctx := context.Background()
	page, err := svc.ListPokemons(ctx, 0, 10, "")

	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Greater(t, page.TotalElements, int64(0))
}

func TestPokemonServiceGetPokemonDetails(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewPokemonService(pokemonRepo, favoriteRepo)

	tests := []struct {
		name      string
		pokemonID string
		wantErr   bool
		errType   error
	}{
		{
			name:      "existing pokemon",
			pokemonID: "1",
			wantErr:   false,
		},
		{
			name:      "non-existing pokemon",
			pokemonID: "999",
			wantErr:   true,
			errType:   domain.ErrPokemonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			detail, err := svc.GetPokemonDetails(ctx, tt.pokemonID, "")

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, detail)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, detail)
				assert.Equal(t, "Bulbasaur", detail.Name)
			}
		})
	}
}

func TestPokemonServiceSearchPokemons(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewPokemonService(pokemonRepo, favoriteRepo)

	tests := []struct {
		name      string
		query     string
		wantFound bool
	}{
		{
			name:      "search for Pikachu",
			query:     "Pikachu",
			wantFound: true,
		},
		{
			name:      "search for number",
			query:     "025",
			wantFound: true,
		},
		{
			name:      "search for non-existing",
			query:     "Xyz123",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			page, err := svc.SearchPokemons(ctx, tt.query, 0, 10, "")

			assert.NoError(t, err)
			assert.NotNil(t, page)

			if tt.wantFound {
				assert.Greater(t, page.TotalElements, int64(0))
			} else {
				assert.Equal(t, int64(0), page.TotalElements)
			}
		})
	}
}

func TestFavoriteServiceAddFavorite(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	tests := []struct {
		name      string
		userID    string
		pokemonID string
		wantErr   bool
	}{
		{
			name:      "add valid favorite",
			userID:    "user123",
			pokemonID: "1",
			wantErr:   false,
		},
		{
			name:      "add non-existing pokemon",
			userID:    "user123",
			pokemonID: "999",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := svc.AddFavorite(ctx, tt.userID, tt.pokemonID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				isFav, err := favoriteRepo.IsFavorite(ctx, tt.userID, tt.pokemonID)
				assert.NoError(t, err)
				assert.True(t, isFav)
			}
		})
	}
}

func TestFavoriteServiceRemoveFavorite(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	ctx := context.Background()
	userID := "user123"
	pokemonID := "25"

	err := svc.AddFavorite(ctx, userID, pokemonID)
	assert.NoError(t, err)

	err = svc.RemoveFavorite(ctx, userID, pokemonID)
	assert.NoError(t, err)

	isFav, err := favoriteRepo.IsFavorite(ctx, userID, pokemonID)
	assert.NoError(t, err)
	assert.False(t, isFav)
}

func TestFavoriteServiceGetUserFavorites(t *testing.T) {
	pokemonRepo := mocks.NewMockPokemonRepository()
	favoriteRepo := mocks.NewMockFavoriteRepository()
	svc := service.NewFavoriteService(favoriteRepo, pokemonRepo)

	ctx := context.Background()
	userID := "user123"

	err := svc.AddFavorite(ctx, userID, "1")
	assert.NoError(t, err)
	err = svc.AddFavorite(ctx, userID, "25")
	assert.NoError(t, err)

	favs, err := svc.GetUserFavorites(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, favs, 2)
}
