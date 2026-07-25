package domain

// TypeColor retorna a cor hexadecimal (com prefixo #) associada a um tipo de Pokémon.
// Aceita nomes em português e inglês. Retorna "#A9AC86" (cor padrão) para tipos desconhecidos.
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
	"Fire":      "#F08030",
	"Água":      "#6390F0",
	"Water":     "#6890F0",
	"Elétrico":  "#F7D02C",
	"Electric":  "#F8D030",
	"Grama":     "#7AC74C",
	"Grass":     "#78C850",
	"Gelo":      "#96D9D6",
	"Ice":       "#98D8D8",
	"Lutador":   "#C22E28",
	"Fighting":  "#C03028",
	"Venenoso":  "#A33EA1",
	"Poison":    "#A040A0",
	"Terrestre": "#E2BF65",
	"Ground":    "#E0C068",
	"Voador":    "#A98FF3",
	"Flying":    "#A890F0",
	"Psíquico":  "#F95587",
	"Psychic":   "#F85888",
	"Inseto":    "#A6B91A",
	"Bug":       "#A8B820",
	"Pedra":     "#B6A136",
	"Rock":      "#B8A038",
	"Fantasma":  "#735797",
	"Ghost":     "#705898",
	"Dragão":    "#6F35FC",
	"Dragon":    "#7038F8",
	"Sombrio":   "#705746",
	"Noturno":   "#705746",
	"Dark":      "#705848",
	"Aço":       "#B7B7CE",
	"Metal":     "#B7B7CE",
	"Steel":     "#B8B8D0",
	"Fada":      "#D685AD",
	"Fairy":     "#EE99AC",
}
