package unit

import (
	"encoding/json"
	"testing"

	httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http"
	"pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http/dto"
	"pokedex-platform/core/bff/mobile-bff/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestResponseBuilderBuildPokemonDetailDTO(t *testing.T) {
	rb := httpadapter.NewResponseBuilder()

	detail := &domain.PokemonDetail{
		Number:      "025",
		Name:        "Pikachu",
		ImageURL:    "https://example.com/pikachu.png",
		Height:      0.41,
		Weight:      6.0,
		Description: "Electric mouse",
		Element: domain.Element{
			Color: "#F8D030",
			Type:  "Electric",
		},
		Types: []domain.Type{
			{Name: "Electric", Color: "#F8D030"},
		},
		IsFavorite: true,
	}

	result := rb.BuildPokemonDetailDTO(detail)

	assert.NotNil(t, result)
	assert.Equal(t, "025", result.Number)
	assert.Equal(t, "Pikachu", result.Name)
	assert.Equal(t, true, result.IsFavorite)
	assert.Equal(t, 1, len(result.Types))
}

func TestResponseBuilderBuildRichPokemonResponse(t *testing.T) {
	rb := httpadapter.NewResponseBuilder()

	pokemon := &domain.Pokemon{
		ID:           "25",
		Name:         "Pikachu",
		Number:       "025",
		Types:        []string{"Electric"},
		ImageURL:     "https://example.com/pikachu.png",
		ElementColor: "#F8D030",
		ElementType:  "Electric",
	}

	result := rb.BuildRichPokemonResponse(pokemon)

	assert.NotNil(t, result)
	assert.Equal(t, "025", result.Number)
	assert.Equal(t, "Pikachu", result.Name)
	assert.Equal(t, 1, len(result.Types))
	assert.Equal(t, "Electric", result.Types[0].Name)
}

func TestResponseBuilderBuildHealthResponse(t *testing.T) {
	rb := httpadapter.NewResponseBuilder()

	result := rb.BuildHealthResponse()

	assert.NotNil(t, result)
	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, "mobile-bff", result.Service)
}

func TestDTOJSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		dto  interface{}
	}{
		{
			name: "health response",
			dto: &dto.HealthResponse{
				Status:  "ok",
				Service: "mobile-bff",
			},
		},
		{
			name: "error response",
			dto: &dto.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: "pokemon not found",
				Code:    404,
			},
		},
		{
			name: "message response",
			dto: &dto.MessageResponse{
				Message: "Success",
			},
		},
		{
			name: "favorite response",
			dto: &dto.FavoriteResponse{
				Message:    "Added to favorites",
				PokemonID:  "25",
				IsFavorite: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.dto)
			assert.NoError(t, err)
			assert.NotEmpty(t, jsonBytes)

			var result interface{}
			err = json.Unmarshal(jsonBytes, &result)
			assert.NoError(t, err)
		})
	}
}

func TestPokemonListResponseMarshaling(t *testing.T) {
	response := &dto.PokemonListResponse{
		Content: []dto.PokemonDetailDTO{
			{
				Number: "025",
				Name:   "Pikachu",
				Image: dto.ImageDTO{
					URL: "https://example.com/pikachu.png",
					Element: dto.ElementDTO{
						Color: "#F8D030",
						Type:  "Electric",
					},
				},
				Types: []dto.TypeDTO{
					{Name: "Electric", Color: "#F8D030"},
				},
			},
		},
		TotalElements: 1,
		CurrentPage:   0,
		TotalPages:    1,
		HasNext:       false,
	}

	jsonBytes, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled dto.PokemonListResponse
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response.TotalElements, unmarshaled.TotalElements)
	assert.Len(t, unmarshaled.Content, 1)
}
