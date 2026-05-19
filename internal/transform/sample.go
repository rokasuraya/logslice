// Package transform provides line transformation utilities for logslice.
package transform

import (
	"fmt"
	"sync/atomic"
)

// Sampler retains every Nth log line, dropping the rest.
// It is safe for concurrent use.
type Sampler struct {
	n       uint64
	counter atomic.Uint64
}

// NewSampler creates a Sampler that keeps every nth line.
// n must be >= 1; n=1 keeps every line (no-op sampling).
func NewSampler(n uint64) (*Sampler, error) {
	if n == 0 {
		return nil, fmt.Errorf("sample rate must be >= 1, got %d", n)
	}
	return &Sampler{n: n}, nil
}

// Apply returns (line, true) when the line should be kept, or ("", false)
// when it should be dropped. The counter increments on every call.
func (s *Sampler) Apply(line string) (string, bool) {
	// counter is 0-based; keep when (counter % n) == 0
	c := s.counter.Add(1) - 1
	if c%s.n == 0 {
		return line, true
	}
	return "", false
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() uint64 {
	return s.n
}

// Reset resets the internal counter back to zero.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}
