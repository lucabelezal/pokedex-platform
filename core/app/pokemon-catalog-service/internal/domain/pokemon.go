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

// Type representa um tipo canônico do catálogo com sua cor visual.
type Type struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Region representa uma regiao e sua geracao para a UI.
type Region struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
}

// Evolution representa um item da cadeia evolutiva.
type Evolution struct {
	ID       string `json:"id"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
	Types    []Type `json:"types"`
	Trigger  string `json:"trigger,omitempty"`
}

// PokemonDetail representa os dados ricos de um pokemon no catalogo.
type PokemonDetail struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Number       string      `json:"number"`
	Types        []Type      `json:"types"`
	Description  string      `json:"description"`
	ImageURL     string      `json:"imageUrl"`
	ElementColor string      `json:"elementColor"`
	Height       float64     `json:"height"`
	Weight       float64     `json:"weight"`
	Category     string      `json:"category"`
	Abilities    []string    `json:"abilities"`
	GenderMale   *float64    `json:"genderMale,omitempty"`
	GenderFemale *float64    `json:"genderFemale,omitempty"`
	Weaknesses   []Type      `json:"weaknesses"`
	Evolutions   []Evolution `json:"evolutions"`
	Region       string      `json:"region,omitempty"`
	Generation   string      `json:"generation,omitempty"`
}
