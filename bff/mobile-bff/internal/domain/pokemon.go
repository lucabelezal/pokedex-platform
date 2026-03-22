package domain

import (
	"errors"
	"time"
)

// Pokemon represents a Pokémon in the domain
type Pokemon struct {
	ID           string
	Name         string
	Number       string
	Types        []string
	Height       float64
	Weight       float64
	Description  string
	ImageURL     string
	ElementColor string
	ElementType  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Type represents a Pokémon type with color information
type Type struct {
	Name  string
	Color string
}

// PokemonDetail holds enriched Pokémon data for detailed views
type PokemonDetail struct {
	Number      string
	Name        string
	ImageURL    string
	Element     Element
	Types       []Type
	Height      float64
	Weight      float64
	Description string
	IsFavorite  bool
}

// Element represents the primary type element with color
type Element struct {
	Color string
	Type  string
}

// PokemonPage represents paginated Pokémon results
type PokemonPage struct {
	Content       []Pokemon
	TotalElements int64
	CurrentPage   int
	TotalPages    int
	HasNext       bool
}

// User represents a user for tracking favorites
type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Favorite represents a user's favorite Pokémon
type Favorite struct {
	ID        string
	UserID    string
	PokemonID string
	CreatedAt time.Time
}

// Validation methods
func (p *Pokemon) Validate() error {
	if p.Name == "" {
		return errors.New("nome do pokemon e obrigatorio")
	}
	if p.Number == "" {
		return errors.New("numero do pokemon e obrigatorio")
	}
	return nil
}

func (p *PokemonDetail) Validate() error {
	if p.Name == "" {
		return errors.New("nome do pokemon e obrigatorio")
	}
	if p.Number == "" {
		return errors.New("numero do pokemon e obrigatorio")
	}
	return nil
}
