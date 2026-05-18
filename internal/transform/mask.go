// Package transform provides line transformation utilities for logslice.
package transform

import (
	"fmt"
	"regexp"
	"strings"
)

// Masker partially masks field values, revealing only a configurable number
// of leading characters and replacing the rest with a mask character.
type Masker struct {
	fields   map[string]struct{}
	reveal   int
	maskChar string
}

// NewMasker creates a Masker for the given field names. reveal controls how
// many leading characters of each value remain visible. maskChar is the
// string used to replace hidden characters (defaults to "*").
func NewMasker(fields []string, reveal int, maskChar string) (*Masker, error) {
	if reveal < 0 {
		return nil, fmt.Errorf("reveal must be >= 0, got %d", reveal)
	}
	if maskChar == "" {
		maskChar = "*"
	}
	f := make(map[string]struct{}, len(fields))
	for _, name := range fields {
		if name == "" {
			continue
		}
		f[name] = struct{}{}
	}
	return &Masker{fields: f, reveal: reveal, maskChar: maskChar}, nil
}

// Apply returns a copy of line with targeted field values partially masked.
// It handles both key=value and key="value" formats as well as JSON-style
// "key":"value" pairs.
func (m *Masker) Apply(line string) string {
	if len(m.fields) == 0 {
		return line
	}
	for field := range m.fields {
		line = m.maskKV(line, field)
		line = m.maskJSON(line, field)
	}
	return line
}

// Fields returns the set of field names targeted by this Masker.
func (m *Masker) Fields() []string {
	out := make([]string, 0, len(m.fields))
	for f := range m.fields {
		out = append(out, f)
	}
	return out
}

var kvMaskCache = map[string]*regexp.Regexp{}

func (m *Masker) maskKV(line, field string) string {
	pattern := fmt.Sprintf(`(?i)(%s=)("?)([^"\s]+)("?)`, regexp.QuoteMeta(field))
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(line, func(match string) string {
		subs := re.FindStringSubmatch(match)
		if len(subs) < 5 {
			return match
		}
		prefix, openQ, val, closeQ := subs[1], subs[2], subs[3], subs[4]
		return prefix + openQ + m.maskValue(val) + closeQ
	})
}

func (m *Masker) maskJSON(line, field string) string {
	pattern := fmt.Sprintf(`(?i)("?%s"?\s*:\s*"?)([^"\s,}]+)("?)`, regexp.QuoteMeta(field))
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(line, func(match string) string {
		subs := re.FindStringSubmatch(match)
		if len(subs) < 4 {
			return match
		}
		prefix, val, closeQ := subs[1], subs[2], subs[3]
		return prefix + m.maskValue(val) + closeQ
	})
}

func (m *Masker) maskValue(val string) string {
	runes := []rune(val)
	if m.reveal >= len(runes) {
		return val
	}
	visible := string(runes[:m.reveal])
	hidden := strings.Repeat(m.maskChar, len(runes)-m.reveal)
	return visible + hidden
}
