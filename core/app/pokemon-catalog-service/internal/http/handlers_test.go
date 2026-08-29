package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/domain"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/repository"
)

type stubPinger struct {
	err error
}

func (s stubPinger) Ping(ctx context.Context) error {
	return s.err
}

type stubPokemonRepo struct {
	getAllFn         func(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	searchFn         func(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	getByTypeFn      func(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	getByIDFn        func(ctx context.Context, id string) (*domain.Pokemon, error)
	getByIDsFn       func(ctx context.Context, ids []string) ([]domain.Pokemon, error)
	getDetailByIDFn  func(ctx context.Context, id string) (*domain.PokemonDetail, error)
	listTypesFn      func(ctx context.Context) ([]domain.Type, error)
	listRegionsFn    func(ctx context.Context) ([]domain.Region, error)
	addFavoriteFn    func(ctx context.Context, userID, pokemonID string) error
	removeFavoriteFn func(ctx context.Context, userID, pokemonID string) error
}

func (s *stubPokemonRepo) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, page, pageSize)
	}
	return &domain.PokemonPage{Content: []domain.Pokemon{}}, nil
}

func (s *stubPokemonRepo) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	if s.searchFn != nil {
		return s.searchFn(ctx, query, page, pageSize)
	}
	return &domain.PokemonPage{Content: []domain.Pokemon{}}, nil
}

func (s *stubPokemonRepo) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	if s.getByTypeFn != nil {
		return s.getByTypeFn(ctx, typeFilter, page, pageSize)
	}
	return &domain.PokemonPage{Content: []domain.Pokemon{}}, nil
}

func (s *stubPokemonRepo) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return &domain.Pokemon{ID: "1", Name: "Bulbasaur", Number: "001"}, nil
}

func (s *stubPokemonRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
	if s.getByIDsFn != nil {
		return s.getByIDsFn(ctx, ids)
	}
	return []domain.Pokemon{}, nil
}

func (s *stubPokemonRepo) GetDetailByID(ctx context.Context, id string) (*domain.PokemonDetail, error) {
	if s.getDetailByIDFn != nil {
		return s.getDetailByIDFn(ctx, id)
	}
	return &domain.PokemonDetail{Name: "Pikachu", Number: "025"}, nil
}

func (s *stubPokemonRepo) ListTypes(ctx context.Context) ([]domain.Type, error) {
	if s.listTypesFn != nil {
		return s.listTypesFn(ctx)
	}
	return []domain.Type{}, nil
}

func (s *stubPokemonRepo) ListRegions(ctx context.Context) ([]domain.Region, error) {
	if s.listRegionsFn != nil {
		return s.listRegionsFn(ctx)
	}
	return []domain.Region{}, nil
}

func (s *stubPokemonRepo) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	if s.addFavoriteFn != nil {
		return s.addFavoriteFn(ctx, userID, pokemonID)
	}
	return nil
}

func (s *stubPokemonRepo) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	if s.removeFavoriteFn != nil {
		return s.removeFavoriteFn(ctx, userID, pokemonID)
	}
	return nil
}

func pokemonPage(items []domain.Pokemon) *domain.PokemonPage {
	return &domain.PokemonPage{
		Content:       items,
		TotalElements: int64(len(items)),
		CurrentPage:   0,
		TotalPages:    1,
		HasNext:       false,
	}
}

func TestListPokemons(t *testing.T) {
	repo := &stubPokemonRepo{
		getAllFn: func(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
			return pokemonPage([]domain.Pokemon{
				{ID: "1", Name: "Bulbasaur", Number: "001", Types: []string{"Grass", "Poison"}},
				{ID: "25", Name: "Pikachu", Number: "025", Types: []string{"Electric"}},
			}), nil
		},
	}
	mux := NewMux(repo)

	req := httptest.NewRequest(http.MethodGet, "/v1/pokemons", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var page domain.PokemonPage
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatalf("erro ao decodificar: %v", err)
	}
	if len(page.Content) != 2 {
		t.Errorf("esperava 2 pokemons, obteve %d", len(page.Content))
	}
}

func TestGetPokemonByID(t *testing.T) {
	repo := &stubPokemonRepo{
		getByIDFn: func(ctx context.Context, id string) (*domain.Pokemon, error) {
			if id == "999" {
				return nil, repository.ErrPokemonNotFound
			}
			return &domain.Pokemon{ID: id, Name: "Bulbasaur", Number: "001"}, nil
		},
	}
	mux := NewMux(repo)

	t.Run("existente", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("nao encontrado", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/999", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})
}

func TestSearchPokemons(t *testing.T) {
	repo := &stubPokemonRepo{
		searchFn: func(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
			return pokemonPage([]domain.Pokemon{
				{ID: "25", Name: "Pikachu", Number: "025"},
			}), nil
		},
	}
	mux := NewMux(repo)

	req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/search?q=pikachu", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestGetFavoritesBatch(t *testing.T) {
	repo := &stubPokemonRepo{
		getByIDsFn: func(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
			return []domain.Pokemon{
				{ID: "1", Name: "Bulbasaur", Number: "001"},
				{ID: "25", Name: "Pikachu", Number: "025"},
			}, nil
		},
	}
	mux := NewMux(repo)

	t.Run("com ids", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/favorites?ids=1,25", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}

		var items []domain.Pokemon
		if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
			t.Fatalf("erro ao decodificar: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("esperava 2 pokemons, obteve %d", len(items))
		}
	})

	t.Run("lista vazia", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/favorites?ids=", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}

		var items []domain.Pokemon
		if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
			t.Fatalf("erro ao decodificar: %v", err)
		}
		if len(items) != 0 {
			t.Errorf("esperava lista vazia, obteve %d", len(items))
		}
	})
}

