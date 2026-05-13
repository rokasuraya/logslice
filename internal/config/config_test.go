package config

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	cfg := New()
	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("expected output format %q, got %q", DefaultOutputFormat, cfg.OutputFormat)
	}
	if cfg.BufferSize != DefaultBufferSize {
		t.Errorf("expected buffer size %d, got %d", DefaultBufferSize, cfg.BufferSize)
	}
	if cfg.CountOnly {
		t.Error("expected CountOnly to be false")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := New()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_BadBufferSize(t *testing.T) {
	cfg := New()
	cfg.BufferSize = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero buffer size")
	}
}

func TestValidate_BadOutputFormat(t *testing.T) {
	cfg := New()
	cfg.OutputFormat = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unsupported output format")
	}
}

func TestValidate_EndBeforeStart(t *testing.T) {
	cfg := New()
	start := time.Now()
	end := start.Add(-time.Hour)
	cfg.StartTime = &start
	cfg.EndTime = &end
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when end is before start")
	}
}

func TestHasTimeRange(t *testing.T) {
	cfg := New()
	if cfg.HasTimeRange() {
		t.Error("expected no time range on fresh config")
	}
	now := time.Now()
	cfg.StartTime = &now
	if !cfg.HasTimeRange() {
		t.Error("expected time range when StartTime is set")
	}
}
