package circuitbreaker

import (
	"testing"
	"time"
)

func newBreaker(max int) *Breaker {
	b := New(max, 5*time.Second)
	return b
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := newBreaker(3)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpensAfterMaxFailures(t *testing.T) {
	b := newBreaker(3)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != StateOpen {
		t.Fatalf("expected Open, got %v", b.State())
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestDoesNotOpenBeforeThreshold(t *testing.T) {
	b := newBreaker(3)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatalf("expected Closed after 2 failures with threshold 3")
	}
}

func TestHalfOpenAfterCooldown(t *testing.T) {
	b := newBreaker(1)
	fixed := time.Now()
	b.now = func() time.Time { return fixed }
	b.RecordFailure()

	// still within cooldown
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen during cooldown")
	}

	// advance past cooldown
	b.now = func() time.Time { return fixed.Add(10 * time.Second) }
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", b.State())
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := newBreaker(1)
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected Open")
	}
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected Closed after success, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestRecordSuccess_ResetsFailureCount(t *testing.T) {
	b := newBreaker(3)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()
	// only 1 failure after reset, should still be closed
	if b.State() != StateClosed {
		t.Fatalf("expected Closed, got %v", b.State())
	}
}
