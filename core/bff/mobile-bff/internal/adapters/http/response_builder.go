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
	return rb.BuildHomePageResponseWithTypes(page, nil, nil, nil, "", "Todos os tipos", "Menor número", "")
}

func (rb *ResponseBuilder) BuildHomePageResponseWithTypes(
	page *domain.PokemonPage,
	types []domain.Type,
	regions []domain.Region,
	favoriteSet map[string]struct{},
	searchValue string,
	selectedType string,
	selectedOrdering string,
	selectedRegion string,
) *dto.HomeResponse {
	pokemons := make([]dto.HomePokemonDTO, len(page.Content))
	for i, p := range page.Content {
		pokemons[i] = rb.BuildHomePokemonDTO(&p, favoriteSet)
	}

	typeItems := make([]dto.HomeFilterItemDTO, 0, len(types)+1)
	typeItems = append(typeItems, dto.HomeFilterItemDTO{Title: "Todos os tipos"})
	for _, t := range types {
		typeItems = append(typeItems, dto.HomeFilterItemDTO{Title: displayTypeTitle(t.Name)})
	}

	regionItems := make([]dto.HomeFilterItemDTO, len(regions))
	for i, region := range regions {
		regionItems[i] = dto.HomeFilterItemDTO{Title: region.Name}
	}

	if strings.TrimSpace(selectedType) == "" {
		selectedType = "Todos os tipos"
	}
	if strings.TrimSpace(selectedOrdering) == "" {
		selectedOrdering = "Menor número"
	}

	return &dto.HomeResponse{
		Title: "Pokédex",
		Search: dto.HomeSearchDTO{
			Placeholder: "Procurar Pokémon...",
			Value:       searchValue,
		},
		Filters: dto.HomeFiltersDTO{
			Types: dto.HomeFilterGroupDTO{
				Title:    "Tipos",
				Selected: selectedType,
				Items:    typeItems,
			},
			Ordering: dto.HomeFilterGroupDTO{
				Title:    "Ordenação",
				Selected: selectedOrdering,
				Items: []dto.HomeFilterItemDTO{
					{Title: "Menor número"},
					{Title: "Maior número"},
					{Title: "A-Z"},
					{Title: "Z-A"},
				},
			},
			Region: dto.HomeFilterGroupDTO{
				Title:    "Regiões",
				Selected: selectedRegion,
				Items:    regionItems,
			},
		},
		Pokemons: pokemons,
	}
}

func (rb *ResponseBuilder) BuildHomePokemonDTO(
	p *domain.Pokemon,
	favoriteSet map[string]struct{},
) dto.HomePokemonDTO {
	types := buildHomePokemonTypes(p.Number, p.Types)
	backgroundColor := normalizeHexColor(p.ElementColor)
	if len(types) > 0 {
		backgroundColor = types[0].Color
	}

	_, isFavorite := favoriteSet[normalizePokemonID(p.Number)]

	return dto.HomePokemonDTO{
		Number: formatHomePokemonNumber(p.Number),
		Name:   p.Name,
		Types:  types,
		Sprites: dto.HomePokemonSpritesDTO{
			URL:             p.ImageURL,
			BackgroundColor: backgroundColor,
		},
		IsFavorite: isFavorite,
	}
}

func (rb *ResponseBuilder) BuildRegionsResponse(regions []domain.Region) *dto.RegionsResponse {
	items := make([]dto.RegionItemDTO, len(regions))
	for i, region := range regions {
		items[i] = dto.RegionItemDTO{
			ID:         region.ID,
			Name:       region.Name,
			Generation: region.Generation,
		}
	}

	return &dto.RegionsResponse{
		Title:   "Regiões",
		Regions: items,
	}
}

func (rb *ResponseBuilder) BuildFavoritesResponse(
	page *domain.PokemonPage,
	favoriteSet map[string]struct{},
	authenticated bool,
) *dto.FavoritesResponse {
	if !authenticated {
		return &dto.FavoritesResponse{
			Title: "Favoritos",
			State: "unauthenticated",
			Message: &dto.ScreenMessageDTO{
				Title:       "Você precisa entrar em uma conta para fazer isso.",
				Description: "Para acessar essa funcionalidade, é necessário fazer login ou criar uma conta. Faça isso agora!",
				CTA:         &dto.ScreenActionDTO{Label: "Entre ou Cadastre-se", Variant: "primary"},
			},
			Pokemons: []dto.HomePokemonDTO{},
		}
	}

	pokemons := make([]dto.HomePokemonDTO, len(page.Content))
	for i, p := range page.Content {
		pokemons[i] = rb.BuildHomePokemonDTO(&p, favoriteSet)
		pokemons[i].IsFavorite = true
	}

	if len(pokemons) == 0 {
		return &dto.FavoritesResponse{
			Title: "Favoritos",
			State: "empty",
			Message: &dto.ScreenMessageDTO{
				Title:       "Você não favoritou nenhum Pokémon :(",
				Description: "Clique no ícone de coração dos seus pokémons favoritos e eles aparecerão aqui.",
			},
			Pokemons: []dto.HomePokemonDTO{},
		}
	}

	return &dto.FavoritesResponse{
		Title:    "Favoritos",
		State:    "has_data",
		Pokemons: pokemons,
	}
}

