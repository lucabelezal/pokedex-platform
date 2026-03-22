// Package preflight performs cross-file referential integrity checks before
// SQL generation begins. Warnings are always printed; with --strict the tool
// exits non-zero on the first violation.
package preflight

import (
	"fmt"
)

// Result holds all warnings found during preflight.
type Result struct {
	Warnings []string
}

// HasIssues returns true when at least one warning was recorded.
func (r *Result) HasIssues() bool { return len(r.Warnings) > 0 }

func (r *Result) warn(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

// Run validates referential integrity across all parsed datasets.
// datasets is a map of tableName → records (already loaded from JSON).
func Run(datasets map[string][]map[string]any) *Result {
	res := &Result{}

	// Collect valid IDs per table.
	ids := func(table string) map[int]bool {
		m := make(map[int]bool)
		for _, r := range datasets[table] {
			if id := toInt(r["id"]); id > 0 {
				m[id] = true
			}
		}
		return m
	}

	regionIDs := ids("regions")
	typeIDs := ids("types")
	eggGroupIDs := ids("egg_groups")
	generationIDs := ids("generations")
	abilityIDs := ids("abilities")
	speciesIDs := ids("species")
	statsIDs := ids("stats")
	pokemonIDs := ids("pokemons")
	evolutionChainIDs := ids("evolution_chains")

	// generations.region_id → regions
	for _, r := range datasets["generations"] {
		if rid := toInt(r["region_id"]); rid > 0 && !regionIDs[rid] {
			res.warn("generations id=%d references unknown region_id=%d", toInt(r["id"]), rid)
		}
	}

	// species.generation_id → generations
	for _, r := range datasets["species"] {
		if gid := toInt(r["generation_id"]); gid > 0 && !generationIDs[gid] {
			res.warn("species id=%d references unknown generation_id=%d", toInt(r["id"]), gid)
		}
	}

	// abilities.introduced_generation_id → generations
	for _, r := range datasets["abilities"] {
		if gid := toInt(r["introduced_generation_id"]); gid > 0 && !generationIDs[gid] {
			res.warn("abilities id=%d references unknown introduced_generation_id=%d", toInt(r["id"]), gid)
		}
	}

	// pokemons FKs
	for _, r := range datasets["pokemons"] {
		pid := toInt(r["id"])
		if sid := toInt(r["stats_id"]); sid > 0 && !statsIDs[sid] {
			res.warn("pokemon id=%d references unknown stats_id=%d", pid, sid)
		}
		if gid := toInt(r["generation_id"]); gid > 0 && !generationIDs[gid] {
			res.warn("pokemon id=%d references unknown generation_id=%d", pid, gid)
		}
		if sid := toInt(r["species_id"]); sid > 0 && !speciesIDs[sid] {
			res.warn("pokemon id=%d references unknown species_id=%d", pid, sid)
		}
		if rid := toInt(r["region_id"]); rid > 0 && !regionIDs[rid] {
			res.warn("pokemon id=%d references unknown region_id=%d", pid, rid)
		}
		if ecid := toInt(r["evolution_chain_id"]); ecid > 0 && !evolutionChainIDs[ecid] {
			res.warn("pokemon id=%d references unknown evolution_chain_id=%d", pid, ecid)
		}
		// type_ids
		if typeIDsRaw, ok := r["type_ids"].([]any); ok {
			for _, tid := range typeIDsRaw {
				if t := toInt(tid); t > 0 && !typeIDs[t] {
					res.warn("pokemon id=%d references unknown type_id=%d", pid, t)
				}
			}
		}
		// abilities
		if abilitiesRaw, ok := r["abilities"].([]any); ok {
			for _, a := range abilitiesRaw {
				if ab, ok := a.(map[string]any); ok {
					if aid := toInt(ab["ability_id"]); aid > 0 && !abilityIDs[aid] {
						res.warn("pokemon id=%d references unknown ability_id=%d", pid, aid)
					}
				}
			}
		}
		// egg_group_ids
		if eggIDsRaw, ok := r["egg_group_ids"].([]any); ok {
			for _, eid := range eggIDsRaw {
				if e := toInt(eid); e > 0 && !eggGroupIDs[e] {
					res.warn("pokemon id=%d references unknown egg_group_id=%d", pid, e)
				}
			}
		}
	}

	// pokemon_weaknesses.pokemon_id → pokemons
	for _, r := range datasets["pokemon_weaknesses"] {
		if pid := toInt(r["pokemon_id"]); pid > 0 && !pokemonIDs[pid] {
			res.warn("pokemon_weaknesses references unknown pokemon_id=%d", pid)
		}
	}

	// Deep validation: Check evolution_chains for orphaned pokemon IDs in nested chain structures
	for _, r := range datasets["evolution_chains"] {
		chainID := toInt(r["id"])
		if chainRaw, ok := r["chain"].(map[string]any); ok {
			orphans := walkChain(chainRaw, pokemonIDs)
			for _, oid := range orphans {
				res.warn("evolution_chain id=%d references unknown pokemon_id=%d in nested chain", chainID, oid)
			}
		}
	}

	return res
}

// walkChain recursively finds all pokemon IDs in a chain structure and returns
// any IDs not found in validPokemonIDs.
func walkChain(chain map[string]any, validPokemonIDs map[int]bool) []int {
	var orphans []int

	// Check the pokemon at this level
	if pmon, ok := chain["pokemon"].(map[string]any); ok {
		if pid := toInt(pmon["id"]); pid > 0 && !validPokemonIDs[pid] {
			orphans = append(orphans, pid)
		}
	}

	// Recursively check evolution branches
	if evolutions, ok := chain["evolutions_to"].([]any); ok {
		for _, evolution := range evolutions {
			if evo, ok := evolution.(map[string]any); ok {
				orphans = append(orphans, walkChain(evo, validPokemonIDs)...)
			}
		}
	}

	return orphans
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
