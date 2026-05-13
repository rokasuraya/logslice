package filter

import (
	"testing"
)

func TestNewFieldFilter_Empty(t *testing.T) {
	ff, err := NewFieldFilter(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ff.Len() != 0 {
		t.Errorf("expected 0 filters, got %d", ff.Len())
	}
}

func TestNewFieldFilter_InvalidExpr(t *testing.T) {
	_, err := NewFieldFilter([]string{"notavalidfilter"})
	if err == nil {
		t.Fatal("expected error for invalid filter expression, got nil")
	}
}

func TestFieldFilter_Accepts_NoFilters(t *testing.T) {
	ff, _ := NewFieldFilter(nil)
	if !ff.Accepts(`{"level":"error","msg":"boom"}`) {
		t.Error("expected line to be accepted when no filters are set")
	}
}

func TestFieldFilter_Accepts_MatchingField(t *testing.T) {
	ff, err := NewFieldFilter([]string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"error","msg":"something failed"}`
	if !ff.Accepts(line) {
		t.Errorf("expected line to be accepted, got rejected")
	}
}

func TestFieldFilter_Accepts_NonMatchingField(t *testing.T) {
	ff, err := NewFieldFilter([]string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"info","msg":"all good"}`
	if ff.Accepts(line) {
		t.Errorf("expected line to be rejected, got accepted")
	}
}

func TestFieldFilter_Accepts_MultipleFilters_AllMatch(t *testing.T) {
	ff, err := NewFieldFilter([]string{"level=error", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"error","service":"api","msg":"fail"}`
	if !ff.Accepts(line) {
		t.Errorf("expected line to be accepted when all filters match")
	}
}

func TestFieldFilter_Accepts_MultipleFilters_PartialMatch(t *testing.T) {
	ff, err := NewFieldFilter([]string{"level=error", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"error","service":"worker","msg":"fail"}`
	if ff.Accepts(line) {
		t.Errorf("expected line to be rejected when only partial filters match")
	}
}

func TestFieldFilter_Len(t *testing.T) {
	ff, err := NewFieldFilter([]string{"level=error", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ff.Len() != 2 {
		t.Errorf("expected Len()=2, got %d", ff.Len())
	}
}
