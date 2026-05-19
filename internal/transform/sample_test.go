package transform

import (
	"testing"
)

func TestNewSampler_InvalidRate(t *testing.T) {
	_, err := NewSampler(0)
	if err == nil {
		t.Fatal("expected error for rate=0, got nil")
	}
}

func TestNewSampler_ValidRate(t *testing.T) {
	s, err := NewSampler(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 3 {
		t.Errorf("expected rate 3, got %d", s.Rate())
	}
}

func TestSampler_RateOne_KeepsAll(t *testing.T) {
	s, _ := NewSampler(1)
	for i := 0; i < 5; i++ {
		_, kept := s.Apply("line")
		if !kept {
			t.Errorf("rate=1: expected line %d to be kept", i)
		}
	}
}

func TestSampler_RateThree_KeepsEveryThird(t *testing.T) {
	s, _ := NewSampler(3)
	expected := []bool{true, false, false, true, false, false, true}
	for i, want := range expected {
		_, kept := s.Apply("line")
		if kept != want {
			t.Errorf("index %d: expected kept=%v, got %v", i, want, kept)
		}
	}
}

func TestSampler_Apply_ReturnsLine_WhenKept(t *testing.T) {
	s, _ := NewSampler(2)
	const input = "hello world"
	out, kept := s.Apply(input)
	if !kept {
		t.Fatal("expected first line to be kept")
	}
	if out != input {
		t.Errorf("expected %q, got %q", input, out)
	}
}

func TestSampler_Apply_ReturnsEmpty_WhenDropped(t *testing.T) {
	s, _ := NewSampler(2)
	s.Apply("first") // kept
	out, kept := s.Apply("second")
	if kept {
		t.Fatal("expected second line to be dropped")
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestSampler_Reset_RestartsCounter(t *testing.T) {
	s, _ := NewSampler(3)
	s.Apply("a") // 0 -> kept
	s.Apply("b") // 1 -> dropped
	s.Reset()
	_, kept := s.Apply("c") // should be kept again (counter reset to 0)
	if !kept {
		t.Error("expected line to be kept after reset")
	}
}

func TestSampler_RateTwo_HalfKept(t *testing.T) {
	s, _ := NewSampler(2)
	kept := 0
	for i := 0; i < 100; i++ {
		_, ok := s.Apply("line")
		if ok {
			kept++
		}
	}
	if kept != 50 {
		t.Errorf("expected 50 kept lines, got %d", kept)
	}
}
