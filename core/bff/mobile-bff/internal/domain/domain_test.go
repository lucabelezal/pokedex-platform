package domain_test

import (
	"testing"

	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestPokemonValidation(t *testing.T) {
	tests := []struct {
		name    string
		pokemon *domain.Pokemon
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid pokemon",
			pokemon: &domain.Pokemon{
				Name:   "Pikachu",
				Number: "025",
				Types:  []string{"Electric"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			pokemon: &domain.Pokemon{
				Name:   "",
				Number: "025",
			},
			wantErr: true,
			errMsg:  "nome do pokemon e obrigatorio",
		},
		{
			name: "missing number",
			pokemon: &domain.Pokemon{
				Name:   "Pikachu",
				Number: "",
			},
			wantErr: true,
			errMsg:  "numero do pokemon e obrigatorio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pokemon.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPokemonDetailValidation(t *testing.T) {
	tests := []struct {
		name    string
		detail  *domain.PokemonDetail
		wantErr bool
	}{
		{
			name: "valid detail",
			detail: &domain.PokemonDetail{
				Name:   "Pikachu",
				Number: "025",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			detail: &domain.PokemonDetail{
				Name:   "",
				Number: "025",
			},
			wantErr: true,
		},
		{
			name: "missing number",
			detail: &domain.PokemonDetail{
				Name:   "Pikachu",
				Number: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.detail.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTypeColor(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     string
	}{
		{name: "tipo normal PT", typeName: "Normal", want: "#A8A878"},
		{name: "tipo fogo PT", typeName: "Fogo", want: "#EE8130"},
		{name: "tipo água PT", typeName: "Água", want: "#6390F0"},
		{name: "tipo elétrico PT", typeName: "Elétrico", want: "#F7D02C"},
		{name: "tipo grama PT", typeName: "Grama", want: "#7AC74C"},
		{name: "tipo dragão PT", typeName: "Dragão", want: "#6F35FC"},
		{name: "tipo fada PT", typeName: "Fada", want: "#D685AD"},
		{name: "tipo aço PT", typeName: "Aço", want: "#B7B7CE"},
		{name: "tipo noturno PT", typeName: "Noturno", want: "#705746"},
		{name: "tipo sombrio PT", typeName: "Sombrio", want: "#705746"},
		{name: "tipo metal PT", typeName: "Metal", want: "#B7B7CE"},
		{name: "tipo fire EN", typeName: "Fire", want: "#F08030"},
		{name: "tipo water EN", typeName: "Water", want: "#6890F0"},
		{name: "tipo grass EN", typeName: "Grass", want: "#78C850"},
		{name: "tipo dragon EN", typeName: "Dragon", want: "#7038F8"},
		{name: "tipo fairy EN", typeName: "Fairy", want: "#EE99AC"},
		{name: "tipo desconhecido", typeName: "Desconhecido", want: "#A9AC86"},
		{name: "string vazia", typeName: "", want: "#A9AC86"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.TypeColor(tt.typeName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllTypeColors(t *testing.T) {
	colors := domain.AllTypeColors()
	assert.NotEmpty(t, colors)
	assert.Greater(t, len(colors), 30, "deve ter mais de 30 entradas (PT + EN + aliases)")
}
