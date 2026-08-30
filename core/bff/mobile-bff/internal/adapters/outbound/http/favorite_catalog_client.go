package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

// FavoriteCatalogClient busca detalhes de Pokémon favoritos via catalog-service.
type FavoriteCatalogClient struct {
	baseURL string
	client  *CircuitBreakerClient
}

// NewFavoriteCatalogClient cria um cliente para o endpoint de batch favorites.
func NewFavoriteCatalogClient(baseURL string) *FavoriteCatalogClient {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	cfg := DefaultCircuitBreakerConfig("pokemon-catalog-service-favorites")
	return &FavoriteCatalogClient{
		baseURL: trimmed,
		client: NewCircuitBreakerClient(&http.Client{
			Timeout: 5 * time.Second,
		}, cfg),
	}
}

// GetFavoriteDetails busca detalhes de múltiplos Pokémon por seus IDs.
func (c *FavoriteCatalogClient) GetFavoriteDetails(ctx context.Context, ids []string) ([]domain.Pokemon, error) {
	if len(ids) == 0 {
		return []domain.Pokemon{}, nil
	}

	idsParam := strings.Join(ids, ",")
	endpoint := fmt.Sprintf("%s/v1/pokemons/favorites?ids=%s", c.baseURL, url.QueryEscape(idsParam))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisicao de batch favorites: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar favoritos em batch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("catalog-service retornou status %d para batch favorites", resp.StatusCode)
	}

	var result []domain.Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta de batch favorites: %w", err)
	}

	return result, nil
}

// AddFavorite adiciona um Pokémon aos favoritos via catalog-service.
func (c *FavoriteCatalogClient) AddFavorite(ctx context.Context, userID, pokemonID string) error {
	endpoint := fmt.Sprintf("%s/v1/pokemons/%s/favorite", c.baseURL, url.PathEscape(pokemonID))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisicao de adicionar favorito: %w", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao adicionar favorito: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusConflict:
		return domain.ErrFavoriteAlreadyExists
	default:
		return fmt.Errorf("catalog-service retornou status %d ao adicionar favorito", resp.StatusCode)
	}
}

// RemoveFavorite remove um Pokémon dos favoritos via catalog-service.
func (c *FavoriteCatalogClient) RemoveFavorite(ctx context.Context, userID, pokemonID string) error {
	endpoint := fmt.Sprintf("%s/v1/pokemons/%s/favorite", c.baseURL, url.PathEscape(pokemonID))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisicao de remover favorito: %w", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao remover favorito: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return domain.ErrFavoriteNotFound
	default:
		return fmt.Errorf("catalog-service retornou status %d ao remover favorito", resp.StatusCode)
	}
}

// GetUserFavorites lista os IDs de Pokémon favoritos do usuário via catalog-service.
func (c *FavoriteCatalogClient) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
	if strings.TrimSpace(userID) == "" {
		return []string{}, nil
	}

	endpoint := fmt.Sprintf("%s/v1/favorites?user_id=%s", c.baseURL, url.QueryEscape(userID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisicao de listar favoritos: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar favoritos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return []string{}, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("catalog-service retornou status %d ao listar favoritos", resp.StatusCode)
	}

	var ids []string
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("erro ao decodificar lista de favoritos: %w", err)
	}
	if ids == nil {
		ids = []string{}
	}

	return ids, nil
}

// IsFavorite verifica se um Pokémon está nos favoritos do usuário via catalog-service.
func (c *FavoriteCatalogClient) IsFavorite(ctx context.Context, userID, pokemonID string) (bool, error) {
	ids, err := c.GetUserFavorites(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if id == pokemonID {
			return true, nil
		}
	}
	return false, nil
}

var _ outbound.FavoriteRepository = (*FavoriteCatalogClient)(nil)
var _ outbound.FavoriteCatalogProvider = (*FavoriteCatalogClient)(nil)
