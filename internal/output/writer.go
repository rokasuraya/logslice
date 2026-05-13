package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format defines the output format for log lines.
type Format string

const (
	FormatRaw  Format = "raw"
	FormatJSON Format = "json"
)

// Writer wraps an io.Writer with buffering and format support.
type Writer struct {
	w      *bufio.Writer
	format Format
	count  int
}

// NewWriter creates a Writer targeting the given io.Writer.
func NewWriter(w io.Writer, format Format) *Writer {
	if format == "" {
		format = FormatRaw
	}
	return &Writer{
		w:      bufio.NewWriterSize(w, 64*1024),
		format: format,
	}
}

// NewStdoutWriter creates a Writer targeting stdout.
func NewStdoutWriter(format Format) *Writer {
	return NewWriter(os.Stdout, format)
}

// WriteLine writes a single log line according to the configured format.
func (w *Writer) WriteLine(line string) error {
	var err error
	switch w.format {
	case FormatJSON:
		_, err = fmt.Fprintf(w.w, "{\"line\":%q}\n", line)
	default:
		_, err = fmt.Fprintln(w.w, line)
	}
	if err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	w.count++
	return nil
}

// Count returns the number of lines written.
func (w *Writer) Count() int {
	return w.count
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return nil
}
