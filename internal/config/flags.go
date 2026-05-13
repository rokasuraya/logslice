package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// ParseFlags parses command-line flags into a Config and returns any
// remaining positional arguments (e.g. input file paths).
func ParseFlags(args []string) (*Config, []string, error) {
	cfg := New()

	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var startStr, endStr, format, outputFmt string
	var fields sliceFlag

	fs.StringVar(&startStr, "start", "", "start timestamp (inclusive)")
	fs.StringVar(&endStr, "end", "", "end timestamp (inclusive)")
	fs.StringVar(&format, "format", "", "timestamp format (Go layout or named format)")
	fs.StringVar(&outputFmt, "output", DefaultOutputFormat, "output format: raw or json")
	fs.Var(&fields, "field", "field filter key=value (repeatable)")
	fs.BoolVar(&cfg.CountOnly, "count", false, "print match count instead of lines")
	fs.IntVar(&cfg.BufferSize, "buf", DefaultBufferSize, "line buffer size in bytes")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	cfg.OutputFormat = outputFmt
	cfg.Format = format
	cfg.Fields = []string(fields)

	if startStr != "" {
		t, err := parser.ParseTimestampWithFormat(startStr, format)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid --start: %w", err)
		}
		cfg.StartTime = &t
	}
	if endStr != "" {
		t, err := parser.ParseTimestampWithFormat(endStr, format)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid --end: %w", err)
		}
		_ = time.Time{} // ensure time import used
		cfg.EndTime = &t
	}

	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}

	return cfg, fs.Args(), nil
}

// sliceFlag is a flag.Value that accumulates repeated string flags.
type sliceFlag []string

func (s *sliceFlag) String() string { return fmt.Sprintf("%v", *s) }
func (s *sliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}
