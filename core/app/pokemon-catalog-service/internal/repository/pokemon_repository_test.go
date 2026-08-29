package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	return mock
}

func TestPostgresPokemonRepository_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockSetup  func(pgxmock.PgxPoolIface)
		wantErr    error
		wantName   string
		wantNumber string
	}{
		{
			name: "encontrado por number",
			id:   "001",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				cols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
				m.ExpectQuery(regexp.QuoteMeta("SELECT")).
					WithArgs("001").
					WillReturnRows(pgxmock.NewRows(cols).
						AddRow("uuid-1", "Bulbasaur", "001", []string{"Grass", "Poison"}, 0.7, 6.9, "seed", "https://img/1.png", "#78C850", "Grass", time.Now(), time.Now()))
			},
			wantName:   "Bulbasaur",
			wantNumber: "001",
		},
		{
			name: "nao encontrado retorna ErrPokemonNotFound",
			id:   "999",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta("SELECT")).
					WithArgs("999").
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}))
			},
			wantErr: ErrPokemonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			mock := mustMock(t)
			defer func() { mock.Close() }()
			tt.mockSetup(mock)

			repo := NewPostgresPokemonRepositoryWithPool(mock)
			p, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				is.Nil(p)
				return
			}
			require.NoError(t, err)
			is.Equal(tt.wantName, p.Name)
			is.Equal(tt.wantNumber, p.Number)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresPokemonRepository_GetAll(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	cols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(20, 0).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("uuid-1", "Bulbasaur", "001", []string{"Grass"}, 0.7, 6.9, "seed", "https://img/1.png", "#78C850", "Grass", time.Now(), time.Now()).
			AddRow("uuid-4", "Charmander", "004", []string{"Fire"}, 0.6, 8.5, "fire", "https://img/4.png", "#F08030", "Fire", time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT")).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(2)))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	page, err := repo.GetAll(context.Background(), 0, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(2), page.TotalElements)
	assert.Len(t, page.Content, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_Search(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	cols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("%pika%", 20, 0).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("uuid-25", "Pikachu", "025", []string{"Electric"}, 0.4, 6.0, "mouse", "https://img/25.png", "#F8D030", "Electric", time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT")).
		WithArgs("%pika%").
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(1)))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	page, err := repo.Search(context.Background(), "pika", 0, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(1), page.TotalElements)
	assert.Len(t, page.Content, 1)
	assert.Equal(t, "Pikachu", page.Content[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_GetByType(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	cols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("Fire", 20, 0).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("uuid-4", "Charmander", "004", []string{"Fire"}, 0.6, 8.5, "fire", "https://img/4.png", "#F08030", "Fire", time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT")).
		WithArgs("Fire").
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(1)))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	page, err := repo.GetByType(context.Background(), "Fire", 0, 20)
	require.NoError(t, err)
	assert.Len(t, page.Content, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_GetByIDs(t *testing.T) {
	t.Run("lista vazia retorna slice vazio", func(t *testing.T) {
		mock := mustMock(t)
		defer func() { mock.Close() }()
		repo := NewPostgresPokemonRepositoryWithPool(mock)
		result, err := repo.GetByIDs(context.Background(), []string{})
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("com ids retorna pokemons", func(t *testing.T) {
		mock := mustMock(t)
		defer func() { mock.Close() }()

		cols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs([]string{"1", "25"}).
			WillReturnRows(pgxmock.NewRows(cols).
				AddRow("uuid-1", "Bulbasaur", "001", []string{"Grass"}, 0.7, 6.9, "seed", "https://img/1.png", "#78C850", "Grass", time.Now(), time.Now()).
				AddRow("uuid-25", "Pikachu", "025", []string{"Electric"}, 0.4, 6.0, "mouse", "https://img/25.png", "#F8D030", "Electric", time.Now(), time.Now()))

		repo := NewPostgresPokemonRepositoryWithPool(mock)
		result, err := repo.GetByIDs(context.Background(), []string{"1", "25"})
		require.NoError(t, err)
		assert.Len(t, result, 2)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresPokemonRepository_AddFavorite(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(pgxmock.PgxPoolIface)
		wantErr error
	}{
		{
			name: "sucesso",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO user_favorites")).
					WithArgs("user-1", "25").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
		},
		{
			name: "já existe retorna ErrFavoriteAlreadyExists",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO user_favorites")).
					WithArgs("user-1", "25").
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
			wantErr: ErrFavoriteAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mustMock(t)
			defer func() { mock.Close() }()
			tt.setup(mock)
			repo := NewPostgresPokemonRepositoryWithPool(mock)
			err := repo.AddFavorite(context.Background(), "user-1", "25")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresPokemonRepository_RemoveFavorite(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(pgxmock.PgxPoolIface)
		wantErr error
	}{
		{
			name: "sucesso",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta("DELETE FROM user_favorites")).
					WithArgs("user-1", "25").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
		{
			name: "nao existe retorna ErrFavoriteNotFound",
			setup: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta("DELETE FROM user_favorites")).
					WithArgs("user-1", "25").
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: ErrFavoriteNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mustMock(t)
			defer func() { mock.Close() }()
			tt.setup(mock)
			repo := NewPostgresPokemonRepositoryWithPool(mock)
			err := repo.RemoveFavorite(context.Background(), "user-1", "25")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresPokemonRepository_ListTypes(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, color")).
		WillReturnRows(pgxmock.NewRows([]string{"name", "color"}).
			AddRow("Fire", "#F08030").
			AddRow("Water", "#6890F0"))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	types, err := repo.ListTypes(context.Background())
	require.NoError(t, err)
	assert.Len(t, types, 2)
	assert.Equal(t, "Fire", types[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_ListRegions(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT r.id, r.name, g.id")).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "g_id"}).
			AddRow(int64(1), "Kanto", int64(1)).
			AddRow(int64(2), "Johto", int64(2)))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	regions, err := repo.ListRegions(context.Background())
	require.NoError(t, err)
	assert.Len(t, regions, 2)
	assert.Equal(t, "Kanto", regions[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInMemoryPokemonRepository_GetByID(t *testing.T) {
	repo := NewInMemoryPokemonRepository()

	p, err := repo.GetByID(context.Background(), "001")
	require.NoError(t, err)
	assert.Equal(t, "Bulbasaur", p.Name)

	_, err = repo.GetByID(context.Background(), "999")
	require.ErrorIs(t, err, ErrPokemonNotFound)
}

func TestInMemoryPokemonRepository_GetAll(t *testing.T) {
	repo := NewInMemoryPokemonRepository()

	page, err := repo.GetAll(context.Background(), 0, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), page.TotalElements)
	assert.Len(t, page.Content, 1)
	assert.True(t, page.HasNext)
}

func TestInMemoryPokemonRepository_Search(t *testing.T) {
	repo := NewInMemoryPokemonRepository()

	page, err := repo.Search(context.Background(), "bulb", 0, 20)
	require.NoError(t, err)
	assert.Len(t, page.Content, 1)
	assert.Equal(t, "Bulbasaur", page.Content[0].Name)
}

func TestInMemoryPokemonRepository_GetByType(t *testing.T) {
	repo := NewInMemoryPokemonRepository()

	page, err := repo.GetByType(context.Background(), "Grass", 0, 20)
	require.NoError(t, err)
	assert.Len(t, page.Content, 1)
}

func TestSanitizePage(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantPage   int
		wantSize   int
	}{
		{"negativo", -1, 20, 0, 20},
		{"zero size", 0, 0, 0, 20},
		{"muito grande", 0, 999, 0, 100},
		{"valido", 2, 10, 2, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotSize := sanitizePage(tt.page, tt.pageSize)
			assert.Equal(t, tt.wantPage, gotPage)
			assert.Equal(t, tt.wantSize, gotSize)
		})
	}
}
