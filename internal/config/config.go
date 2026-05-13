package config

import (
	"errors"
	"time"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	StartTime    *time.Time
	EndTime      *time.Time
	Format       string
	Fields       []string
	OutputFormat string
	CountOnly    bool
	BufferSize   int
	InputFile    string
}

// DefaultBufferSize is the default line buffer size in bytes.
const DefaultBufferSize = 64 * 1024

// DefaultOutputFormat is the default output format.
const DefaultOutputFormat = "raw"

// New returns a Config with sensible defaults applied.
func New() *Config {
	return &Config{
		OutputFormat: DefaultOutputFormat,
		BufferSize:   DefaultBufferSize,
	}
}

// Validate checks that the configuration is internally consistent.
func (c *Config) Validate() error {
	if c.BufferSize <= 0 {
		return errors.New("buffer size must be greater than zero")
	}
	if c.OutputFormat != "raw" && c.OutputFormat != "json" {
		return errors.New("output format must be \"raw\" or \"json\"")
	}
	if c.StartTime != nil && c.EndTime != nil {
		if c.EndTime.Before(*c.StartTime) {
			return errors.New("end time must not be before start time")
		}
	}
	return nil
}

// HasTimeRange reports whether at least one time bound is set.
func (c *Config) HasTimeRange() bool {
	return c.StartTime != nil || c.EndTime != nil
}
