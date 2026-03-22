package mapper

// FileOrder is the canonical processing order of source JSON files.
var FileOrder = []string{
	"01_region.json",
	"02_type.json",
	"03_egg_group.json",
	"04_generation.json",
	"05_ability.json",
	"06_species.json",
	"07_stats.json",
	"08_evolution_chains.json",
	"09_pokemon.json",
	"10_weaknesses.json",
}

// FileToTable maps each source file to its primary target table.
var FileToTable = map[string]string{
	"01_region.json":         "regions",
	"02_type.json":           "types",
	"03_egg_group.json":      "egg_groups",
	"04_generation.json":     "generations",
	"05_ability.json":        "abilities",
	"06_species.json":        "species",
	"07_stats.json":          "stats",
	"08_evolution_chains.json": "evolution_chains",
	"09_pokemon.json":        "pokemons",
	"10_weaknesses.json":     "pokemon_weaknesses",
}

// TableValidFields lists the columns accepted by each table.
// Fields present in the JSON but absent here are silently dropped.
var TableValidFields = map[string][]string{
	"regions":    {"id", "name"},
	"types":      {"id", "name", "color"},
	"egg_groups": {"id", "name"},
	"generations": {"id", "name", "region_id"},
	"abilities":  {"id", "name", "description", "introduced_generation_id"},
	"species": {
		"id", "pokemon_number", "name", "species_en", "species_pt",
		"description", "color", "generation_id",
	},
	"stats": {
		"id", "total", "hp", "attack", "defense",
		"sp_atk", "sp_def", "speed",
	},
	"evolution_chains": {"id", "chain_data"},
	"pokemons": {
		"id", "number", "name", "height", "weight", "description", "sprites",
		"gender_male", "gender_female", "gender_rate_value", "egg_cycles",
		"stats_id", "generation_id", "species_id", "region_id", "evolution_chain_id",
	},
	"pokemon_weaknesses": {"pokemon_id", "type_id"},
}

// SpecialTables is the set of tables that need custom generation logic.
var SpecialTables = map[string]bool{
	"evolution_chains":   true,
	"pokemons":           true,
	"pokemon_weaknesses": true,
}

// FilterFields returns only the fields from record that are valid for table.
func FilterFields(record map[string]any, table string) map[string]any {
	valid, ok := TableValidFields[table]
	if !ok {
		return record
	}
	out := make(map[string]any, len(valid))
	for _, f := range valid {
		if v, exists := record[f]; exists {
			out[f] = v
		}
	}
	return out
}
