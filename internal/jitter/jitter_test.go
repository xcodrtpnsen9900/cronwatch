package jitter_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/jitter"
)

// fixedSource always returns the same value so tests are deterministic.
type fixedSource struct{ val int64 }

func (f *fixedSource) Int63n(_ int64) int64 { return f.val }

func TestJitterWith_AddsOffset(t *testing.T) {
	base := 10 * time.Second
	max := 5 * time.Second
	src := &fixedSource{val: 3_000_000_000} // 3 s in nanoseconds

	got := jitter.JitterWith(base, max, src)
	want := 13 * time.Second

	if got != want {
		t.Fatalf("JitterWith: got %v, want %v", got, want)
	}
}

func TestJitterWith_ZeroMaxReturnsBase(t *testing.T) {
	base := 7 * time.Second
	src := &fixedSource{val: 999}

	got := jitter.JitterWith(base, 0, src)
	if got != base {
		t.Fatalf("expected base %v unchanged, got %v", base, got)
	}
}

func TestJitterWith_NegativeMaxReturnsBase(t *testing.T) {
	base := 4 * time.Second
	src := &fixedSource{val: 1}

	got := jitter.JitterWith(base, -1*time.Second, src)
	if got != base {
		t.Fatalf("expected base %v unchanged, got %v", base, got)
	}
}

func TestJitter_ResultInRange(t *testing.T) {
	base := 1 * time.Second
	max := 500 * time.Millisecond

	for i := 0; i < 200; i++ {
		got := jitter.Jitter(base, max)
		if got < base || got >= base+max {
			t.Fatalf("iteration %d: result %v out of range [%v, %v)", i, got, base, base+max)
		}
	}
}

func TestFull_ResultInRange(t *testing.T) {
	max := 2 * time.Second

	for i := 0; i < 200; i++ {
		got := jitter.Full(max)
		if got < 0 || got >= max {
			t.Fatalf("iteration %d: Full result %v out of range [0, %v)", i, got, max)
		}
	}
}

func TestFull_ZeroMaxReturnsZero(t *testing.T) {
	if got := jitter.Full(0); got != 0 {
		t.Fatalf("Full(0): got %v, want 0", got)
	}
}
