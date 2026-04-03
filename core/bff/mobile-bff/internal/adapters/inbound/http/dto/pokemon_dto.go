package dto

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

// HealthResponse representa uma resposta de verificação de saúde
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// PokemonDTO representa um Pokémon em respostas HTTP (básico)
type PokemonDTO struct {
	ID       string   `json:"id"`
	Number   string   `json:"number"`
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	ImageURL string   `json:"image_url"`
	Height   float64  `json:"height,omitempty"`
	Weight   float64  `json:"weight,omitempty"`
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

// TypeDTO representa um tipo de Pokémon com cor
type TypeDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
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

// SearchMetadata contém metadados relacionados à busca
type SearchMetadata struct {
	Placeholder string `json:"placeholder"`
}

// RichPokemonResponse representa Pokémon em formato rico para endpoints de home/lista
type RichPokemonResponse struct {
	Number     string    `json:"number"`
	Name       string    `json:"name"`
	Image      ImageDTO  `json:"image"`
	Types      []TypeDTO `json:"types"`
	IsFavorite bool      `json:"is_favorite,omitempty"`
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

// HomeSearchDTO representa a área de busca do home
type HomeSearchDTO struct {
	Placeholder string `json:"placeholder"`
	Value       string `json:"value,omitempty"`
}

// HomeFilterItemDTO representa um item de filtro
type HomeFilterItemDTO struct {
	Title string `json:"title"`
}

// HomeFilterGroupDTO representa um grupo de filtros
type HomeFilterGroupDTO struct {
	Title    string              `json:"title"`
	Selected string              `json:"selected,omitempty"`
	Items    []HomeFilterItemDTO `json:"items"`
}

// HomeFiltersDTO agrupa todos os filtros do home
type HomeFiltersDTO struct {
	Types    HomeFilterGroupDTO `json:"types"`
	Ordering HomeFilterGroupDTO `json:"ordering"`
	Region   HomeFilterGroupDTO `json:"region,omitempty"`
}

// HomePokemonTypeDTO representa o tipo de Pokémon para exibição no home
type HomePokemonTypeDTO struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

// HomePokemonSpritesDTO representa os sprites do Pokémon no home
type HomePokemonSpritesDTO struct {
	URL             string `json:"url"`
	BackgroundColor string `json:"backgroundColor"`
}

// HomePokemonDTO representa um Pokémon na listagem do home
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

// ScreenActionDTO representa uma ação de tela (botão, link, etc.)
type ScreenActionDTO struct {
	Label   string `json:"label"`
	Variant string `json:"variant,omitempty"`
}

// ScreenMessageDTO representa uma mensagem de tela com título e descrição
type ScreenMessageDTO struct {
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	CTA         *ScreenActionDTO `json:"cta,omitempty"`
}

// RegionItemDTO representa uma região de Pokémon
type RegionItemDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
}

// RegionsResponse representa a resposta com lista de regiões
type RegionsResponse struct {
	Title   string          `json:"title"`
	Regions []RegionItemDTO `json:"regions"`
}

// FavoritesResponse representa a resposta da tela de favoritos
type FavoritesResponse struct {
	Title    string            `json:"title"`
	State    string            `json:"state"`
	Message  *ScreenMessageDTO `json:"message,omitempty"`
	Pokemons []HomePokemonDTO  `json:"pokemons"`
}

// ProfileHeaderDTO representa o cabeçalho do perfil
type ProfileHeaderDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// ProfileUserDTO representa os dados do usuário no perfil
type ProfileUserDTO struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// ProfileSectionItemDTO representa um item de seção do perfil
type ProfileSectionItemDTO struct {
	Label       string `json:"label"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// ProfileSectionDTO representa uma seção do perfil
type ProfileSectionDTO struct {
	Title string                  `json:"title"`
	Items []ProfileSectionItemDTO `json:"items"`
}

// ProfileResponse representa a resposta da tela de perfil
type ProfileResponse struct {
	Title         string              `json:"title"`
	Authenticated bool                `json:"authenticated"`
	Header        *ProfileHeaderDTO   `json:"header,omitempty"`
	User          *ProfileUserDTO     `json:"user,omitempty"`
	Sections      []ProfileSectionDTO `json:"sections,omitempty"`
	Actions       []ScreenActionDTO   `json:"actions,omitempty"`
	Footer        *ScreenMessageDTO   `json:"footer,omitempty"`
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

// DetailInfoValueDTO representa um par label/valor nos detalhes
type DetailInfoValueDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// DetailAbilitiesDTO representa as habilidades nos detalhes
type DetailAbilitiesDTO struct {
	Label string   `json:"label"`
	Items []string `json:"items"`
}

// DetailGenderDTO representa o gênero nos detalhes
type DetailGenderDTO struct {
	Label  string `json:"label"`
	Male   string `json:"male,omitempty"`
	Female string `json:"female,omitempty"`
}

// DetailAboutDTO representa as informações gerais nos detalhes
type DetailAboutDTO struct {
	Weight    DetailInfoValueDTO `json:"weight"`
	Height    DetailInfoValueDTO `json:"height"`
	Category  DetailInfoValueDTO `json:"category"`
	Abilities DetailAbilitiesDTO `json:"abilities"`
	Gender    DetailGenderDTO    `json:"gender"`
}

// DetailEvolutionDTO representa uma evolução na tela de detalhes
type DetailEvolutionDTO struct {
	Number  string                `json:"number"`
	Name    string                `json:"name"`
	Types   []HomePokemonTypeDTO  `json:"types"`
	Sprites HomePokemonSpritesDTO `json:"sprites"`
	Trigger *DetailInfoValueDTO   `json:"trigger,omitempty"`
}

// PokemonDetailScreenResponse representa a tela de detalhes completa do Pokémon
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
