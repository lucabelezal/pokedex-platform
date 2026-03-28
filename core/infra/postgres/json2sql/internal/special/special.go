// Package special handles the three tables that need custom SQL generation:
//   - evolution_chains  – maps JSON "chain" field to "chain_data" JSONB column
//   - pokemons          – also emits rows into pokemon_types, pokemon_abilities,
//     pokemon_egg_groups
//   - pokemon_weaknesses – resolves type name → type_id via the types lookup map
package special

import (
	"fmt"

	"pokedex-platform/core/infra/postgres/json2sql/internal/mapper"
	"pokedex-platform/core/infra/postgres/json2sql/internal/writer"
)

// TypeNameToID maps canonical Portuguese type names to their DB IDs.
// Built dynamically from 02_type.json by BuildTypeMap; populated here as
// fallback so the package compiles standalone in tests.
var TypeNameToID = map[string]int{
	"Normal":    1,
	"Fogo":      2,
	"Água":      3,
	"Elétrico":  4,
	"Grama":     5,
	"Gelo":      6,
	"Lutador":   7,
	"Venenoso":  8,
	"Terrestre": 9,
	"Voador":    10,
	"Psíquico":  11,
	"Inseto":    12,
	"Pedra":     13,
	"Fantasma":  14,
	"Dragão":    15,
	"Sombrio":   16,
	"Aço":       17,
	"Fada":      18,
}

// BuildTypeMap derives TypeNameToID from the parsed type records so that the
// map always reflects whatever is in 02_type.json.
func BuildTypeMap(types []map[string]any) {
	for _, t := range types {
		name, ok1 := t["name"].(string)
		idRaw, ok2 := t["id"]
		if !ok1 || !ok2 {
			continue
		}
		id := toInt(idRaw)
		if id > 0 {
			TypeNameToID[name] = id
		}
	}
}

// EvolutionChains generates INSERT statements for the evolution_chains table.
func EvolutionChains(records []map[string]any) []string {
	cols := mapper.TableValidFields["evolution_chains"] // ["id","chain_data"]
	var stmts []string
	for _, r := range records {
		row := map[string]any{
			"id":         r["id"],
			"chain_data": r["chain"], // rename: "chain" → "chain_data"
		}
		if s := writer.InsertRow("evolution_chains", row, cols); s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

// Pokemons generates INSERT statements for pokemons and its 3 join tables.
func Pokemons(records []map[string]any) []string {
	mainCols := mapper.TableValidFields["pokemons"]
	var stmts []string

	for _, r := range records {
		// Flatten gender sub-object into main record.
		row := make(map[string]any, len(r))
		for k, v := range r {
			row[k] = v
		}
		if gObj, ok := r["gender"].(map[string]any); ok {
			row["gender_male"] = gObj["male"]
			row["gender_female"] = gObj["female"]
		}

		filtered := mapper.FilterFields(row, "pokemons")
		if s := writer.InsertRow("pokemons", filtered, mainCols); s != "" {
			stmts = append(stmts, s)
		}

		pokemonID := toInt(r["id"])

		// pokemon_types
		if typeIDs, ok := r["type_ids"].([]any); ok {
			for _, tid := range typeIDs {
				stmts = append(stmts,
					fmt.Sprintf("INSERT INTO pokemon_types (pokemon_id, type_id) VALUES (%d, %d);",
						pokemonID, toInt(tid)))
			}
		}

		// pokemon_abilities
		if abilities, ok := r["abilities"].([]any); ok {
			for _, a := range abilities {
				ab, ok := a.(map[string]any)
				if !ok {
					continue
				}
				abilityID := toInt(ab["ability_id"])
				isHidden := ab["is_hidden"]
				hiddenVal := "FALSE"
				if b, ok := isHidden.(bool); ok && b {
					hiddenVal = "TRUE"
				}
				stmts = append(stmts,
					fmt.Sprintf("INSERT INTO pokemon_abilities (pokemon_id, ability_id, is_hidden) VALUES (%d, %d, %s);",
						pokemonID, abilityID, hiddenVal))
			}
		}

		// pokemon_egg_groups
		if eggIDs, ok := r["egg_group_ids"].([]any); ok {
			for _, eid := range eggIDs {
				stmts = append(stmts,
					fmt.Sprintf("INSERT INTO pokemon_egg_groups (pokemon_id, egg_group_id) VALUES (%d, %d);",
						pokemonID, toInt(eid)))
			}
		}
	}
	return stmts
}

// PokemonWeaknesses generates INSERT statements for pokemon_weaknesses
// by resolving type names to IDs using TypeNameToID.
func PokemonWeaknesses(records []map[string]any) ([]string, []string) {
	var stmts []string
	var warns []string

	for _, r := range records {
		pokemonID := toInt(r["pokemon_id"])
		weaknesses, ok := r["weaknesses"].([]any)
		if !ok {
			continue
		}
		for _, w := range weaknesses {
			name, ok := w.(string)
			if !ok {
				continue
			}
			typeID, found := TypeNameToID[name]
			if !found {
				warns = append(warns, fmt.Sprintf(
					"unknown weakness type %q for pokemon_id %d – skipped", name, pokemonID))
				continue
			}
			stmts = append(stmts,
				fmt.Sprintf("INSERT INTO pokemon_weaknesses (pokemon_id, type_id) VALUES (%d, %d);",
					pokemonID, typeID))
		}
	}
	return stmts, warns
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}
