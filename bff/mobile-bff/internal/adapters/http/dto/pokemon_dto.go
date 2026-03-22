package dto

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// PokemonDTO represents a Pokémon in HTTP responses (basic)
type PokemonDTO struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Number   string   `json:"number"`
	Types    []string `json:"types"`
	ImageURL string   `json:"image_url"`
	Height   float64  `json:"height,omitempty"`
	Weight   float64  `json:"weight,omitempty"`
}

// TypeDTO represents a Pokémon type with color
type TypeDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ElementDTO represents the primary element/type with color
type ElementDTO struct {
	Color string `json:"color"`
	Type  string `json:"type"`
}

// ImageDTO represents image information with element
type ImageDTO struct {
	URL     string     `json:"url"`
	Element ElementDTO `json:"element"`
}

// PokemonDetailDTO represents enriched Pokémon details for rich responses
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

// PokemonListResponse represents a paginated list of Pokémons
type PokemonListResponse struct {
	Content       []PokemonDetailDTO `json:"content"`
	TotalElements int64              `json:"totalElements"`
	CurrentPage   int                `json:"currentPage"`
	TotalPages    int                `json:"totalPages"`
	HasNext       bool               `json:"hasNext"`
}

// RichPokemonResponse represents Pokémon in rich format for home/list endpoints
type RichPokemonResponse struct {
	Number     string    `json:"number"`
	Name       string    `json:"name"`
	Image      ImageDTO  `json:"image"`
	Types      []TypeDTO `json:"types"`
	IsFavorite bool      `json:"is_favorite,omitempty"`
}

// SearchMetadata holds search-related metadata
type SearchMetadata struct {
	Placeholder string `json:"placeholder"`
}

// RichPokemonListResponse represents enriched list response with search/filters
type RichPokemonListResponse struct {
	Content       []RichPokemonResponse `json:"content"`
	TotalElements int64                 `json:"totalElements"`
	CurrentPage   int                   `json:"currentPage"`
	TotalPages    int                   `json:"totalPages"`
	HasNext       bool                  `json:"hasNext"`
	Search        SearchMetadata        `json:"search"`
	Filters       []interface{}         `json:"filters"`
}

// HomeResponse represents data for the home screen
type HomeResponse struct {
	Status             string                   `json:"status"`
	Message            string                   `json:"message"`
	SearchPlaceholder  string                   `json:"searchPlaceholder"`
	RecommendedFilters []string                 `json:"recommendedFilters"`
	Data               *RichPokemonListResponse `json:"data"`
}

// FavoriteRequest represents a request to add a favorite
type FavoriteRequest struct {
	PokemonID string `json:"pokemon_id"`
}

// FavoriteResponse represents a response after favoriting
type FavoriteResponse struct {
	Message    string `json:"message"`
	PokemonID  string `json:"pokemon_id"`
	IsFavorite bool   `json:"is_favorite"`
}
