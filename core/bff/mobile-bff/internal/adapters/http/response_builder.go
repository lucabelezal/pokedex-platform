package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"
	"strconv"
	"strings"
)

type ResponseBuilder struct{}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

func (rb *ResponseBuilder) BuildPokemonDetailDTO(p *domain.PokemonDetail) *dto.PokemonDetailDTO {
	types := make([]dto.TypeDTO, len(p.Types))
	for i, t := range p.Types {
		types[i] = dto.TypeDTO{
			Name:  t.Name,
			Color: t.Color,
		}
	}

	return &dto.PokemonDetailDTO{
		Number: p.Number,
		Name:   p.Name,
		Image: dto.ImageDTO{
			URL: p.ImageURL,
			Element: dto.ElementDTO{
				Color: p.Element.Color,
				Type:  p.Element.Type,
			},
		},
		Types:       types,
		Height:      p.Height,
		Weight:      p.Weight,
		Description: p.Description,
		IsFavorite:  p.IsFavorite,
	}
}

func (rb *ResponseBuilder) BuildRichPokemonResponse(p *domain.Pokemon) *dto.RichPokemonResponse {
	types := make([]dto.TypeDTO, len(p.Types))
	for i, t := range p.Types {
		color := getTypeColor(t)
		types[i] = dto.TypeDTO{
			Name:  t,
			Color: color,
		}
	}

	return &dto.RichPokemonResponse{
		Number: p.Number,
		Name:   p.Name,
		Image: dto.ImageDTO{
			URL: p.ImageURL,
			Element: dto.ElementDTO{
				Color: p.ElementColor,
				Type:  p.ElementType,
			},
		},
		Types: types,
	}
}

func (rb *ResponseBuilder) BuildRichPokemonListResponse(page *domain.PokemonPage) *dto.RichPokemonListResponse {
	pokemons := make([]dto.RichPokemonResponse, len(page.Content))
	for i, p := range page.Content {
		rich := rb.BuildRichPokemonResponse(&p)
		pokemons[i] = *rich
	}

	return &dto.RichPokemonListResponse{
		Content:       pokemons,
		TotalElements: page.TotalElements,
		CurrentPage:   page.CurrentPage,
		TotalPages:    page.TotalPages,
		HasNext:       page.HasNext,
		Search: dto.SearchMetadata{
			Placeholder: "Procure por um Pokémon...",
		},
		Filters: []interface{}{},
	}
}

func (rb *ResponseBuilder) BuildHomeResponse(page *domain.PokemonPage) *dto.RichPokemonListResponse {
	pokemons := make([]dto.RichPokemonResponse, len(page.Content))
	for i, p := range page.Content {
		rich := rb.BuildRichPokemonResponse(&p)
		pokemons[i] = *rich
	}

	return &dto.RichPokemonListResponse{
		Content:       pokemons,
		TotalElements: page.TotalElements,
		CurrentPage:   page.CurrentPage,
		TotalPages:    page.TotalPages,
		HasNext:       page.HasNext,
		Search: dto.SearchMetadata{
			Placeholder: "Procure por um Pokémon...",
		},
		Filters: []interface{}{},
	}
}

func (rb *ResponseBuilder) BuildHomePageResponse(page *domain.PokemonPage) *dto.HomeResponse {
	return rb.BuildHomePageResponseWithTypes(page, nil, nil)
}

func (rb *ResponseBuilder) BuildHomePageResponseWithTypes(
	page *domain.PokemonPage,
	types []domain.Type,
	favoriteSet map[string]struct{},
) *dto.HomeResponse {
	pokemons := make([]dto.HomePokemonDTO, len(page.Content))
	for i, p := range page.Content {
		pokemons[i] = rb.BuildHomePokemonDTO(&p, favoriteSet)
	}

	typeItems := make([]dto.HomeFilterItemDTO, len(types))
	for i, t := range types {
		typeItems[i] = dto.HomeFilterItemDTO{Title: t.Name}
	}

	return &dto.HomeResponse{
		Search: dto.HomeSearchDTO{
			Placeholder: "Procurar Pokémon...",
		},
		Filters: dto.HomeFiltersDTO{
			Types: dto.HomeFilterGroupDTO{
				Title: "Todos os tipos",
				Items: typeItems,
			},
			Ordering: dto.HomeFilterGroupDTO{
				Title: "Ordenação",
				Items: []dto.HomeFilterItemDTO{
					{Title: "Menor número"},
				},
			},
		},
		Pokemons: pokemons,
	}
}