func (rb *ResponseBuilder) BuildProfileResponse(authenticated bool, email string) *dto.ProfileResponse {
	if !authenticated {
		return &dto.ProfileResponse{
			Title:         "Conta",
			Authenticated: false,
			Header: &dto.ProfileHeaderDTO{
				Title:       "Entre ou Cadastre-se",
				Description: "Mantenha sua Pokédex atualizada e participe desse mundo.",
			},
			Actions: []dto.ScreenActionDTO{
				{Label: "Entre ou Cadastre-se", Variant: "primary"},
			},
		}
	}

	displayName := buildDisplayName(email)

	return &dto.ProfileResponse{
		Title:         "Conta",
		Authenticated: true,
		User: &dto.ProfileUserDTO{
			Name:  displayName,
			Email: email,
		},
		Sections: []dto.ProfileSectionDTO{
			{
				Title: "Informações da conta",
				Items: []dto.ProfileSectionItemDTO{
					{Label: "Nome", Value: displayName},
					{Label: "Email", Value: email},
					{Label: "Senha", Value: "••••••••••••••••"},
				},
			},
			{
				Title: "Pokédex",
				Items: []dto.ProfileSectionItemDTO{
					{Label: "Mega evoluções", Description: "Habilita a exibição de mega evoluções.", Type: "toggle", Value: "false"},
					{Label: "Outras formas", Description: "Habilita a exibição de formas alternativas de pokémon.", Type: "toggle", Value: "false"},
				},
			},
			{
				Title: "Notificações",
				Items: []dto.ProfileSectionItemDTO{
					{Label: "Atualizações na pokédex", Description: "Novos Pokémons, habilidades, informações, etc.", Type: "toggle", Value: "false"},
					{Label: "Mundo Pokémon", Description: "Acontecimentos e informações do mundo de pokémon.", Type: "toggle", Value: "false"},
				},
			},
			{
				Title: "Idioma",
				Items: []dto.ProfileSectionItemDTO{
					{Label: "Idioma da interface", Value: "Português (PT-BR)"},
					{Label: "Idioma de informações em jogo", Value: "English (US)"},
				},
			},
			{
				Title: "Geral",
				Items: []dto.ProfileSectionItemDTO{
					{Label: "Versão", Value: "0.8.12"},
					{Label: "Termos e condições", Description: "Tudo o que você precisa saber."},
					{Label: "Central de ajuda", Description: "Precisa de ajuda? Fale conosco."},
					{Label: "Sobre", Description: "Saiba mais sobre o app."},
				},
			},
		},
		Actions: []dto.ScreenActionDTO{{Label: "Sair", Variant: "danger"}},
		Footer:  &dto.ScreenMessageDTO{Title: fmt.Sprintf("Você entrou como %s.", displayName)},
	}
}

