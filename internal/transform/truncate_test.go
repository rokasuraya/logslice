package transform

import (
	"strings"
	"testing"
)

func TestNewTruncator_InvalidMaxLen(t *testing.T) {
	_, err := NewTruncator(0, nil, "")
	if err == nil {
		t.Fatal("expected error for maxLen=0")
	}
	_, err = NewTruncator(-5, nil, "")
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestNewTruncator_DefaultSuffix(t *testing.T) {
	tr, err := NewTruncator(10, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := tr.Apply("this is a very long line that should be cut")
	if !strings.HasSuffix(result, "...") {
		t.Errorf("expected default suffix '...', got: %q", result)
	}
}

func TestTruncator_Apply_NoFields_ShortLine(t *testing.T) {
	tr, _ := NewTruncator(100, nil, "...")
	line := "short line"
	got := tr.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestTruncator_Apply_NoFields_LongLine(t *testing.T) {
	tr, _ := NewTruncator(10, nil, "[cut]")
	got := tr.Apply("abcdefghijklmnopqrstuvwxyz")
	expected := "abcdefghij[cut]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestTruncator_Apply_KeyValue_Quoted(t *testing.T) {
	tr, _ := NewTruncator(5, []string{"msg"}, "...")
	line := `level=info msg="this is a long message" ts=2024-01-01`
	got := tr.Apply(line)
	if !strings.Contains(got, `msg="this ..."`) {
		t.Errorf("expected truncated quoted value, got: %q", got)
	}
	if !strings.Contains(got, "ts=2024-01-01") {
		t.Errorf("expected ts field untouched, got: %q", got)
	}
}

func TestTruncator_Apply_KeyValue_Unquoted(t *testing.T) {
	tr, _ := NewTruncator(4, []string{"token"}, "~")
	line := "user=alice token=supersecretvalue status=ok"
	got := tr.Apply(line)
	if !strings.Contains(got, "token=supe~") {
		t.Errorf("expected truncated token, got: %q", got)
	}
	if !strings.Contains(got, "user=alice") {
		t.Errorf("expected user field untouched, got: %q", got)
	}
}

func TestTruncator_Apply_FieldNotPresent(t *testing.T) {
	tr, _ := NewTruncator(5, []string{"secret"}, "...")
	line := "level=info msg=hello"
	got := tr.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line when field absent, got %q", got)
	}
}

func TestTruncator_Apply_ShortValueUnchanged(t *testing.T) {
	tr, _ := NewTruncator(20, []string{"msg"}, "...")
	line := `msg="hi there"`
	got := tr.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line for short value, got %q", got)
	}
}

func TestTruncator_Fields(t *testing.T) {
	fields := []string{"foo", "bar"}
	tr, _ := NewTruncator(10, fields, "")
	got := tr.Fields()
	if len(got) != len(fields) {
		t.Fatalf("expected %d fields, got %d", len(fields), len(got))
	}
	for i, f := range fields {
		if got[i] != f {
			t.Errorf("field[%d]: expected %q, got %q", i, f, got[i])
		}
	}
}
