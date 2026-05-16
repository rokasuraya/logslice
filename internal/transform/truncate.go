// Package transform provides utilities for modifying log line content
// before output, such as redacting sensitive fields or truncating long values.
package transform

import (
	"fmt"
	"strings"
)

// Truncator truncates field values in log lines that exceed a maximum length.
type Truncator struct {
	maxLen  int
	fields  []string
	suffix  string
}

// NewTruncator creates a Truncator that truncates the given fields to maxLen
// characters. If fields is empty, all detected values are truncated.
// A suffix (e.g. "...") is appended to indicate truncation.
func NewTruncator(maxLen int, fields []string, suffix string) (*Truncator, error) {
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be positive, got %d", maxLen)
	}
	if suffix == "" {
		suffix = "..."
	}
	copied := make([]string, len(fields))
	copy(copied, fields)
	return &Truncator{
		maxLen: maxLen,
		fields: copied,
		suffix: suffix,
	}, nil
}

// Apply returns the log line with targeted field values truncated.
// It supports key=value and key="value" styles as well as JSON-style "key":"value".
func (t *Truncator) Apply(line string) string {
	if len(t.fields) == 0 {
		return t.truncateAll(line)
	}
	for _, field := range t.fields {
		line = t.truncateField(line, field)
	}
	return line
}

// Fields returns the list of fields this truncator targets.
func (t *Truncator) Fields() []string {
	out := make([]string, len(t.fields))
	copy(out, t.fields)
	return out
}

func (t *Truncator) truncateField(line, field string) string {
	// Try key="value" style
	quotedPrefix := field + `="`
	if idx := strings.Index(line, quotedPrefix); idx != -1 {
		start := idx + len(quotedPrefix)
		end := strings.Index(line[start:], `"`)
		if end != -1 {
			val := line[start : start+end]
			return line[:start] + t.truncateValue(val) + line[start+end:]
		}
	}
	// Try key=value style
	kvPrefix := field + "="
	if idx := strings.Index(line, kvPrefix); idx != -1 {
		start := idx + len(kvPrefix)
		end := strings.IndexAny(line[start:], " \t,}")
		if end == -1 {
			end = len(line[start:])
		}
		val := line[start : start+end]
		return line[:start] + t.truncateValue(val) + line[start+end:]
	}
	return line
}

func (t *Truncator) truncateAll(line string) string {
	if len(line) <= t.maxLen {
		return line
	}
	return line[:t.maxLen] + t.suffix
}

func (t *Truncator) truncateValue(val string) string {
	if len(val) <= t.maxLen {
		return val
	}
	return val[:t.maxLen] + t.suffix
}
