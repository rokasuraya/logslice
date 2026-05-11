package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp_KnownFormats(t *testing.T) {
	cases := []struct {
		input  string
		wantOK bool
	}{
		{"2024-05-01T12:00:00Z", true},
		{"2024-05-01T12:00:00.123456789Z", true},
		{"2024-05-01T12:00:00.123+02:00", true},
		{"2024-05-01 12:00:00", true},
		{"2024-05-01 12:00:00.000", true},
		{"2024/05/01 12:00:00", true},
		{"not-a-timestamp", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			_, _, err := ParseTimestamp(tc.input)
			if tc.wantOK && err != nil {
				t.Errorf("expected success, got error: %v", err)
			}
			if !tc.wantOK && err == nil {
				t.Errorf("expected error for input %q, got none", tc.input)
			}
		})
	}
}

func TestParseTimestampWithFormat(t *testing.T) {
	_, err := ParseTimestampWithFormat("2024-05-01 12:00:00", "2006-01-02 15:04:05")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ParseTimestampWithFormat("bad", "2006-01-02")
	if err == nil {
		t.Fatal("expected error for bad input, got none")
	}
}

func TestInRange(t *testing.T) {
	base := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	before := base.Add(-time.Hour)
	after := base.Add(time.Hour)

	if !InRange(base, before, after) {
		t.Error("expected base to be in range [before, after]")
	}
	if InRange(before.Add(-time.Second), before, after) {
		t.Error("expected value before 'from' to be out of range")
	}
	if InRange(after.Add(time.Second), before, after) {
		t.Error("expected value after 'to' to be out of range")
	}
	// Unbounded from
	if !InRange(before, time.Time{}, after) {
		t.Error("expected unbounded-from range to include before")
	}
	// Unbounded to
	if !InRange(after, before, time.Time{}) {
		t.Error("expected unbounded-to range to include after")
	}
	// Both unbounded
	if !InRange(base, time.Time{}, time.Time{}) {
		t.Error("expected fully unbounded range to include any value")
	}
}
