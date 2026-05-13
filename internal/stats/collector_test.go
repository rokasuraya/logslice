package stats

import (
	"testing"
	"time"
)

func TestNewCollector_InitialZeros(t *testing.T) {
	c := NewCollector()
	s := c.Summary()
	if s.TotalLines != 0 || s.MatchedLines != 0 || s.SkippedLines != 0 {
		t.Errorf("expected all zero counters, got total=%d matched=%d skipped=%d",
			s.TotalLines, s.MatchedLines, s.SkippedLines)
	}
}

func TestCollector_RecordTotal(t *testing.T) {
	c := NewCollector()
	c.RecordTotal()
	c.RecordTotal()
	if got := c.Summary().TotalLines; got != 2 {
		t.Errorf("expected TotalLines=2, got %d", got)
	}
}

func TestCollector_RecordMatched(t *testing.T) {
	c := NewCollector()
	c.RecordMatched()
	if got := c.Summary().MatchedLines; got != 1 {
		t.Errorf("expected MatchedLines=1, got %d", got)
	}
}

func TestCollector_RecordSkipped(t *testing.T) {
	c := NewCollector()
	c.RecordSkipped()
	c.RecordSkipped()
	c.RecordSkipped()
	if got := c.Summary().SkippedLines; got != 3 {
		t.Errorf("expected SkippedLines=3, got %d", got)
	}
}

func TestCollector_ElapsedPositive(t *testing.T) {
	c := NewCollector()
	time.Sleep(2 * time.Millisecond)
	if c.Summary().Elapsed <= 0 {
		t.Error("expected positive elapsed duration")
	}
}
