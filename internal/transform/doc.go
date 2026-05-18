// Package transform provides composable, line-level transformations for
// structured log output. Transformations operate on raw log lines (strings)
// and are applied after filtering but before writing to the output.
//
// Available transformations:
//
//	- Redactor  – replaces sensitive field values with a placeholder.
//	- Truncator – truncates long field values or entire lines to a max length.
//	- Renamer   – renames fields, supporting key=value and JSON-style formats.
//
// Transformations can be chained by applying them sequentially to a line.
package transform
