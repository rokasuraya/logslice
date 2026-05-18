// Package transform provides line-level transformations for log output.
package transform

import (
	"fmt"
	"strings"
)

// Renamer renames fields in log lines, supporting both key=value and JSON-style
// log formats.
type Renamer struct {
	// mappings is a map from old field name to new field name.
	mappings map[string]string
}

// NewRenamer creates a Renamer from a slice of "old=new" expressions.
// Returns an error if any expression is malformed.
func NewRenamer(exprs []string) (*Renamer, error) {
	mappings := make(map[string]string, len(exprs))
	for _, expr := range exprs {
		parts := strings.SplitN(expr, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("rename: invalid expression %q: expected old=new", expr)
		}
		mappings[parts[0]] = parts[1]
	}
	return &Renamer{mappings: mappings}, nil
}

// Apply renames fields in the given log line. It handles both key=value and
// "key": value (JSON-style) formats. Returns the (possibly modified) line.
func (r *Renamer) Apply(line string) string {
	if len(r.mappings) == 0 {
		return line
	}
	for oldName, newName := range r.mappings {
		// key=value style
		kvPrefix := oldName + "="
		if idx := strings.Index(line, kvPrefix); idx != -1 {
			line = line[:idx] + newName + "=" + line[idx+len(kvPrefix):]
			continue
		}
		// JSON-style: "key":
		jsonKey := fmt.Sprintf("%q:", oldName)
		newKey := fmt.Sprintf("%q:", newName)
		line = strings.ReplaceAll(line, jsonKey, newKey)
	}
	return line
}

// Fields returns the set of old field names that will be renamed.
func (r *Renamer) Fields() []string {
	fields := make([]string, 0, len(r.mappings))
	for k := range r.mappings {
		fields = append(fields, k)
	}
	return fields
}
