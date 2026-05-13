package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestMain_NoArgs_ReadsStdin(t *testing.T) {
	bin := buildBinary(t)
	input := `{"time":"2024-01-01T10:00:00Z","level":"info","msg":"hello"}` + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", out)
	}
}

func TestMain_FileArg_FiltersByStart(t *testing.T) {
	bin := buildBinary(t)
	lines := []string{
		`{"time":"2024-01-01T09:00:00Z","msg":"before"}`,
		`{"time":"2024-01-01T11:00:00Z","msg":"after"}`,
	}
	logFile := writeTempLog(t, lines)
	cmd := exec.Command(bin, "-start", "2024-01-01T10:00:00Z", logFile)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(out), "before") {
		t.Errorf("expected 'before' to be filtered out, got: %s", out)
	}
	if !strings.Contains(string(out), "after") {
		t.Errorf("expected 'after' in output, got: %s", out)
	}
}

func TestMain_CountFlag(t *testing.T) {
	bin := buildBinary(t)
	lines := []string{
		`{"time":"2024-01-01T10:00:00Z","msg":"one"}`,
		`{"time":"2024-01-01T11:00:00Z","msg":"two"}`,
	}
	logFile := writeTempLog(t, lines)
	cmd := exec.Command(bin, "-count", logFile)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(string(out))
	if result != "2" {
		t.Errorf("expected count 2, got %q", result)
	}
}

func TestMain_InvalidFormat_Exits(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-format", "xml")
	cmd.Stdin = bytes.NewBufferString("")
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit for invalid format")
	}
}
