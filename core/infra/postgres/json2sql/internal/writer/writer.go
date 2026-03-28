package writer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EscapeValue converts a Go value to a SQL literal.
func EscapeValue(v any) string {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case bool:
		if val {
			return "TRUE"
		}
		return "FALSE"
	case json.Number:
		return val.String()
	case float64:
		// json.Unmarshal uses float64 for numbers; avoid scientific notation for ints.
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case string:
		escaped := strings.ReplaceAll(val, "'", "''")
		escaped = strings.ReplaceAll(escaped, "\n", `\n`)
		return "'" + escaped + "'"
	case map[string]any, []any:
		b, _ := json.Marshal(val)
		s := strings.ReplaceAll(string(b), "'", "''")
		return "'" + s + "'"
	default:
		s := fmt.Sprintf("%v", val)
		s = strings.ReplaceAll(s, "'", "''")
		return "'" + s + "'"
	}
}

// InsertRow builds a single INSERT statement for table using the ordered
// columns and values maps.
func InsertRow(table string, row map[string]any, columns []string) string {
	cols := make([]string, 0, len(columns))
	vals := make([]string, 0, len(columns))
	for _, c := range columns {
		if v, ok := row[c]; ok {
			cols = append(cols, c)
			vals = append(vals, EscapeValue(v))
		}
	}
	if len(cols) == 0 {
		return ""
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		table,
		strings.Join(cols, ", "),
		strings.Join(vals, ", "),
	)
}
