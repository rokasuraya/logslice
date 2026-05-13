package filter

import (
	"github.com/user/logslice/internal/parser"
)

// FieldFilter holds compiled field filter expressions and applies them
// to parsed log lines.
type FieldFilter struct {
	filters []parser.FieldFilter
}

// NewFieldFilter creates a FieldFilter from a slice of raw filter strings
// (e.g. ["level=error", "service=api"]). Returns an error if any expression
// cannot be parsed.
func NewFieldFilter(exprs []string) (*FieldFilter, error) {
	filters := make([]parser.FieldFilter, 0, len(exprs))
	for _, expr := range exprs {
		f, err := parser.ParseFieldFilter(expr)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return &FieldFilter{filters: filters}, nil
}

// Accepts returns true when the given raw log line satisfies all field
// filters. Lines are accepted unconditionally when no filters are configured.
func (ff *FieldFilter) Accepts(line string) bool {
	if len(ff.filters) == 0 {
		return true
	}
	fields := parser.ExtractFields(line)
	return parser.MatchesAllFilters(fields, ff.filters)
}

// Len returns the number of active field filters.
func (ff *FieldFilter) Len() int {
	return len(ff.filters)
}
