package reader

import (
	"bufio"
	"io"
	"os"
)

// LineReader reads lines from a file or stdin with optional byte-range support.
type LineReader struct {
	scanner *bufio.Scanner
	closer  io.Closer
}

// Options configures the LineReader.
type Options struct {
	// BufferSize sets the max line buffer size (default: 1MB)
	BufferSize int
}

const defaultBufferSize = 1024 * 1024 // 1MB

// NewFileReader opens a file and returns a LineReader for it.
func NewFileReader(path string, opts *Options) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return newReader(f, f, opts), nil
}

// NewStdinReader returns a LineReader that reads from stdin.
func NewStdinReader(opts *Options) *LineReader {
	return newReader(os.Stdin, nil, opts)
}

func newReader(r io.Reader, closer io.Closer, opts *Options) *LineReader {
	bufSize := defaultBufferSize
	if opts != nil && opts.BufferSize > 0 {
		bufSize = opts.BufferSize
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, bufSize), bufSize)

	return &LineReader{
		scanner: scanner,
		closer:  closer,
	}
}

// Next advances to the next line. Returns false when done or on error.
func (lr *LineReader) Next() bool {
	return lr.scanner.Scan()
}

// Line returns the current line text.
func (lr *LineReader) Line() string {
	return lr.scanner.Text()
}

// Err returns any scanning error (excluding io.EOF).
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

// Close releases any underlying file resources.
func (lr *LineReader) Close() error {
	if lr.closer != nil {
		return lr.closer.Close()
	}
	return nil
}
