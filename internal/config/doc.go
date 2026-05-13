// Package config provides runtime configuration for logslice, including
// flag parsing and validation.
//
// Usage:
//
//	cfg, args, err := config.ParseFlags(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Config fields map directly to CLI flags:
//
//	-start       lower timestamp bound (inclusive)
//	-end         upper timestamp bound (inclusive)
//	-format      Go time layout or named format for timestamp parsing
//	-output      output format: "raw" (default) or "json"
//	-field       key=value field filter; may be repeated
//	-count       emit match count only, suppress line output
//	-buf         internal line-buffer size in bytes
package config
