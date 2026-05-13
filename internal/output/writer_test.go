package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewWriter_DefaultsToRaw(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "")
	if w.format != FormatRaw {
		t.Fatalf("expected FormatRaw, got %q", w.format)
	}
}

func TestWriter_WriteLine_Raw(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatRaw)

	if err := w.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if got != "hello world\n" {
		t.Errorf("expected %q, got %q", "hello world\n", got)
	}
}

func TestWriter_WriteLine_JSON(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON)

	if err := w.WriteLine(`{"ts":"2024-01-01","msg":"ok"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = w.Flush()

	got := buf.String()
	if !strings.HasPrefix(got, `{"line":`) {
		t.Errorf("expected JSON wrapper, got %q", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "}") {
		t.Errorf("expected closing brace, got %q", got)
	}
}

func TestWriter_Count(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatRaw)

	for i := 0; i < 5; i++ {
		_ = w.WriteLine("line")
	}
	_ = w.Flush()

	if w.Count() != 5 {
		t.Errorf("expected count 5, got %d", w.Count())
	}
}

func TestWriter_MultipleLines_Raw(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatRaw)

	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}
	_ = w.Flush()

	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got[i])
		}
	}
}
