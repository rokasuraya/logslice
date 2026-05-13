package stats

import (
	"fmt"
	"io"
	"time"
)

// Summary holds a snapshot of processing statistics.
type Summary struct {
	TotalLines   int64
	MatchedLines int64
	SkippedLines int64
	Elapsed      time.Duration
}

// MatchRate returns the percentage of lines that matched, or 0 if no lines.
func (s Summary) MatchRate() float64 {
	if s.TotalLines == 0 {
		return 0.0
	}
	return float64(s.MatchedLines) / float64(s.TotalLines) * 100.0
}

// Print writes a human-readable summary to w.
func (s Summary) Print(w io.Writer) {
	fmt.Fprintf(w, "--- logslice summary ---\n")
	fmt.Fprintf(w, "total lines   : %d\n", s.TotalLines)
	fmt.Fprintf(w, "matched lines : %d\n", s.MatchedLines)
	fmt.Fprintf(w, "skipped lines : %d\n", s.SkippedLines)
	fmt.Fprintf(w, "match rate    : %.1f%%\n", s.MatchRate())
	fmt.Fprintf(w, "elapsed       : %s\n", s.Elapsed.Round(time.Millisecond))
}
