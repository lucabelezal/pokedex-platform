package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/domain"

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
		name     string
		page     int
		pageSize int
		wantPage int
		wantSize int
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

func TestMapTypes(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		wantLen  int
		wantName string
	}{
		{"nomes normais", []string{"Fire", "Water"}, 2, "Fire"},
		{"acerto aço", []string{"Aço"}, 1, "Metal"},
		{"acerto sombrio", []string{"Sombrio"}, 1, "Noturno"},
		{"vazio", []string{}, 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types := mapTypes(tt.names)
			assert.Len(t, types, tt.wantLen)
			if tt.wantLen > 0 {
				assert.Equal(t, tt.wantName, types[0].Name)
				assert.NotEmpty(t, types[0].Color)
			}
		})
	}
}

func TestMapTypesWithOverrides(t *testing.T) {
	types := mapTypesWithOverrides("6", []string{"Fire", "Flying"})
	assert.Len(t, types, 2)
}

func TestNormalizeTypeName(t *testing.T) {
	assert.Equal(t, "Metal", normalizeTypeName("Aço"))
	assert.Equal(t, "Noturno", normalizeTypeName("Sombrio"))
	assert.Equal(t, "Fire", normalizeTypeName("  Fire  "))
}

func TestNormalizeCategory(t *testing.T) {
	assert.Equal(t, "Seed", normalizeCategory("001", "Seed Pokémon"))
	assert.Equal(t, "", normalizeCategory("999", ""))
	assert.Equal(t, "Rato", normalizeCategory("999", "Rato Pokémon"))
}

func TestNormalizeCondition(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Subir nivel com felicidade", "Nível de Amizade"},
		{"Uso de thunder stone", "Pedra do Trovão"},
		{"Uso de moon stone", "Pedra da Lua"},
		{"Uso de dusk stone", "Pedra do Anoitecer"},
		{"Troca", "Trocas"},
		{"Condicao nao mapeada", "Subir de Nível c/ Rollout"},
		{"Nivel 32", "Nível 36"},
		{"Nivel 16", "Nível 16"},
		{"qualquer coisa", "qualquer coisa"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeCondition(tt.in))
		})
	}
}

func TestNormalizeDescription(t *testing.T) {
	assert.Equal(t, "descrição", normalizeDescription("999", "  descrição  "))
}

func TestNormalizeAbilities(t *testing.T) {
	assert.Equal(t, []string{"Overgrow"}, normalizeAbilities("001", []string{"Overgrow"}))
}

func TestTypeColor(t *testing.T) {
	assert.Equal(t, "#EE8130", typeColor("Fogo"))
	assert.Equal(t, "#A9AC86", typeColor("unknown"))
}

func TestGenerationLabel(t *testing.T) {
	assert.Equal(t, "1º Geração", generationLabel(1))
}

