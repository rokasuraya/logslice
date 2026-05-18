// Package transform provides line transformation utilities for logslice.
//
// # Deduplicator
//
// Deduplicator removes repeated log lines within a configurable sliding
// window. This is useful when tailing noisy logs that emit the same
// message many times in succession.
//
// Usage:
//
//	d, err := transform.NewDeduplicator(1000) // remember last 1000 lines
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, line := range lines {
//		if !d.IsDuplicate(line) {
//			fmt.Println(line)
//		}
//	}
//
// A window size of 0 disables eviction, keeping all seen lines in memory
// for the lifetime of the Deduplicator. Call Reset to free memory and
// start a fresh deduplication pass.
package transform
