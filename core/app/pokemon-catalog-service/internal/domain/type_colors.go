package domain

// TypeColor retorna a cor hexadecimal (com prefixo #) associada a um tipo de Pokémon.
// Aceita nomes em português. Retorna "#A9AC86" (cor padrão) para tipos desconhecidos.
func TypeColor(name string) string {
	if color, ok := typeColors[name]; ok {
		return color
	}
	return defaultTypeColor
}

func AllTypeColors() map[string]string {
	result := make(map[string]string, len(typeColors))
	for k, v := range typeColors {
		result[k] = v
	}
	return result
}

const defaultTypeColor = "#A9AC86"

var typeColors = map[string]string{
	"Normal":    "#A8A878",
	"Fogo":      "#EE8130",
	"Água":      "#6390F0",
	"Elétrico":  "#F7D02C",
	"Grama":     "#7AC74C",
	"Gelo":      "#96D9D6",
	"Lutador":   "#C22E28",
	"Venenoso":  "#A33EA1",
	"Terrestre": "#E2BF65",
	"Voador":    "#A98FF3",
	"Psíquico":  "#F95587",
	"Inseto":    "#A6B91A",
	"Pedra":     "#B6A136",
	"Fantasma":  "#735797",
	"Dragão":    "#6F35FC",
	"Noturno":   "#705746",
	"Sombrio":   "#705746",
	"Metal":     "#B7B7CE",
	"Aço":       "#B7B7CE",
	"Fada":      "#D685AD",
}