func (rb *ResponseBuilder) BuildHomePokemonDTO(
	p *domain.Pokemon,
	favoriteSet map[string]struct{},
) dto.HomePokemonDTO {
	types := make([]dto.HomePokemonTypeDTO, len(p.Types))
	for i, t := range p.Types {
		types[i] = dto.HomePokemonTypeDTO{
			Title: t,
			Color: normalizeHexColor(getTypeColor(t)),
		}
	}

	_, isFavorite := favoriteSet[normalizePokemonID(p.Number)]

	return dto.HomePokemonDTO{
		Number: formatHomePokemonNumber(p.Number),
		Name:   p.Name,
		Types:  types,
		Sprites: dto.HomePokemonSpritesDTO{
			URL:             p.ImageURL,
			BackgroundColor: normalizeHexColor(p.ElementColor),
		},
		IsFavorite: isFavorite,
	}
}

func (rb *ResponseBuilder) BuildHealthResponse() *dto.HealthResponse {
	return &dto.HealthResponse{
		Status:  "ok",
		Service: "mobile-bff",
	}
}

func getTypeColor(typeStr string) string {
	typeColors := map[string]string{
		"Normal":    "#A8A878",
		"Fogo":      "#EE8130",
		"Fire":      "#F08030",
		"Água":      "#6390F0",
		"Water":     "#6890F0",
		"Elétrico":  "#F7D02C",
		"Electric":  "#F8D030",
		"Grama":     "#7AC74C",
		"Grass":     "#78C850",
		"Gelo":      "#96D9D6",
		"Ice":       "#98D8D8",
		"Lutador":   "#C22E28",
		"Fighting":  "#C03028",
		"Venenoso":  "#A33EA1",
		"Poison":    "#A040A0",
		"Terrestre": "#E2BF65",
		"Ground":    "#E0C068",
		"Voador":    "#A98FF3",
		"Flying":    "#A890F0",
		"Psíquico":  "#F95587",
		"Psychic":   "#F85888",
		"Inseto":    "#A6B91A",
		"Bug":       "#A8B820",
		"Pedra":     "#B6A136",
		"Rock":      "#B8A038",
		"Fantasma":  "#735797",
		"Ghost":     "#705898",
		"Dragão":    "#6F35FC",
		"Dragon":    "#7038F8",
		"Sombrio":   "#705746",
		"Dark":      "#705848",
		"Aço":       "#B7B7CE",
		"Steel":     "#B8B8D0",
		"Fada":      "#D685AD",
		"Fairy":     "#EE99AC",
	}

	if color, exists := typeColors[typeStr]; exists {
		return color
	}
	return "#A9AC86"
}

func normalizeHexColor(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "#")
}

func formatHomePokemonNumber(number string) string {
	trimmed := strings.TrimSpace(number)
	if parsed, err := strconv.Atoi(trimmed); err == nil {
		return fmt.Sprintf("Nº%03d", parsed)
	}

	normalized := strings.TrimLeft(trimmed, "0")
	if normalized == "" {
		normalized = "0"
	}

	if parsed, err := strconv.Atoi(normalized); err == nil {
		return fmt.Sprintf("Nº%03d", parsed)
	}

	return "Nº" + trimmed
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func RespondError(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := dto.ErrorResponse{
		Error:   code,
		Message: message,
		Code:    status,
	}
	_ = json.NewEncoder(w).Encode(err)
}
