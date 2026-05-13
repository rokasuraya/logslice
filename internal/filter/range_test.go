package filter

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewTimeRange_Unbounded(t *testing.T) {
	tr, err := NewTimeRange("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.IsUnbounded() {
		t.Error("expected unbounded range")
	}
}

func TestNewTimeRange_InvalidStart(t *testing.T) {
	_, err := NewTimeRange("not-a-time", "")
	if err == nil {
		t.Error("expected error for invalid start time")
	}
}

func TestTimeRange_Contains(t *testing.T) {
	start := mustTime("2024-01-01T10:00:00Z")
	end := mustTime("2024-01-01T12:00:00Z")
	tr := TimeRange{Start: &start, End: &end}

	cases := []struct {
		ts   time.Time
		want bool
	}{
		{mustTime("2024-01-01T09:59:59Z"), false},
		{mustTime("2024-01-01T10:00:00Z"), true},
		{mustTime("2024-01-01T11:00:00Z"), true},
		{mustTime("2024-01-01T12:00:00Z"), true},
		{mustTime("2024-01-01T12:00:01Z"), false},
	}

	for _, c := range cases {
		got := tr.Contains(c.ts)
		if got != c.want {
			t.Errorf("Contains(%v) = %v, want %v", c.ts, got, c.want)
		}
	}
}

func TestLineFilter_Accepts_NoFilters(t *testing.T) {
	lf, err := NewLineFilter("", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !lf.Accepts(nil, map[string]string{"level": "info"}) {
		t.Error("expected unbounded filter to accept any line")
	}
}

func TestLineFilter_Accepts_OutOfRange(t *testing.T) {
	lf, err := NewLineFilter("2024-01-01T10:00:00Z", "2024-01-01T11:00:00Z", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := mustTime("2024-01-01T09:00:00Z")
	if lf.Accepts(&ts, nil) {
		t.Error("expected line before range to be rejected")
	}
}

func TestLineFilter_Accepts_FieldMismatch(t *testing.T) {
	lf, err := NewLineFilter("", "", []string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lf.Accepts(nil, map[string]string{"level": "info"}) {
		t.Error("expected field mismatch to be rejected")
	}
}

func TestLineFilter_Accepts_Match(t *testing.T) {
	lf, err := NewLineFilter("2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z", []string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := mustTime("2024-01-01T11:00:00Z")
	if !lf.Accepts(&ts, map[string]string{"level": "error"}) {
		t.Error("expected matching line to be accepted")
	}
}
