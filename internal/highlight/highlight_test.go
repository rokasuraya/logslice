package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestNew_ModeOff(t *testing.T) {
	h := highlight.New(highlight.ModeOff, true)
	if h.Enabled() {
		t.Fatal("expected highlighting disabled")
	}
}

func TestNew_ModeOn(t *testing.T) {
	h := highlight.New(highlight.ModeOn, false)
	if !h.Enabled() {
		t.Fatal("expected highlighting enabled")
	}
}

func TestNew_ModeAuto_TTY(t *testing.T) {
	h := highlight.New(highlight.ModeAuto, true)
	if !h.Enabled() {
		t.Fatal("expected highlighting enabled when isTTY=true")
	}
}

func TestNew_ModeAuto_NoTTY(t *testing.T) {
	h := highlight.New(highlight.ModeAuto, false)
	if h.Enabled() {
		t.Fatal("expected highlighting disabled when isTTY=false")
	}
}

func TestLine_Disabled(t *testing.T) {
	h := highlight.New(highlight.ModeOff, false)
	line := "level=info msg=hello"
	if got := h.Line(line, highlight.Green); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestLine_Enabled(t *testing.T) {
	h := highlight.New(highlight.ModeOn, false)
	line := "level=info msg=hello"
	got := h.Line(line, highlight.Green)
	if !strings.Contains(got, highlight.Green) {
		t.Error("expected ANSI green code in output")
	}
	if !strings.Contains(got, highlight.Reset) {
		t.Error("expected ANSI reset code in output")
	}
	if !strings.Contains(got, line) {
		t.Error("expected original line content preserved")
	}
}

func TestField_Disabled(t *testing.T) {
	h := highlight.New(highlight.ModeOff, false)
	line := "level=info msg=hello"
	if got := h.Field(line, "level"); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestField_Enabled_KeyFound(t *testing.T) {
	h := highlight.New(highlight.ModeOn, false)
	line := "level=info msg=hello"
	got := h.Field(line, "level")
	if !strings.Contains(got, highlight.Bold) {
		t.Error("expected bold code in highlighted field output")
	}
	if !strings.Contains(got, "level=info") {
		t.Error("expected field content preserved")
	}
}

func TestField_Enabled_KeyNotFound(t *testing.T) {
	h := highlight.New(highlight.ModeOn, false)
	line := "level=info msg=hello"
	got := h.Field(line, "missing")
	if got != line {
		t.Errorf("expected unchanged line when key not found, got %q", got)
	}
}

func TestField_Enabled_KeyAtEnd(t *testing.T) {
	h := highlight.New(highlight.ModeOn, false)
	line := "ts=2024-01-01 level=error"
	got := h.Field(line, "level")
	if !strings.Contains(got, "level=error") {
		t.Error("expected field at end of line to be highlighted")
	}
}
