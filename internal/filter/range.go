package filter

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// TimeRange holds an optional start and end boundary for log filtering.
type TimeRange struct {
	Start *time.Time
	End   *time.Time
}

// NewTimeRange constructs a TimeRange from optional start/end strings.
// An empty string means the boundary is unbounded.
func NewTimeRange(startStr, endStr string) (TimeRange, error) {
	var tr TimeRange

	if startStr != "" {
		t, err := parser.ParseTimestamp(startStr)
		if err != nil {
			return tr, fmt.Errorf("invalid start time %q: %w", startStr, err)
		}
		tr.Start = &t
	}

	if endStr != "" {
		t, err := parser.ParseTimestamp(endStr)
		if err != nil {
			return tr, fmt.Errorf("invalid end time %q: %w", endStr, err)
		}
		tr.End = &t
	}

	return tr, nil
}

// Contains reports whether the given time falls within the range (inclusive).
func (tr TimeRange) Contains(t time.Time) bool {
	if tr.Start != nil && t.Before(*tr.Start) {
		return false
	}
	if tr.End != nil && t.After(*tr.End) {
		return false
	}
	return true
}

// IsUnbounded reports whether the range has no boundaries set.
func (tr TimeRange) IsUnbounded() bool {
	return tr.Start == nil && tr.End == nil
}

// LineFilter combines a TimeRange with field filters to decide whether
// a parsed log line should be included in output.
type LineFilter struct {
	Range       TimeRange
	FieldFilter []parser.FieldFilter
}

// NewLineFilter builds a LineFilter from raw CLI arguments.
func NewLineFilter(startStr, endStr string, fields []string) (LineFilter, error) {
	tr, err := NewTimeRange(startStr, endStr)
	if err != nil {
		return LineFilter{}, err
	}

	var ff []parser.FieldFilter
	for _, raw := range fields {
		f, err := parser.ParseFieldFilter(raw)
		if err != nil {
			return LineFilter{}, fmt.Errorf("invalid field filter %q: %w", raw, err)
		}
		ff = append(ff, f)
	}

	return LineFilter{Range: tr, FieldFilter: ff}, nil
}

// Accepts returns true when the line's timestamp is within the range AND
// all field filters match.
func (lf LineFilter) Accepts(ts *time.Time, fields map[string]string) bool {
	if ts != nil && !lf.Range.IsUnbounded() {
		if !lf.Range.Contains(*ts) {
			return false
		}
	}
	return parser.MatchesAllFilters(fields, lf.FieldFilter)
}
