package transform

import (
	"sort"
	"testing"
)

func TestNewRenamer_Empty(t *testing.T) {
	r, err := NewRenamer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Fields()) != 0 {
		t.Errorf("expected no fields, got %v", r.Fields())
	}
}

func TestNewRenamer_InvalidExpr(t *testing.T) {
	cases := []string{"noequals", "=newname", "oldname=", "="}
	for _, expr := range cases {
		_, err := NewRenamer([]string{expr})
		if err == nil {
			t.Errorf("expected error for expr %q, got nil", expr)
		}
	}
}

func TestRenamer_Apply_NoMappings(t *testing.T) {
	r, _ := NewRenamer(nil)
	line := "level=info msg=hello"
	if got := r.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestRenamer_Apply_KeyValue(t *testing.T) {
	r, err := NewRenamer([]string{"level=severity"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "level=info msg=hello"
	want := "severity=info msg=hello"
	if got := r.Apply(line); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestRenamer_Apply_JSONStyle(t *testing.T) {
	r, err := NewRenamer([]string{"level=severity"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"info","msg":"hello"}`
	want := `{"severity":"info","msg":"hello"}`
	if got := r.Apply(line); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestRenamer_Apply_NoMatch(t *testing.T) {
	r, _ := NewRenamer([]string{"missing=other"})
	line := "level=info msg=hello"
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line %q, got %q", line, got)
	}
}

func TestRenamer_Apply_MultipleFields(t *testing.T) {
	r, err := NewRenamer([]string{"level=severity", "msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "level=warn msg=oops"
	want := "severity=warn message=oops"
	if got := r.Apply(line); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestRenamer_Fields(t *testing.T) {
	r, _ := NewRenamer([]string{"a=x", "b=y"})
	fields := r.Fields()
	sort.Strings(fields)
	if len(fields) != 2 || fields[0] != "a" || fields[1] != "b" {
		t.Errorf("unexpected fields: %v", fields)
	}
}
