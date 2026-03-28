package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// FileRecords holds the parsed records from a single JSON file.
type FileRecords struct {
	FileName string
	Records  []map[string]any
}

// LoadAll reads every JSON file in dir whose name matches a key in fileOrder,
// returning them in sorted order (01_… before 02_… etc.).
func LoadAll(dir string, fileOrder []string) ([]FileRecords, []error) {
	var results []FileRecords
	var errs []error

	sorted := make([]string, len(fileOrder))
	copy(sorted, fileOrder)
	sort.Strings(sorted)

	for _, name := range sorted {
		path := filepath.Join(dir, name)
		records, err := loadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			continue
		}
		results = append(results, FileRecords{FileName: name, Records: records})
	}
	return results, errs
}

// LoadFile loads a single JSON file by path (exported for dynamic lookups).
func LoadFile(path string) ([]map[string]any, error) {
	return loadFile(path)
}

func loadFile(path string) ([]map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	// Try as array first.
	var arr []map[string]any
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	// Fallback: single object wrapped as array.
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return []map[string]any{obj}, nil
}
