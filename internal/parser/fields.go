package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldFilter represents a key=value filter to apply against log fields.
type FieldFilter struct {
	Key   string
	Value string
}

// ParseFieldFilter parses a filter expression in the form "key=value".
func ParseFieldFilter(expr string) (FieldFilter, error) {
	parts := strings.SplitN(expr, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		return FieldFilter{}, fmt.Errorf("invalid filter expression %q: expected key=value", expr)
	}
	return FieldFilter{Key: parts[0], Value: parts[1]}, nil
}

// ExtractFields parses a JSON log line and returns its fields as a map.
func ExtractFields(line string) (map[string]string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fields := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			fields[k] = val
		case nil:
			fields[k] = ""
		default:
			fields[k] = fmt.Sprintf("%v", val)
		}
	}
	return fields, nil
}

// MatchesFilter reports whether the given fields satisfy the FieldFilter.
func MatchesFilter(fields map[string]string, f FieldFilter) bool {
	v, ok := fields[f.Key]
	return ok && v == f.Value
}

// MatchesAllFilters reports whether the given fields satisfy all provided filters.
func MatchesAllFilters(fields map[string]string, filters []FieldFilter) bool {
	for _, f := range filters {
		if !MatchesFilter(fields, f) {
			return false
		}
	}
	return true
}
