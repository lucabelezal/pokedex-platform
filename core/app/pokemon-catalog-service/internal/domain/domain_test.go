package domain

import (
	"testing"
	"time"
)

func TestTypeColor(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"fogo", "Fogo", "#EE8130"},
		{"agua", "Água", "#6390F0"},
		{"metal via aco", "Aço", "#B7B7CE"},
		{"desconhecido usa default", "Desconhecido", defaultTypeColor},
		{"vazio usa default", "", defaultTypeColor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeColor(tt.in); got != tt.want {
				t.Errorf("TypeColor(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAllTypeColors(t *testing.T) {
	colors := AllTypeColors()
	if len(colors) == 0 {
		t.Fatal("AllTypeColors() retornou mapa vazio")
	}
	// cópia não deve mutar o mapa original
	colors["Fogo"] = "#000000"
	if typeColors["Fogo"] == "#000000" {
		t.Error("AllTypeColors retornou referência ao mapa interno")
	}
}

func TestPokemonJSONTags(t *testing.T) {
	p := Pokemon{
		ID:     "1",
		Name:   "Bulbasaur",
		Number: "001",
	}
	got := map[string]any{
		"id":     p.ID,
		"name":   p.Name,
		"number": p.Number,
	}
	if got["id"] != "1" || got["name"] != "Bulbasaur" {
		t.Error("campos básicos do Pokemon inválidos")
	}
}

func TestPokemonPageFields(t *testing.T) {
	now := time.Now()
	page := PokemonPage{
		Content:       []Pokemon{{ID: "1", CreatedAt: now}},
		TotalElements: 1,
		CurrentPage:   0,
		TotalPages:    1,
		HasNext:       false,
	}
	if page.TotalElements != 1 || page.CurrentPage != 0 || page.TotalPages != 1 {
		t.Error("campos do PokemonPage inválidos")
	}
	if page.HasNext {
		t.Error("HasNext deveria ser false com 1 página")
	}
}

func TestPokemonDetailPointers(t *testing.T) {
	male := 88.0
	female := 12.0
	d := PokemonDetail{
		GenderMale:   &male,
		GenderFemale: &female,
	}
	if d.GenderMale == nil || *d.GenderMale != 88.0 {
		t.Error("GenderMale deveria apontar para 88.0")
	}
	if d.GenderFemale == nil || *d.GenderFemale != 12.0 {
		t.Error("GenderFemale deveria apontar para 12.0")
	}
}

func TestTypeAndRegionAndEvolution(t *testing.T) {
	tt := Type{Name: "Fogo", Color: "#EE8130"}
	if tt.Name == "" || tt.Color == "" {
		t.Error("Type com campos vazios")
	}

	r := Region{ID: "kanto", Name: "Kanto", Generation: "1º Geração"}
	if r.ID == "" || r.Name == "" {
		t.Error("Region com campos vazios")
	}

	e := Evolution{ID: "1", Number: "001", Name: "Bulbasaur", Types: []Type{tt}}
	if e.Number != "001" || len(e.Types) != 1 {
		t.Error("Evolution com campos inválidos")
	}
}
