// Package transform provides post-parse log line transformations applied
// before output. Current transformations include:
//
//   - Redaction: replace sensitive field values (e.g. passwords, tokens)
//     with a configurable placeholder string so that filtered output never
//     leaks credentials into terminals or downstream pipelines.
//
// Transformers are designed to be composable and zero-allocation when no
// fields are configured, making them safe to include unconditionally in the
// hot path of the line-processing loop.
package transform
