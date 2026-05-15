// Package highlight provides ANSI color highlighting for matched log lines and fields.
package highlight

import (
	"fmt"
	"strings"
)

// ANSI color codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Mode controls whether highlighting is enabled.
type Mode int

const (
	ModeAuto Mode = iota // enable if stdout is a TTY
	ModeOn               // always enable
	ModeOff              // always disable
)

// Highlighter applies ANSI colors to log output.
type Highlighter struct {
	enabled bool
}

// New creates a Highlighter with the given mode.
// When mode is ModeAuto, enabled is set by the caller via isTTY.
func New(mode Mode, isTTY bool) *Highlighter {
	var enabled bool
	switch mode {
	case ModeOn:
		enabled = true
	case ModeOff:
		enabled = false
	default:
		enabled = isTTY
	}
	return &Highlighter{enabled: enabled}
}

// Line wraps an entire log line with the given color.
func (h *Highlighter) Line(line, color string) string {
	if !h.enabled {
		return line
	}
	return fmt.Sprintf("%s%s%s", color, line, Reset)
}

// Field highlights a specific key=value substring within a line.
func (h *Highlighter) Field(line, key string) string {
	if !h.enabled {
		return line
	}
	idx := strings.Index(line, key+"=")
	if idx == -1 {
		return line
	}
	// find end of value (space or end of string)
	start := idx
	end := strings.IndexByte(line[start:], ' ')
	if end == -1 {
		end = len(line)
	} else {
		end = start + end
	}
	return line[:start] + Bold + Cyan + line[start:end] + Reset + line[end:]
}

// Enabled reports whether highlighting is active.
func (h *Highlighter) Enabled() bool {
	return h.enabled
}
