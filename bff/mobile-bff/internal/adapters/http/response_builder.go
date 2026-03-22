package http

import (
	"encoding/json"
	"net/http"
	"pokedex-platform/bff/mobile-bff/internal/adapters/http/dto"
	"pokedex-platform/bff/mobile-bff/internal/domain"
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
	richList := rb.BuildHomeResponse(page)

	return &dto.HomeResponse{
		Status:             "sucesso",
		Message:            "Bem-vindo a Pokedex",
		SearchPlaceholder:  "Busque Pokemon por nome ou ID",
		RecommendedFilters: []string{"Fire", "Water", "Grass", "Electric", "Flying"},
		Data:               richList,
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
		"Normal":   "#A8A878",
		"Fire":     "#F08030",
		"Water":    "#6890F0",
		"Electric": "#F8D030",
		"Grass":    "#78C850",
		"Ice":      "#98D8D8",
		"Fighting": "#C03028",
		"Poison":   "#A040A0",
		"Ground":   "#E0C068",
		"Flying":   "#A890F0",
		"Psychic":  "#F85888",
		"Bug":      "#A8B820",
		"Rock":     "#B8A038",
		"Ghost":    "#705898",
		"Dragon":   "#7038F8",
		"Dark":     "#705848",
		"Steel":    "#B8B8D0",
		"Fairy":    "#EE99AC",
	}

	if color, exists := typeColors[typeStr]; exists {
		return color
	}
	return "#A9AC86"
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
