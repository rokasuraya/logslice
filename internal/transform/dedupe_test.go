package transform

import (
	"testing"
)

func TestNewDeduplicator_InvalidWindow(t *testing.T) {
	_, err := NewDeduplicator(-1)
	if err == nil {
		t.Fatal("expected error for negative window size")
	}
}

func TestNewDeduplicator_ValidWindow(t *testing.T) {
	d, err := NewDeduplicator(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Len() != 0 {
		t.Errorf("expected empty deduplicator, got len=%d", d.Len())
	}
}

func TestDeduplicator_FirstOccurrence_NotDuplicate(t *testing.T) {
	d, _ := NewDeduplicator(0)
	if d.IsDuplicate("hello world") {
		t.Error("first occurrence should not be a duplicate")
	}
}

func TestDeduplicator_SecondOccurrence_IsDuplicate(t *testing.T) {
	d, _ := NewDeduplicator(0)
	d.IsDuplicate("hello world")
	if !d.IsDuplicate("hello world") {
		t.Error("second occurrence should be a duplicate")
	}
}

func TestDeduplicator_DifferentLines_NotDuplicate(t *testing.T) {
	d, _ := NewDeduplicator(0)
	d.IsDuplicate("line one")
	if d.IsDuplicate("line two") {
		t.Error("different lines should not be duplicates")
	}
}

func TestDeduplicator_WindowEviction(t *testing.T) {
	d, _ := NewDeduplicator(2)
	d.IsDuplicate("line A") // slot 1
	d.IsDuplicate("line B") // slot 2
	d.IsDuplicate("line C") // evicts "line A"

	// "line A" should have been evicted, so it should not be a duplicate
	if d.IsDuplicate("line A") {
		t.Error("evicted line should not be considered a duplicate")
	}
	if d.Len() != 2 {
		t.Errorf("expected window size 2, got %d", d.Len())
	}
}

func TestDeduplicator_Reset(t *testing.T) {
	d, _ := NewDeduplicator(0)
	d.IsDuplicate("some line")
	d.Reset()
	if d.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", d.Len())
	}
	if d.IsDuplicate("some line") {
		t.Error("line should not be duplicate after reset")
	}
}

func TestDeduplicator_EmptyLine(t *testing.T) {
	d, _ := NewDeduplicator(0)
	if d.IsDuplicate("") {
		t.Error("first empty line should not be a duplicate")
	}
	if !d.IsDuplicate("") {
		t.Error("second empty line should be a duplicate")
	}
}
