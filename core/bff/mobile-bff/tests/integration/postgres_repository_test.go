package integration

import (
	"context"
	"os"
	"testing"
	"time"

	repository "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/postgres"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func urlBancoTeste() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5433/pokedex_test?sslmode=disable"
}

func setupTestDB(t *testing.T) *repository.Database {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repository.NewDatabase(ctx, urlBancoTeste())
	if err != nil {
		t.Skipf("Pulando testes de integração (banco indisponível): %v", err)
	}

	_, _ = db.Pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto`)
	_, err = db.Pool.Exec(ctx, `TRUNCATE TABLE favorites, users, pokemons RESTART IDENTITY CASCADE`)
	if err != nil {
		db.Close()
		t.Skipf("Pulando testes de integração (falha ao preparar banco): %v", err)
	}

	_, err = db.Pool.Exec(ctx, `INSERT INTO users (id, email) VALUES ('user-teste', 'teste@local') ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		db.Close()
		t.Skipf("Pulando testes de integração (falha ao inserir usuário): %v", err)
	}

	return db
}

func seedPokemonBasico(t *testing.T, db *repository.Database) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO pokemons (id, name, number, types, height, weight, description, image_url, element_color, element_type, created_at, updated_at)
		VALUES
		('1', 'Bulbasaur', '001', ARRAY['Grass','Poison'], 0.7, 6.9, 'Pokemon Semente', 'https://img/1.png', 'green', 'Grass', NOW(), NOW()),
		('25', 'Pikachu', '025', ARRAY['Electric'], 0.4, 6.0, 'Pokemon Rato', 'https://img/25.png', 'yellow', 'Electric', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`
	_, err := db.Pool.Exec(ctx, query)
	require.NoError(t, err)
}

func TestPostgresPokemonRepositoryGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)

	repo := repository.NewPostgresPokemonRepository(db.Pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := repo.GetByID(ctx, "1")
	require.NoError(t, err)
	assert.Equal(t, "Bulbasaur", p.Name)
	assert.Equal(t, "001", p.Number)
}

func TestPostgresPokemonRepositoryGetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewPostgresPokemonRepository(db.Pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := repo.GetByID(ctx, "999")
	assert.Error(t, err)
	assert.Nil(t, p)
	assert.Equal(t, domain.ErrPokemonNotFound, err)
}

func TestPostgresFavoriteRepositoryFluxoBasico(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedPokemonBasico(t, db)

	repo := repository.NewPostgresFavoriteRepository(db.Pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.AddFavorite(ctx, "user-teste", "1")
	require.NoError(t, err)

	isFav, err := repo.IsFavorite(ctx, "user-teste", "1")
	require.NoError(t, err)
	assert.True(t, isFav)

	err = repo.RemoveFavorite(ctx, "user-teste", "1")
	require.NoError(t, err)

	isFav, err = repo.IsFavorite(ctx, "user-teste", "1")
	require.NoError(t, err)
	assert.False(t, isFav)
}
