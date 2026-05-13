package stats

import (
	"sync/atomic"
	"time"
)

// Collector tracks processing statistics for a logslice run.
type Collector struct {
	totalLines   atomic.Int64
	matchedLines atomic.Int64
	skippedLines atomic.Int64
	startTime     time.Time
}

// NewCollector creates a new Collector and records the start time.
func NewCollector() *Collector {
	return &Collector{
		startTime: time.Now(),
	}
}

// RecordTotal increments the total lines seen counter.
func (c *Collector) RecordTotal() {
	c.totalLines.Add(1)
}

// RecordMatched increments the matched lines counter.
func (c *Collector) RecordMatched() {
	c.matchedLines.Add(1)
}

// RecordSkipped increments the skipped lines counter.
func (c *Collector) RecordSkipped() {
	c.skippedLines.Add(1)
}

// Summary returns a snapshot of the current statistics.
func (c *Collector) Summary() Summary {
	return Summary{
		TotalLines:   c.totalLines.Load(),
		MatchedLines: c.matchedLines.Load(),
		SkippedLines: c.skippedLines.Load(),
		Elapsed:      time.Since(c.startTime),
	}
}
