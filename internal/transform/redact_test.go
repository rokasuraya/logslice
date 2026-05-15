package transform_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/transform"
)

func TestNewRedactor_EmptyFields(t *testing.T) {
	r, err := transform.NewRedactor(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `level=info msg="login" password=secret`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestRedactor_Apply_KeyValue(t *testing.T) {
	r, err := transform.NewRedactor([]string{"password", "token"}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := `level=info msg=login password=s3cr3t token=abc123 user=alice`
	got := r.Apply(input)
	if contains(got, "s3cr3t") {
		t.Errorf("password value should be redacted, got: %s", got)
	}
	if contains(got, "abc123") {
		t.Errorf("token value should be redacted, got: %s", got)
	}
	if !contains(got, "alice") {
		t.Errorf("user value should be preserved, got: %s", got)
	}
}

func TestRedactor_Apply_JSONStyle(t *testing.T) {
	r, err := transform.NewRedactor([]string{"password"}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := `{"level":"info","password":"hunter2","user":"bob"}`
	got := r.Apply(input)
	if contains(got, "hunter2") {
		t.Errorf("password value should be redacted, got: %s", got)
	}
	if !contains(got, "bob") {
		t.Errorf("user value should be preserved, got: %s", got)
	}
}

func TestRedactor_Apply_DefaultPlaceholder(t *testing.T) {
	r, err := transform.NewRedactor([]string{"secret"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("secret=myvalue")
	if !contains(got, "[REDACTED]") {
		t.Errorf("expected default placeholder, got: %s", got)
	}
}

func TestRedactor_Fields(t *testing.T) {
	fields := []string{"password", "token"}
	r, _ := transform.NewRedactor(fields, "")
	got := r.Fields()
	if len(got) != len(fields) {
		t.Fatalf("expected %d fields, got %d", len(fields), len(got))
	}
	for i, f := range fields {
		if got[i] != f {
			t.Errorf("field[%d]: want %q, got %q", i, f, got[i])
		}
	}
}

func TestRedactor_Apply_CaseInsensitive(t *testing.T) {
	r, err := transform.NewRedactor([]string{"Password"}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("PASSWORD=topsecret")
	if contains(got, "topsecret") {
		t.Errorf("expected case-insensitive redaction, got: %s", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
