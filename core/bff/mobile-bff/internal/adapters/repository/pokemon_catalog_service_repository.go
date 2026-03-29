package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

type PokemonCatalogServiceRepository struct {
	baseURL string
	client  *http.Client
}

func NewPokemonCatalogServiceRepository(baseURL string) *PokemonCatalogServiceRepository {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	return &PokemonCatalogServiceRepository{
		baseURL: trimmed,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (r *PokemonCatalogServiceRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
	path := fmt.Sprintf("%s/v1/pokemons/%s", r.baseURL, url.PathEscape(id))

	var out domain.Pokemon
	status, err := r.getJSON(ctx, path, &out)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, domain.ErrPokemonNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (r *PokemonCatalogServiceRepository) GetDetailByID(ctx context.Context, id string) (*domain.PokemonScreenDetail, error) {
	path := fmt.Sprintf("%s/v1/pokemon-details/%s", r.baseURL, url.PathEscape(id))

	var out domain.PokemonScreenDetail
	status, err := r.getJSON(ctx, path, &out)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, domain.ErrPokemonNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (r *PokemonCatalogServiceRepository) GetAll(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
	endpoint := fmt.Sprintf("%s/v1/pokemons?page=%d&size=%d", r.baseURL, page, pageSize)
	var out domain.PokemonPage
	_, err := r.getJSON(ctx, endpoint, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PokemonCatalogServiceRepository) Search(ctx context.Context, query string, page, pageSize int) (*domain.PokemonPage, error) {
	endpoint := fmt.Sprintf("%s/v1/pokemons/search?q=%s&page=%d&size=%d", r.baseURL, url.QueryEscape(query), page, pageSize)
	var out domain.PokemonPage
	_, err := r.getJSON(ctx, endpoint, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PokemonCatalogServiceRepository) GetByType(ctx context.Context, typeFilter string, page, pageSize int) (*domain.PokemonPage, error) {
	endpoint := fmt.Sprintf("%s/v1/pokemons/type/%s?page=%d&size=%d", r.baseURL, url.PathEscape(typeFilter), page, pageSize)
	var out domain.PokemonPage
	_, err := r.getJSON(ctx, endpoint, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PokemonCatalogServiceRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	endpoint := fmt.Sprintf("%s/v1/types", r.baseURL)
	var out []domain.Type
	_, err := r.getJSON(ctx, endpoint, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PokemonCatalogServiceRepository) ListRegions(ctx context.Context) ([]domain.Region, error) {
	endpoint := fmt.Sprintf("%s/v1/regions", r.baseURL)
	var out []domain.Region
	_, err := r.getJSON(ctx, endpoint, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PokemonCatalogServiceRepository) GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]string, error) {
	_ = ctx
	_ = userID
	_ = page
	_ = pageSize
	return []string{}, nil
}

func (r *PokemonCatalogServiceRepository) getJSON(ctx context.Context, endpoint string, out any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("pokemon-catalog-service retornou status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}
