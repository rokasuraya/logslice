package parser

import (
	"strings"
	"time"
)

// LogLine represents a parsed log line with its timestamp and fields.
type LogLine struct {
	Raw       string
	Timestamp time.Time
	Fields    map[string]string
	HasTime   bool
}

// LineParser holds configuration for parsing log lines.
type LineParser struct {
	Formats []string
	Filters []FieldFilter
}

// NewLineParser creates a LineParser with the given timestamp formats and field filters.
func NewLineParser(formats []string, filters []FieldFilter) *LineParser {
	if len(formats) == 0 {
		formats = nil // will use auto-detection
	}
	return &LineParser{
		Formats: formats,
		Filters: filters,
	}
}

// Parse attempts to parse a raw log line into a LogLine.
// Timestamp parsing is attempted with provided formats or auto-detected.
func (lp *LineParser) Parse(raw string) LogLine {
	line := LogLine{
		Raw:    raw,
		Fields: ExtractFields(raw),
	}

	var ts time.Time
	var err error

	if len(lp.Formats) > 0 {
		for _, fmt := range lp.Formats {
			ts, err = ParseTimestampWithFormat(raw, fmt)
			if err == nil {
				line.Timestamp = ts
				line.HasTime = true
				break
			}
		}
	} else {
		ts, err = ParseTimestamp(raw)
		if err == nil {
			line.Timestamp = ts
			line.HasTime = true
		}
	}

	return line
}

// Matches returns true if the log line passes all configured field filters.
func (lp *LineParser) Matches(line LogLine) bool {
	if len(lp.Filters) == 0 {
		return true
	}
	return MatchesAllFilters(line.Fields, lp.Filters)
}

// ParseLines parses a slice of raw log lines and returns LogLine results.
func (lp *LineParser) ParseLines(raws []string) []LogLine {
	result := make([]LogLine, 0, len(raws))
	for _, raw := range raws {
		raw = strings.TrimRight(raw, "\r\n")
		if raw == "" {
			continue
		}
		result = append(result, lp.Parse(raw))
	}
	return result
}
