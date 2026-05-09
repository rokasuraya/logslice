package parser

import (
	"testing"
)

func TestParseFieldFilter(t *testing.T) {
	tests := []struct {
		expr    string
		wantKey string
		wantVal string
		wantErr bool
	}{
		{"level=info", "level", "info", false},
		{"service=auth", "service", "auth", false},
		{"msg=hello=world", "msg", "hello=world", false},
		{"noequals", "", "", true},
		{"=value", "", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			f, err := ParseFieldFilter(tt.expr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.expr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.Key != tt.wantKey || f.Value != tt.wantVal {
				t.Errorf("got {%q, %q}, want {%q, %q}", f.Key, f.Value, tt.wantKey, tt.wantVal)
			}
		})
	}
}

func TestExtractFields(t *testing.T) {
	line := `{"level":"error","msg":"disk full","code":500}`
	fields, err := ExtractFields(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["level"] != "error" {
		t.Errorf("expected level=error, got %q", fields["level"])
	}
	if fields["msg"] != "disk full" {
		t.Errorf("expected msg='disk full', got %q", fields["msg"])
	}
	if fields["code"] != "500" {
		t.Errorf("expected code=500, got %q", fields["code"])
	}

	_, err = ExtractFields("not json")
	if err == nil {
		t.Error("expected error for non-JSON input")
	}

	_, err = ExtractFields("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestMatchesAllFilters(t *testing.T) {
	fields := map[string]string{"level": "info", "service": "api"}

	if !MatchesAllFilters(fields, []FieldFilter{{"level", "info"}, {"service", "api"}}) {
		t.Error("expected all filters to match")
	}
	if MatchesAllFilters(fields, []FieldFilter{{"level", "info"}, {"service", "db"}}) {
		t.Error("expected filter mismatch for service=db")
	}
	if !MatchesAllFilters(fields, []FieldFilter{}) {
		t.Error("expected empty filters to match everything")
	}
}