func TestPostgresPokemonRepository_GetDetailByID(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	baseCols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("010").
		WillReturnRows(pgxmock.NewRows(baseCols).
			AddRow("uuid-10", "Caterpie", "010", []string{"Inseto"}, 0.3, 2.9, "verme", "https://img/10.png", "#A6B91A", "Inseto", time.Now(), time.Now()))

	detailCols := []string{"category", "gender_male", "gender_female", "region", "generation", "evolution_chain_id"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("010").
		WillReturnRows(pgxmock.NewRows(detailCols).AddRow("Verme", -1.0, -1.0, "Kanto", "Geração I", int64(1)))

	abilityCols := []string{"name"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT a.name")).
		WithArgs("010").
		WillReturnRows(pgxmock.NewRows(abilityCols).AddRow("Poeira de Escama"))

	weakCols := []string{"name", "color"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT t.name, t.color")).
		WithArgs("010").
		WillReturnRows(pgxmock.NewRows(weakCols).AddRow("Fogo", "#EE8130"))

	chainCols := []string{"chain_data"}
	chainJSON := `{"pokemon":{"id":10,"name":"Caterpie"},"condition":{"description":""},"evolutions_to":[{"pokemon":{"id":11,"name":"Metapod"},"condition":{"description":"Nivel 7"},"evolutions_to":[]}]}`
	mock.ExpectQuery(regexp.QuoteMeta("SELECT chain_data")).
		WithArgs(int64(1)).
		WillReturnRows(pgxmock.NewRows(chainCols).AddRow(chainJSON))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("11").
		WillReturnRows(pgxmock.NewRows(baseCols).
			AddRow("uuid-11", "Metapod", "011", []string{"Inseto"}, 0.7, 9.9, "casulo", "https://img/11.png", "#A6B91A", "Inseto", time.Now(), time.Now()))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	detail, err := repo.GetDetailByID(context.Background(), "010")
	require.NoError(t, err)
	assert.Equal(t, "Caterpie", detail.Name)
	assert.NotNil(t, detail)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_GetDetailByID_NotFounded(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	baseCols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("999").
		WillReturnRows(pgxmock.NewRows(baseCols))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	_, err := repo.GetDetailByID(context.Background(), "999")
	require.ErrorIs(t, err, ErrPokemonNotFound)
}

func TestFlattenEvolutionChain(t *testing.T) {
	root := evolutionNode{
		Pokemon: struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}{ID: 1, Name: "Bulbasaur"},
		Condition: struct {
			Description string `json:"description"`
		}{Description: ""},
		EvolutionsTo: []evolutionNode{{
			Pokemon: struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}{ID: 2, Name: "Ivysaur"},
			Condition: struct {
				Description string `json:"description"`
			}{Description: "Nivel 16"},
			EvolutionsTo: nil,
		}},
	}

	steps := make([]evolutionStep, 0)
	flattenEvolutionChain(root, "", &steps)
	assert.Len(t, steps, 2)
	assert.Equal(t, int64(1), steps[0].ID)
	assert.Equal(t, "Nivel 16", steps[1].Trigger)
}

func TestBuildPage(t *testing.T) {
	page := buildPage([]domain.Pokemon{{ID: "1"}}, 25, 0, 10)
	assert.Equal(t, int64(25), page.TotalElements)
	assert.Equal(t, 3, page.TotalPages)
	assert.True(t, page.HasNext)
}

func TestPostgresPokemonRepository_BuildEvolutionOverrides_Encontrado(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	baseCols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("1").
		WillReturnRows(pgxmock.NewRows(baseCols).
			AddRow("uuid-1", "Bulbasaur", "001", []string{"Grama", "Venenoso"}, 0.7, 6.9, "seed", "https://img/1.png", "#78C850", "Grama", time.Now(), time.Now()))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	items := []evolutionOverride{
		{Number: "1", Name: "Bulbasaur", Types: []string{"Grama", "Venenoso"}},
	}
	evos := repo.buildEvolutionOverrides(context.Background(), items)
	require.Len(t, evos, 1)
	assert.Equal(t, "Bulbasaur", evos[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_BuildEvolutionOverrides_NaoEncontrado(t *testing.T) {
	mock := mustMock(t)
	defer func() { mock.Close() }()

	baseCols := []string{"id", "name", "number", "types", "height", "weight", "description", "image_url", "element_color", "element_type", "created_at", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs("999").
		WillReturnRows(pgxmock.NewRows(baseCols))

	repo := NewPostgresPokemonRepositoryWithPool(mock)
	items := []evolutionOverride{
		{Number: "999", Name: "Fake", Types: []string{"Fogo"}},
	}
	evos := repo.buildEvolutionOverrides(context.Background(), items)
	require.Len(t, evos, 1)
	assert.Equal(t, "Fake", evos[0].Name)
	// imagem fallback gerada a partir do número
	assert.Contains(t, evos[0].ImageURL, "999")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNewPool_ComURIVazia(t *testing.T) {
	_, err := NewPool(context.Background(), "")
	require.Error(t, err)
}

func TestPostgresPokemonRepository_GetUserFavoriteIDs(t *testing.T) {
	t.Run("user vazio retorna lista vazia", func(t *testing.T) {
		mock := mustMock(t)
		defer func() { mock.Close() }()
		repo := NewPostgresPokemonRepositoryWithPool(mock)
		ids, err := repo.GetUserFavoriteIDs(context.Background(), "")
		require.NoError(t, err)
		assert.Empty(t, ids)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("com favoritos", func(t *testing.T) {
		mock := mustMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT pokemon_id::TEXT")).
			WithArgs("user-1").
			WillReturnRows(pgxmock.NewRows([]string{"pokemon_id"}).AddRow("1").AddRow("25"))

		repo := NewPostgresPokemonRepositoryWithPool(mock)
		ids, err := repo.GetUserFavoriteIDs(context.Background(), "user-1")
		require.NoError(t, err)
		assert.Equal(t, []string{"1", "25"}, ids)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sem favoritos retorna lista vazia", func(t *testing.T) {
		mock := mustMock(t)
		defer func() { mock.Close() }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT pokemon_id::TEXT")).
			WithArgs("user-1").
			WillReturnRows(pgxmock.NewRows([]string{"pokemon_id"}))

		repo := NewPostgresPokemonRepositoryWithPool(mock)
		ids, err := repo.GetUserFavoriteIDs(context.Background(), "user-1")
		require.NoError(t, err)
		assert.Empty(t, ids)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
