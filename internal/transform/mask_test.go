package transform

import (
	"strings"
	"testing"
)

func TestNewMasker_InvalidReveal(t *testing.T) {
	_, err := NewMasker([]string{"token"}, -1, "*")
	if err == nil {
		t.Fatal("expected error for negative reveal, got nil")
	}
}

func TestNewMasker_DefaultMaskChar(t *testing.T) {
	m, err := NewMasker([]string{"token"}, 2, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.maskChar != "*" {
		t.Errorf("expected default maskChar '*', got %q", m.maskChar)
	}
}

func TestMasker_Apply_NoFields(t *testing.T) {
	m, _ := NewMasker(nil, 2, "*")
	line := "token=abc123 status=ok"
	if got := m.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestMasker_Apply_KeyValue_Unquoted(t *testing.T) {
	m, _ := NewMasker([]string{"token"}, 2, "*")
	got := m.Apply("token=abc123 status=ok")
	if !strings.Contains(got, "token=ab****") {
		t.Errorf("expected masked token, got %q", got)
	}
	if !strings.Contains(got, "status=ok") {
		t.Errorf("expected status unchanged, got %q", got)
	}
}

func TestMasker_Apply_KeyValue_Quoted(t *testing.T) {
	m, _ := NewMasker([]string{"secret"}, 3, "#")
	got := m.Apply(`secret="mysecretvalue" other=x`)
	if !strings.Contains(got, `secret="mys##########"`) {
		t.Errorf("expected masked quoted value, got %q", got)
	}
}

func TestMasker_Apply_JSONStyle(t *testing.T) {
	m, _ := NewMasker([]string{"password"}, 0, "*")
	got := m.Apply(`{"password":"hunter2","user":"alice"}`)
	if !strings.Contains(got, `"password":"*******"`) {
		t.Errorf("expected fully masked password, got %q", got)
	}
	if !strings.Contains(got, `"user":"alice"`) {
		t.Errorf("expected user unchanged, got %q", got)
	}
}

func TestMasker_Apply_RevealExceedsLength(t *testing.T) {
	m, _ := NewMasker([]string{"id"}, 100, "*")
	line := "id=abc"
	got := m.Apply(line)
	// value shorter than reveal — should be unchanged
	if !strings.Contains(got, "id=abc") {
		t.Errorf("expected unchanged short value, got %q", got)
	}
}

func TestMasker_Fields(t *testing.T) {
	m, _ := NewMasker([]string{"token", "secret"}, 2, "*")
	fields := m.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
}

func TestMasker_Apply_MultipleFields(t *testing.T) {
	m, _ := NewMasker([]string{"token", "pass"}, 1, "-")
	got := m.Apply("token=abcd pass=xyz")
	if !strings.Contains(got, "token=a---") {
		t.Errorf("token not masked correctly in %q", got)
	}
	if !strings.Contains(got, "pass=x--") {
		t.Errorf("pass not masked correctly in %q", got)
	}
}
