package parser

import (
	"testing"
)

func TestNewLineParser_Defaults(t *testing.T) {
	lp := NewLineParser(nil, nil)
	if lp == nil {
		t.Fatal("expected non-nil LineParser")
	}
	if len(lp.Formats) != 0 {
		t.Errorf("expected no formats, got %d", len(lp.Formats))
	}
}

func TestLineParser_Parse_WithTimestamp(t *testing.T) {
	lp := NewLineParser(nil, nil)
	raw := `2024-01-15T10:30:00Z level=info msg="service started" service=api`
	line := lp.Parse(raw)

	if line.Raw != raw {
		t.Errorf("Raw mismatch: got %q", line.Raw)
	}
	if !line.HasTime {
		t.Error("expected HasTime=true")
	}
	if line.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if line.Fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", line.Fields["level"])
	}
}

func TestLineParser_Parse_NoTimestamp(t *testing.T) {
	lp := NewLineParser(nil, nil)
	raw := `level=warn msg="no timestamp here"`
	line := lp.Parse(raw)

	if line.HasTime {
		t.Error("expected HasTime=false")
	}
	if line.Fields["level"] != "warn" {
		t.Errorf("expected level=warn, got %q", line.Fields["level"])
	}
}

func TestLineParser_Matches_NoFilters(t *testing.T) {
	lp := NewLineParser(nil, nil)
	line := LogLine{Fields: map[string]string{"level": "error"}}
	if !lp.Matches(line) {
		t.Error("expected Matches=true with no filters")
	}
}

func TestLineParser_Matches_WithFilter(t *testing.T) {
	filter, err := ParseFieldFilter("level=info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lp := NewLineParser(nil, []FieldFilter{filter})

	match := LogLine{Fields: map[string]string{"level": "info", "service": "api"}}
	noMatch := LogLine{Fields: map[string]string{"level": "error"}}

	if !lp.Matches(match) {
		t.Error("expected match for level=info")
	}
	if lp.Matches(noMatch) {
		t.Error("expected no match for level=error")
	}
}

func TestLineParser_ParseLines(t *testing.T) {
	lp := NewLineParser(nil, nil)
	raws := []string{
		`2024-01-15T10:00:00Z level=info msg="start"`,
		``,
		`2024-01-15T10:01:00Z level=error msg="fail"`,
		"\r\n",
	}
	lines := lp.ParseLines(raws)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !lines[0].HasTime || !lines[1].HasTime {
		t.Error("expected both lines to have timestamps")
	}
}

func TestLineParser_ParseLines_WithFilter(t *testing.T) {
	filter, err := ParseFieldFilter("level=info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lp := NewLineParser(nil, []FieldFilter{filter})
	raws := []string{
		`2024-01-15T10:00:00Z level=info msg="start"`,
		`2024-01-15T10:01:00Z level=error msg="fail"`,
		`2024-01-15T10:02:00Z level=info msg="done"`,
	}
	lines := lp.ParseLines(raws)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after filtering, got %d", len(lines))
	}
	for _, line := range lines {
		if line.Fields["level"] != "info" {
			t.Errorf("expected only info lines, got level=%q", line.Fields["level"])
		}
	}
}
