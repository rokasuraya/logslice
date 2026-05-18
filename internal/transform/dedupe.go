// Package transform provides line transformation utilities.
package transform

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// Deduplicator filters out duplicate log lines within a sliding window.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]struct{}
	window  int
	order   []string
}

// NewDeduplicator creates a Deduplicator that remembers up to windowSize
// recently seen lines. A windowSize of 0 means unlimited memory.
func NewDeduplicator(windowSize int) (*Deduplicator, error) {
	if windowSize < 0 {
		return nil, fmt.Errorf("dedupe: window size must be >= 0, got %d", windowSize)
	}
	return &Deduplicator{
		seen:   make(map[string]struct{}),
		window: windowSize,
	}, nil
}

// IsDuplicate returns true if the line has been seen recently.
// If it is not a duplicate, the line is recorded.
func (d *Deduplicator) IsDuplicate(line string) bool {
	key := hashLine(line)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[key]; exists {
		return true
	}

	d.seen[key] = struct{}{}
	d.order = append(d.order, key)

	if d.window > 0 && len(d.order) > d.window {
		evict := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, evict)
	}

	return false
}

// Reset clears all recorded lines.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
	d.order = d.order[:0]
}

// Len returns the number of unique lines currently tracked.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

func hashLine(line string) string {
	h := sha256.Sum256([]byte(line))
	return hex.EncodeToString(h[:])
}
