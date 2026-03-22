package domain

import "time"

// Pokemon representa o registro canonico de um Pokemon no catalogo.
type Pokemon struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Number       string    `json:"number"`
	Types        []string  `json:"types"`
	Height       float64   `json:"height"`
	Weight       float64   `json:"weight"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"imageUrl"`
	ElementColor string    `json:"elementColor"`
	ElementType  string    `json:"elementType"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

// PokemonPage representa uma pagina de resultados do catalogo.
type PokemonPage struct {
	Content       []Pokemon `json:"content"`
	TotalElements int64     `json:"totalElements"`
	CurrentPage   int       `json:"currentPage"`
	TotalPages    int       `json:"totalPages"`
	HasNext       bool      `json:"hasNext"`
}
