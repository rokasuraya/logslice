// Package transform provides composable line transformation utilities.
//
// # Sampler
//
// Sampler reduces log volume by retaining only every Nth line and discarding
// the rest. This is useful when dealing with extremely high-throughput logs
// where full fidelity is not required for analysis.
//
// Usage:
//
//	s, err := transform.NewSampler(10) // keep 1 in every 10 lines
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, line := range lines {
//		out, ok := s.Apply(line)
//		if ok {
//			fmt.Println(out)
//		}
//	}
//
// A rate of 1 is a no-op — every line is kept. A rate of 0 is invalid and
// returns an error from NewSampler.
//
// The counter can be reset between logical log segments via Reset().
package transform
