package httphandler

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"
)

// GetHome retorna os dados da tela Home da Pokédex.
func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page, hasPage := getOptionalQueryParamInt(r, "page", 0)
	pageSize, hasSize := getOptionalQueryParamInt(r, "size", 20)
	userID := getUserIDFromContext(ctx)
	searchValue := strings.TrimSpace(r.URL.Query().Get("q"))
	selectedType := strings.TrimSpace(r.URL.Query().Get("type"))
	selectedOrdering := strings.TrimSpace(r.URL.Query().Get("order"))
	selectedRegion := strings.TrimSpace(r.URL.Query().Get("region"))
	paginate := hasPage

	if page < 0 {
		page = 0
	}

	if paginate {
		if !hasSize {
			pageSize = 20
		}
		if pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}

	pokemonPage, err := h.loadHomePokemonPage(ctx, userID, searchValue, selectedType)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter dados da home", "INTERNAL_ERROR")
		return
	}

	types, err := h.pokemonUseCase.ListTypes(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter filtros da home", "INTERNAL_ERROR")
		return
	}

	regions, err := h.pokemonUseCase.ListRegions(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao obter regioes da home", "INTERNAL_ERROR")
		return
	}

	filterHomePokemonByRegion(pokemonPage, selectedRegion)
	sortHomePokemonPage(pokemonPage, selectedOrdering)

	if paginate {
		paginateHomePokemonPage(pokemonPage, page, pageSize)
	}

	favoriteSet := h.buildFavoriteSet(ctx, userID)
	response := h.responseBuilder.BuildHomePageResponseWithTypes(
		pokemonPage,
		types,
		regions,
		favoriteSet,
		searchValue,
		selectedOrDefault(selectedType, "Todos os tipos"),
		selectedOrDefault(selectedOrdering, "Menor número"),
		selectedRegion,
	)
	RespondJSON(w, http.StatusOK, response)
}

// GetRegions retorna a lista de regiões disponíveis.
func (h *Handler) GetRegions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	regions, err := h.pokemonUseCase.ListRegions(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "falha ao listar regioes", "INTERNAL_ERROR")
		return
	}

	RespondJSON(w, http.StatusOK, h.responseBuilder.BuildRegionsResponse(regions))
}

func (h *Handler) loadHomePokemonPage(
	ctx context.Context,
	userID string,
	searchValue string,
	selectedType string,
) (*domain.PokemonPage, error) {
	searchTerm := strings.TrimSpace(searchValue)
	typeFilter := strings.TrimSpace(selectedType)

	const page = 0
	const pageSize = 200

	var result *domain.PokemonPage
	var err error

	if searchTerm != "" {
		result, err = h.pokemonUseCase.SearchPokemons(ctx, searchTerm, page, pageSize, userID)
	} else if typeFilter != "" && typeFilter != "Todos os tipos" {
		result, err = h.pokemonUseCase.FilterByType(ctx, typeFilter, page, pageSize, userID)
	} else {
		result, err = h.pokemonUseCase.ListPokemons(ctx, page, pageSize, userID)
	}

	if err != nil {
		return nil, err
	}

	if result != nil && typeFilter != "" && typeFilter != "Todos os tipos" && searchTerm != "" {
		filtered := make([]domain.Pokemon, 0, len(result.Content))
		for _, pokemon := range result.Content {
			if hasPokemonType(pokemon.Types, typeFilter) {
				filtered = append(filtered, pokemon)
			}
		}
		result.Content = filtered
	}

	return result, nil
}

func paginateHomePokemonPage(page *domain.PokemonPage, currentPage, pageSize int) {
	if page == nil {
		return
	}

	if currentPage < 0 {
		currentPage = 0
	}
	if pageSize < 1 {
		pageSize = 20
	}

	start := currentPage * pageSize
	if start >= len(page.Content) {
		page.Content = []domain.Pokemon{}
		return
	}

	end := start + pageSize
	if end > len(page.Content) {
		end = len(page.Content)
	}

	page.Content = page.Content[start:end]
}

func sortHomePokemonPage(page *domain.PokemonPage, selectedOrdering string) {
	if page == nil {
		return
	}

	switch selectedOrdering {
	case "Maior número":
		sort.Slice(page.Content, func(i, j int) bool { return parseNumber(page.Content[i].Number) > parseNumber(page.Content[j].Number) })
	case "A-Z":
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Name < page.Content[j].Name })
	case "Z-A":
		sort.Slice(page.Content, func(i, j int) bool { return page.Content[i].Name > page.Content[j].Name })
	default:
		sort.Slice(page.Content, func(i, j int) bool { return parseNumber(page.Content[i].Number) < parseNumber(page.Content[j].Number) })
	}
}

func parseNumber(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

func filterHomePokemonByRegion(page *domain.PokemonPage, region string) {
	if page == nil || strings.TrimSpace(region) == "" {
		return
	}

	filtered := make([]domain.Pokemon, 0, len(page.Content))
	for _, pokemon := range page.Content {
		if matchesRegion(pokemon.Number, region) {
			filtered = append(filtered, pokemon)
		}
	}
	page.Content = filtered
}

func matchesRegion(number string, region string) bool {
	parsed, ok := parsePokemonNumber(number)
	if !ok {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(region)) {
	case "kanto":
		return parsed >= 1 && parsed <= 151
	case "johto":
		return parsed >= 152 && parsed <= 251
	case "hoenn":
		return parsed >= 252 && parsed <= 386
	case "sinnoh":
		return parsed >= 387 && parsed <= 493
	case "unova":
		return parsed >= 494 && parsed <= 649
	case "kalos":
		return parsed >= 650 && parsed <= 721
	case "alola":
		return parsed >= 722 && parsed <= 809
	case "galar":
		return parsed >= 810 && parsed <= 905
	default:
		return true
	}
}

func parsePokemonNumber(number string) (int, bool) {
	normalized := normalizePokemonID(number)
	parsed, err := strconv.Atoi(normalized)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func selectedOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