func TestHealthHandler(t *testing.T) {
	repo := &stubPokemonRepo{}
	mux := NewMux(repo)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestListTypesAndRegions(t *testing.T) {
	repo := &stubPokemonRepo{
		listTypesFn: func(ctx context.Context) ([]domain.Type, error) {
			return []domain.Type{{Name: "Fire", Color: "#F08030"}}, nil
		},
		listRegionsFn: func(ctx context.Context) ([]domain.Region, error) {
			return []domain.Region{{ID: "kanto", Name: "Kanto"}}, nil
		},
	}
	mux := NewMux(repo)

	t.Run("tipos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/types", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("regioes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/regions", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})
}

func TestFilterByType(t *testing.T) {
	repo := &stubPokemonRepo{
		getByTypeFn: func(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
			return pokemonPage([]domain.Pokemon{
				{ID: "4", Name: "Charmander", Number: "004", Types: []string{"Fire"}},
			}), nil
		},
	}
	mux := NewMux(repo)

	t.Run("com tipo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/type/Fire", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var page domain.PokemonPage
		if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
			t.Fatalf("erro ao decodificar: %v", err)
		}
		if len(page.Content) != 1 {
			t.Errorf("esperava 1 pokemon, obteve %d", len(page.Content))
		}
	})
}

func TestGetPokemonDetailByID(t *testing.T) {
	repo := &stubPokemonRepo{
		getDetailByIDFn: func(ctx context.Context, id string) (*domain.PokemonDetail, error) {
			if id == "999" {
				return nil, repository.ErrPokemonNotFound
			}
			return &domain.PokemonDetail{ID: "1", Name: "Bulbasaur", Number: "001"}, nil
		},
	}
	mux := NewMux(repo)

	t.Run("existente", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemon-details/1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("nao encontrado", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemon-details/999", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("erro de repositorio", func(t *testing.T) {
		repo2 := &stubPokemonRepo{
			getDetailByIDFn: func(ctx context.Context, id string) (*domain.PokemonDetail, error) {
				return nil, errors.New("db down")
			},
		}
		mux2 := NewMux(repo2)
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemon-details/1", nil)
		w := httptest.NewRecorder()
		mux2.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}

func TestAddFavoriteHandler(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		addErr     error
		wantStatus int
	}{
		{"sem autenticacao", "", nil, http.StatusUnauthorized},
		{"sucesso", "user-1", nil, http.StatusOK},
		{"favorito ja existe", "user-1", repository.ErrFavoriteAlreadyExists, http.StatusConflict},
		{"erro de repositorio", "user-1", errors.New("db down"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubPokemonRepo{
				addFavoriteFn: func(ctx context.Context, userID, pokemonID string) error {
					return tt.addErr
				},
			}
			mux := NewMux(repo)
			req := httptest.NewRequest(http.MethodPost, "/v1/pokemons/1/favorite", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRemoveFavoriteHandler(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		removeErr  error
		wantStatus int
	}{
		{"sem autenticacao", "", nil, http.StatusUnauthorized},
		{"sucesso", "user-1", nil, http.StatusOK},
		{"favorito nao encontrado", "user-1", repository.ErrFavoriteNotFound, http.StatusNotFound},
		{"erro de repositorio", "user-1", errors.New("db down"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubPokemonRepo{
				removeFavoriteFn: func(ctx context.Context, userID, pokemonID string) error {
					return tt.removeErr
				},
			}
			mux := NewMux(repo)
			req := httptest.NewRequest(http.MethodDelete, "/v1/pokemons/1/favorite", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestPingHandler(t *testing.T) {
	repo := &stubPokemonRepo{}
	mux := NewMux(repo)
	req := httptest.NewRequest(http.MethodGet, "/v1/pokemon/ping", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
	}{
		{"banco pronto", nil, http.StatusOK},
		{"banco indisponivel", errors.New("connection refused"), http.StatusServiceUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadyHandler(stubPinger{err: tt.pingErr})
			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			w := httptest.NewRecorder()
			h(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestSearchPokemonsErro(t *testing.T) {
	repo := &stubPokemonRepo{
		searchFn: func(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
			return nil, errors.New("db down")
		},
	}
	mux := NewMux(repo)

	t.Run("sem query retorna 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/search", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("erro de repositorio retorna 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pokemons/search?q=pika", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}
