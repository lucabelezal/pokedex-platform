package dto

// HealthResponse representa uma resposta de verificação de saúde
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// ErrorResponse representa uma resposta de erro padrão
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// MessageResponse representa uma resposta simples com mensagem
type MessageResponse struct {
	Message string `json:"message"`
}

// PokemonDTO representa um Pokémon em respostas HTTP (básico)
type PokemonDTO struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Number   string   `json:"number"`
	Types    []string `json:"types"`
	ImageURL string   `json:"image_url"`
	Height   float64  `json:"height,omitempty"`
	Weight   float64  `json:"weight,omitempty"`
}

// TypeDTO representa um tipo de Pokémon com cor
type TypeDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ElementDTO representa o elemento/tipo primário com cor
type ElementDTO struct {
	Color string `json:"color"`
	Type  string `json:"type"`
}

// ImageDTO representa informações de imagem com elemento
type ImageDTO struct {
	URL     string     `json:"url"`
	Element ElementDTO `json:"element"`
}

// PokemonDetailDTO representa detalhes enriquecidos de Pokémon para respostas ricas
type PokemonDetailDTO struct {
	Number      string    `json:"number"`
	Name        string    `json:"name"`
	Image       ImageDTO  `json:"image"`
	Types       []TypeDTO `json:"types"`
	Height      float64   `json:"height,omitempty"`
	Weight      float64   `json:"weight,omitempty"`
	Description string    `json:"description,omitempty"`
	IsFavorite  bool      `json:"is_favorite"`
}

// PokemonListResponse representa uma lista paginada de Pokémons
type PokemonListResponse struct {
	Content       []PokemonDetailDTO `json:"content"`
	TotalElements int64              `json:"totalElements"`
	CurrentPage   int                `json:"currentPage"`
	TotalPages    int                `json:"totalPages"`
	HasNext       bool               `json:"hasNext"`
}

// RichPokemonResponse representa Pokémon em formato rico para endpoints de home/lista
type RichPokemonResponse struct {
	Number     string    `json:"number"`
	Name       string    `json:"name"`
	Image      ImageDTO  `json:"image"`
	Types      []TypeDTO `json:"types"`
	IsFavorite bool      `json:"is_favorite,omitempty"`
}

// SearchMetadata contém metadados relacionados à busca
type SearchMetadata struct {
	Placeholder string `json:"placeholder"`
}

// RichPokemonListResponse representa resposta de lista enriquecida com busca/filtros
type RichPokemonListResponse struct {
	Content       []RichPokemonResponse `json:"content"`
	TotalElements int64                 `json:"totalElements"`
	CurrentPage   int                   `json:"currentPage"`
	TotalPages    int                   `json:"totalPages"`
	HasNext       bool                  `json:"hasNext"`
	Search        SearchMetadata        `json:"search"`
	Filters       []interface{}         `json:"filters"`
}

type HomeSearchDTO struct {
	Placeholder string `json:"placeholder"`
	Value       string `json:"value,omitempty"`
}

type HomeFilterItemDTO struct {
	Title string `json:"title"`
}

type HomeFilterGroupDTO struct {
	Title    string              `json:"title"`
	Selected string              `json:"selected,omitempty"`
	Items    []HomeFilterItemDTO `json:"items"`
}

type HomeFiltersDTO struct {
	Types    HomeFilterGroupDTO `json:"types"`
	Ordering HomeFilterGroupDTO `json:"ordering"`
	Region   HomeFilterGroupDTO `json:"region,omitempty"`
}

type HomePokemonTypeDTO struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type HomePokemonSpritesDTO struct {
	URL             string `json:"url"`
	BackgroundColor string `json:"backgroundColor"`
}

type HomePokemonDTO struct {
	Number     string                `json:"number"`
	Name       string                `json:"name"`
	Types      []HomePokemonTypeDTO  `json:"types"`
	Sprites    HomePokemonSpritesDTO `json:"sprites"`
	IsFavorite bool                  `json:"isFavorite"`
}

// HomeResponse representa dados para a tela de pokedex/home
type HomeResponse struct {
	Title    string           `json:"title"`
	Search   HomeSearchDTO    `json:"search"`
	Filters  HomeFiltersDTO   `json:"filters"`
	Pokemons []HomePokemonDTO `json:"pokemons"`
}

type ScreenActionDTO struct {
	Label   string `json:"label"`
	Variant string `json:"variant,omitempty"`
}

type ScreenMessageDTO struct {
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	CTA         *ScreenActionDTO `json:"cta,omitempty"`
}

type RegionItemDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
}

type RegionsResponse struct {
	Title   string          `json:"title"`
	Regions []RegionItemDTO `json:"regions"`
}

type FavoritesResponse struct {
	Title    string            `json:"title"`
	State    string            `json:"state"`
	Message  *ScreenMessageDTO `json:"message,omitempty"`
	Pokemons []HomePokemonDTO  `json:"pokemons"`
}

type ProfileSectionItemDTO struct {
	Label       string `json:"label"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type ProfileSectionDTO struct {
	Title string                  `json:"title"`
	Items []ProfileSectionItemDTO `json:"items"`
}

type ProfileHeaderDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ProfileUserDTO struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type ProfileResponse struct {
	Title         string              `json:"title"`
	Authenticated bool                `json:"authenticated"`
	Header        *ProfileHeaderDTO   `json:"header,omitempty"`
	User          *ProfileUserDTO     `json:"user,omitempty"`
	Sections      []ProfileSectionDTO `json:"sections,omitempty"`
	Actions       []ScreenActionDTO   `json:"actions,omitempty"`
	Footer        *ScreenMessageDTO   `json:"footer,omitempty"`
}

type DetailInfoValueDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DetailAbilitiesDTO struct {
	Label string   `json:"label"`
	Items []string `json:"items"`
}

type DetailGenderDTO struct {
	Label  string `json:"label"`
	Male   string `json:"male,omitempty"`
	Female string `json:"female,omitempty"`
}

type DetailAboutDTO struct {
	Weight    DetailInfoValueDTO `json:"weight"`
	Height    DetailInfoValueDTO `json:"height"`
	Category  DetailInfoValueDTO `json:"category"`
	Abilities DetailAbilitiesDTO `json:"abilities"`
	Gender    DetailGenderDTO    `json:"gender"`
}

type DetailEvolutionDTO struct {
	Number  string                `json:"number"`
	Name    string                `json:"name"`
	Types   []HomePokemonTypeDTO  `json:"types"`
	Sprites HomePokemonSpritesDTO `json:"sprites"`
	Trigger *DetailInfoValueDTO   `json:"trigger,omitempty"`
}

type PokemonDetailScreenResponse struct {
	Number      string                `json:"number"`
	Name        string                `json:"name"`
	Types       []HomePokemonTypeDTO  `json:"types"`
	Description string                `json:"description"`
	Sprites     HomePokemonSpritesDTO `json:"sprites"`
	About       DetailAboutDTO        `json:"about"`
	Weaknesses  []HomePokemonTypeDTO  `json:"weaknesses"`
	Evolutions  []DetailEvolutionDTO  `json:"evolutions"`
	IsFavorite  bool                  `json:"isFavorite"`
}

// FavoriteRequest representa uma requisição para adicionar um favorito
type FavoriteRequest struct {
	PokemonID string `json:"pokemon_id"`
}

// FavoriteResponse representa uma resposta após favoritar
type FavoriteResponse struct {
	Message    string `json:"message"`
	PokemonID  string `json:"pokemon_id"`
	IsFavorite bool   `json:"is_favorite"`
}
