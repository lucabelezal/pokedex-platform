package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pokedex-platform/core/app/pokemon-catalog-service/internal/domain"
	"pokedex-platform/core/app/pokemon-catalog-service/internal/repository"
)

type stubPokemonRepo struct {
	getAllFn        func(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)
	searchFn        func(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error)
	getByTypeFn     func(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error)
	getByIDFn       func(ctx context.Context, id string) (*domain.Pokemon, error)
	getByIDsFn      func(ctx context.Context, ids []string) ([]domain.Pokemon, error)
	getDetailByIDFn func(ctx context.Context, id string) (*domain.PokemonDetail, error)
	listTypesFn     func(ctx context.Context) ([]domain.Type, error)
	listRegionsFn   func(ctx context.Context) ([]domain.Region, error)
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
