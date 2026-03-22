package unit

import (
	"testing"

	"pokedex-platform/bff/mobile-bff/internal/domain"

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
