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
