package domain

import (
	"errors"
	"time"
)

// Pokemon representa um Pokémon no domínio
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

// Type representa um tipo de Pokémon com informações de cor
type Type struct {
	Name  string
	Color string
}

// PokemonDetail contém dados enriquecidos de Pokémon para visualização detalhada
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

// Element representa o tipo primário com cor
type Element struct {
	Color string
	Type  string
}

// PokemonPage representa resultados paginados de Pokémons
type PokemonPage struct {
	Content       []Pokemon
	TotalElements int64
	CurrentPage   int
	TotalPages    int
	HasNext       bool
}

// User representa um usuário para rastrear favoritos
type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Favorite representa um Pokémon favorito do usuário
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
