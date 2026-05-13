package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

// TestCollector_FullWorkflow simulates a typical logslice processing loop.
func TestCollector_FullWorkflow(t *testing.T) {
	c := stats.NewCollector()

	lines := []struct {
		matched bool
	}{
		{true}, {true}, {false}, {true}, {false}, {false},
	}

	for _, l := range lines {
		c.RecordTotal()
		if l.matched {
			c.RecordMatched()
		} else {
			c.RecordSkipped()
		}
	}

	s := c.Summary()

	if s.TotalLines != 6 {
		t.Errorf("expected TotalLines=6, got %d", s.TotalLines)
	}
	if s.MatchedLines != 3 {
		t.Errorf("expected MatchedLines=3, got %d", s.MatchedLines)
	}
	if s.SkippedLines != 3 {
		t.Errorf("expected SkippedLines=3, got %d", s.SkippedLines)
	}
	if s.MatchRate() != 50.0 {
		t.Errorf("expected MatchRate=50.0, got %f", s.MatchRate())
	}

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "logslice summary") {
		t.Errorf("expected summary header in output, got:\n%s", out)
	}
}
