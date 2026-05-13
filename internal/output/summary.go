package output

import (
	"fmt"
	"io"
	"time"
)

// Summary holds statistics collected during a log slice run.
type Summary struct {
	LinesRead    int
	LinesMatched int
	Duration     time.Duration
	InputSource  string
}

// Print writes a human-readable summary to the given writer.
func (s *Summary) Print(w io.Writer) {
	fmt.Fprintf(w, "source:  %s\n", s.InputSource)
	fmt.Fprintf(w, "read:    %d lines\n", s.LinesRead)
	fmt.Fprintf(w, "matched: %d lines\n", s.LinesMatched)
	fmt.Fprintf(w, "elapsed: %s\n", s.Duration.Round(time.Millisecond))
}

// MatchRate returns the fraction of lines that matched, or 0 if none were read.
func (s *Summary) MatchRate() float64 {
	if s.LinesRead == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesRead)
}
