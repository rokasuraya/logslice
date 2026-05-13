package config

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, args, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected no positional args, got %v", args)
	}
	if cfg.OutputFormat != "raw" {
		t.Errorf("expected raw output, got %q", cfg.OutputFormat)
	}
	if cfg.CountOnly {
		t.Error("expected CountOnly false")
	}
}

func TestParseFlags_CountFlag(t *testing.T) {
	cfg, _, err := ParseFlags([]string{"-count"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.CountOnly {
		t.Error("expected CountOnly true")
	}
}

func TestParseFlags_OutputJSON(t *testing.T) {
	cfg, _, err := ParseFlags([]string{"-output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected json, got %q", cfg.OutputFormat)
	}
}

func TestParseFlags_InvalidOutput(t *testing.T) {
	_, _, err := ParseFlags([]string{"-output", "csv"})
	if err == nil {
		t.Error("expected error for unsupported output format")
	}
}

func TestParseFlags_PositionalArgs(t *testing.T) {
	_, args, err := ParseFlags([]string{"file1.log", "file2.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 positional args, got %d", len(args))
	}
}

func TestParseFlags_FieldFilter(t *testing.T) {
	cfg, _, err := ParseFlags([]string{"-field", "level=error", "-field", "app=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(cfg.Fields))
	}
}

func TestParseFlags_InvalidStart(t *testing.T) {
	_, _, err := ParseFlags([]string{"-start", "not-a-time"})
	if err == nil {
		t.Error("expected error for invalid start timestamp")
	}
}
