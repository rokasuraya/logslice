package parser

import (
	"fmt"
	"time"
)

// Common timestamp formats found in structured logs.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.000000",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using a set of
// well-known formats. It returns the parsed time and the matched format,
// or an error if no format matched.
func ParseTimestamp(value string) (time.Time, string, error) {
	for _, format := range knownFormats {
		t, err := time.Parse(format, value)
		if err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("parser: unrecognized timestamp format: %q", value)
}

// ParseTimestampWithFormat parses a timestamp using an explicit format string.
func ParseTimestampWithFormat(value, format string) (time.Time, error) {
	t, err := time.Parse(format, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parser: cannot parse %q with format %q: %w", value, format, err)
	}
	return t, nil
}

// InRange reports whether ts falls within [from, to] (inclusive).
// A zero value for from or to means unbounded on that side.
func InRange(ts, from, to time.Time) bool {
	if !from.IsZero() && ts.Before(from) {
		return false
	}
	if !to.IsZero() && ts.After(to) {
		return false
	}
	return true
}