func (rb *ResponseBuilder) BuildPokemonDetailScreenResponse(
	detail *domain.PokemonScreenDetail,
	isFavorite bool,
) *dto.PokemonDetailScreenResponse {
	types := buildScreenPokemonTypes(detail.Number, detail.Types)
	weaknesses := buildScreenPokemonTypes(detail.Number, detail.Weaknesses)

	evolutions := make([]dto.DetailEvolutionDTO, len(detail.Evolutions))
	for i, item := range detail.Evolutions {
		evolutionTypes := buildScreenPokemonTypes(item.Number, item.Types)
		var trigger *dto.DetailInfoValueDTO
		if strings.TrimSpace(item.Trigger) != "" {
			trigger = &dto.DetailInfoValueDTO{Label: item.Trigger, Value: item.Trigger}
		}
		evolutions[i] = dto.DetailEvolutionDTO{
			Number: formatHomePokemonNumber(item.Number),
			Name:   item.Name,
			Types:  evolutionTypes,
			Sprites: dto.HomePokemonSpritesDTO{
				URL: item.ImageURL,
			},
			Trigger: trigger,
		}
	}

	return &dto.PokemonDetailScreenResponse{
		Number:      formatHomePokemonNumber(detail.Number),
		Name:        detail.Name,
		Types:       types,
		Description: detail.Description,
		Sprites: dto.HomePokemonSpritesDTO{
			URL:             detail.ImageURL,
			BackgroundColor: detailBackgroundColor(detail.ElementColor, types),
		},
		About: dto.DetailAboutDTO{
			Weight:    dto.DetailInfoValueDTO{Label: "Peso", Value: formatWeight(detail.Weight)},
			Height:    dto.DetailInfoValueDTO{Label: "Altura", Value: formatHeight(detail.Height)},
			Category:  dto.DetailInfoValueDTO{Label: "Categoria", Value: detail.Category},
			Abilities: dto.DetailAbilitiesDTO{Label: "Habilidade", Items: detail.Abilities},
			Gender: dto.DetailGenderDTO{
				Label:  "Gênero",
				Male:   formatGender(detail.GenderMale),
				Female: formatGender(detail.GenderFemale),
			},
		},
		Weaknesses: weaknesses,
		Evolutions: evolutions,
		IsFavorite: isFavorite,
	}
}

func detailBackgroundColor(fallback string, types []dto.HomePokemonTypeDTO) string {
	if len(types) > 0 {
		return types[0].Color
	}
	return normalizeHexColor(fallback)
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
		"Noturno":   "#705746",
		"Dark":      "#705848",
		"Aço":       "#B7B7CE",
		"Metal":     "#B7B7CE",
		"Steel":     "#B8B8D0",
		"Fada":      "#D685AD",
		"Fairy":     "#EE99AC",
	}

	if color, exists := typeColors[typeStr]; exists {
		return color
	}
	return "#A9AC86"
}

func displayTypeTitle(value string) string {
	switch strings.TrimSpace(value) {
	case "Aço":
		return "Metal"
	case "Sombrio":
		return "Noturno"
	default:
		return strings.TrimSpace(value)
	}
}

func buildHomePokemonTypes(number string, rawTypes []string) []dto.HomePokemonTypeDTO {
	items := make([]dto.HomePokemonTypeDTO, len(rawTypes))
	for i, item := range rawTypes {
		title := displayTypeTitle(item)
		items[i] = dto.HomePokemonTypeDTO{
			Title: title,
			Color: normalizeHexColor(getTypeColor(title)),
		}
	}
	return reorderTypeDTOs(number, items)
}

func buildScreenPokemonTypes(number string, rawTypes []domain.Type) []dto.HomePokemonTypeDTO {
	items := make([]dto.HomePokemonTypeDTO, len(rawTypes))
	for i, item := range rawTypes {
		items[i] = dto.HomePokemonTypeDTO{
			Title: displayTypeTitle(item.Name),
			Color: normalizeHexColor(item.Color),
		}
	}
	return reorderTypeDTOs(number, items)
}

func reorderTypeDTOs(number string, items []dto.HomePokemonTypeDTO) []dto.HomePokemonTypeDTO {
	order, ok := pokemonTypeOrderOverrides[normalizePokemonID(number)]
	if !ok || len(items) <= 1 {
		return items
	}

	indexed := make(map[string]dto.HomePokemonTypeDTO, len(items))
	for _, item := range items {
		indexed[item.Title] = item
	}

	reordered := make([]dto.HomePokemonTypeDTO, 0, len(items))
	used := make(map[string]struct{}, len(items))
	for _, title := range order {
		if item, exists := indexed[title]; exists {
			reordered = append(reordered, item)
			used[title] = struct{}{}
		}
	}

	for _, item := range items {
		if _, exists := used[item.Title]; !exists {
			reordered = append(reordered, item)
		}
	}

	return reordered
}

var pokemonTypeOrderOverrides = map[string][]string{
	"15":  {"Inseto", "Venenoso"},
	"95":  {"Pedra", "Terrestre"},
	"306": {"Metal", "Terrestre"},
	"448": {"Lutador", "Metal"},
	"609": {"Fantasma", "Fogo"},
	"733": {"Voador", "Normal"},
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

func formatWeight(value float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.1f kg", value), ".", ",")
}

func formatHeight(value float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.1f m", value), ".", ",")
}

func formatGender(value *float64) string {
	if value == nil {
		return "Desconhecido"
	}
	return strings.ReplaceAll(fmt.Sprintf("%.1f%%", *value), ".", ",")
}

func buildDisplayName(email string) string {
	local := strings.TrimSpace(strings.Split(email, "@")[0])
	if local == "" {
		return "Treinador"
	}
	return local
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
