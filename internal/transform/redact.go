// Package transform provides log line transformation utilities such as
// field redaction before output.
package transform

import (
	"regexp"
	"strings"
)

// Redactor replaces sensitive field values in log lines with a placeholder.
type Redactor struct {
	fields      []string
	placeholder string
	patterns    []*regexp.Regexp
}

// NewRedactor creates a Redactor that masks the given field names.
// fields should be bare key names (e.g. "password", "token").
// placeholder is the string substituted for matched values (e.g. "[REDACTED]").
func NewRedactor(fields []string, placeholder string) (*Redactor, error) {
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}
	patterns := make([]*regexp.Regexp, 0, len(fields))
	for _, f := range fields {
		if f == "" {
			continue
		}
		// Matches key=value and "key":"value" / "key": "value" styles.
		raw := `(?i)(` + regexp.QuoteMeta(f) + `)(\s*[:=]\s*)("[^"]*"|\S+)`
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, re)
	}
	return &Redactor{
		fields:      fields,
		placeholder: placeholder,
		patterns:    patterns,
	}, nil
}

// Apply returns a copy of line with sensitive field values replaced.
func (r *Redactor) Apply(line string) string {
	if len(r.patterns) == 0 {
		return line
	}
	out := line
	for _, re := range r.patterns {
		out = re.ReplaceAllStringFunc(out, func(match string) string {
			// Preserve the key and separator; replace only the value.
			idx := re.FindStringSubmatchIndex(match)
			if idx == nil || len(idx) < 8 {
				return match
			}
			key := match[idx[2]:idx[3]]
			sep := match[idx[4]:idx[5]]
			// Preserve quoting style if value was quoted.
			val := match[idx[6]:idx[7]]
			if strings.HasPrefix(val, `"`) {
				return key + sep + `"` + r.placeholder + `"`
			}
			return key + sep + r.placeholder
		})
	}
	return out
}

// Fields returns the list of field names being redacted.
func (r *Redactor) Fields() []string {
	return append([]string(nil), r.fields...)
}
