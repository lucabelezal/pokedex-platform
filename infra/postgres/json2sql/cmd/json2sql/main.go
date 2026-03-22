package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pokedex-platform/infra/postgres/json2sql/internal/loader"
	"pokedex-platform/infra/postgres/json2sql/internal/mapper"
	"pokedex-platform/infra/postgres/json2sql/internal/preflight"
	"pokedex-platform/infra/postgres/json2sql/internal/special"
	"pokedex-platform/infra/postgres/json2sql/internal/writer"
)

func main() {
	input := flag.String("input", "infra/postgres/source-json", "directory with source JSON files")
	output := flag.String("output", "infra/postgres/seeds/init-data.sql", "output SQL file path")
	strict := flag.Bool("strict", false, "fail on first referential inconsistency")
	flag.Parse()

	log.SetFlags(0)
	log.Printf("json2sql  input=%s  output=%s  strict=%v", *input, *output, *strict)

	// ── 1. Load all JSON files ────────────────────────────────────────────────
	fileResults, loadErrs := loader.LoadAll(*input, mapper.FileOrder)
	for _, e := range loadErrs {
		log.Printf("LOAD ERROR: %v", e)
	}
	if len(loadErrs) > 0 && *strict {
		os.Exit(1)
	}

	// Build a table→records index used by preflight and special handlers.
	byTable := make(map[string][]map[string]any, len(fileResults))
	for _, fr := range fileResults {
		table := mapper.FileToTable[fr.FileName]
		byTable[table] = fr.Records
	}

	// ── 2. Update type name→ID map from live data ─────────────────────────────
	if types, ok := byTable["types"]; ok {
		special.BuildTypeMap(types)
	}

	// ── 3. Preflight checks ───────────────────────────────────────────────────
	pfResult := preflight.Run(byTable)
	for _, w := range pfResult.Warnings {
		log.Printf("WARN: %s", w)
	}
	if pfResult.HasIssues() && *strict {
		log.Println("Aborting due to --strict mode.")
		os.Exit(1)
	}

	// ── 4. Generate SQL ───────────────────────────────────────────────────────
	var lines []string
	lines = append(lines,
		"-- init-data.sql",
		"-- Generated automatically from source JSON files by json2sql.",
		fmt.Sprintf("-- Generated at: %s", time.Now().Format(time.RFC3339)),
		"",
		"-- Start of data load",
		"",
	)

	successCount, errorCount := 0, 0

	for _, fr := range fileResults {
		table := mapper.FileToTable[fr.FileName]
		log.Printf("  processing %-30s -> %s (%d records)", fr.FileName, table, len(fr.Records))

		lines = append(lines, fmt.Sprintf("-- Table: %s  (source: %s)", table, fr.FileName))

		var stmts []string
		var warns []string
		var err error

		switch table {
		case "evolution_chains":
			stmts = special.EvolutionChains(fr.Records)
		case "pokemons":
			stmts = special.Pokemons(fr.Records)
		case "pokemon_weaknesses":
			stmts, warns = special.PokemonWeaknesses(fr.Records)
			for _, w := range warns {
				log.Printf("  WARN: %s", w)
			}
		default:
			validCols := mapper.TableValidFields[table]
			for _, rec := range fr.Records {
				filtered := mapper.FilterFields(rec, table)
				if s := writer.InsertRow(table, filtered, validCols); s != "" {
					stmts = append(stmts, s)
				}
			}
		}

		if err != nil {
			log.Printf("  ERROR: %v", err)
			errorCount++
			continue
		}

		lines = append(lines, stmts...)
		lines = append(lines, "")
		log.Printf("  -> %d statements", len(stmts))
		successCount++
	}

	lines = append(lines,
		"-- End of data load",
		fmt.Sprintf("-- Summary: %d files OK, %d errors", successCount, errorCount),
	)

	// ── 5. Write output file ──────────────────────────────────────────────────
	if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
		log.Fatalf("cannot create output directory: %v", err)
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(*output, []byte(content), 0o644); err != nil {
		log.Fatalf("cannot write output file: %v", err)
	}

	log.Printf("wrote %s (%d lines, %d bytes)", *output, len(lines), len(content))
	if errorCount > 0 {
		os.Exit(1)
	}
}
