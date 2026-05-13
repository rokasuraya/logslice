package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.log")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestNewFileReader_ReadsLines(t *testing.T) {
	content := "line one\nline two\nline three\n"
	path := writeTempFile(t, content)

	r, err := NewFileReader(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()

	expected := []string{"line one", "line two", "line three"}
	var got []string
	for r.Next() {
		got = append(got, r.Line())
	}
	if err := r.Err(); err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, line := range expected {
		if got[i] != line {
			t.Errorf("line %d: expected %q, got %q", i, line, got[i])
		}
	}
}

func TestNewFileReader_NotFound(t *testing.T) {
	_, err := NewFileReader("/nonexistent/path/to/file.log", nil)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestNewFileReader_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")

	r, err := NewFileReader(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()

	if r.Next() {
		t.Errorf("expected no lines for empty file, got: %q", r.Line())
	}
}

func TestNewFileReader_CustomBufferSize(t *testing.T) {
	content := "short line\n"
	path := writeTempFile(t, content)

	r, err := NewFileReader(path, &Options{BufferSize: 4096})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()

	if !r.Next() {
		t.Fatal("expected at least one line")
	}
	if r.Line() != "short line" {
		t.Errorf("expected %q, got %q", "short line", r.Line())
	}
}
