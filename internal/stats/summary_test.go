package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSummary_MatchRate_Zero(t *testing.T) {
	s := Summary{}
	if s.MatchRate() != 0.0 {
		t.Errorf("expected 0.0 match rate for empty summary, got %f", s.MatchRate())
	}
}

func TestSummary_MatchRate_Half(t *testing.T) {
	s := Summary{TotalLines: 10, MatchedLines: 5}
	if s.MatchRate() != 50.0 {
		t.Errorf("expected 50.0, got %f", s.MatchRate())
	}
}

func TestSummary_MatchRate_Full(t *testing.T) {
	s := Summary{TotalLines: 4, MatchedLines: 4}
	if s.MatchRate() != 100.0 {
		t.Errorf("expected 100.0, got %f", s.MatchRate())
	}
}

func TestSummary_Print_ContainsFields(t *testing.T) {
	s := Summary{
		TotalLines:   20,
		MatchedLines: 15,
		SkippedLines: 5,
		Elapsed:      42 * time.Millisecond,
	}
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"20", "15", "5", "75.0%", "42ms"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestSummary_Print_ZeroElapsed(t *testing.T) {
	s := Summary{
		TotalLines:   10,
		MatchedLines: 10,
		SkippedLines: 0,
		Elapsed:      0,
	}
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"10", "100.0%", "0s"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
