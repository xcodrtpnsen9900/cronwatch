package backoff_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/backoff"
)

func TestDefaultPolicy_Fields(t *testing.T) {
	p := backoff.DefaultPolicy()
	if p.InitialInterval != 500*time.Millisecond {
		t.Errorf("expected 500ms initial interval, got %v", p.InitialInterval)
	}
	if p.MaxInterval != 30*time.Second {
		t.Errorf("expected 30s max interval, got %v", p.MaxInterval)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected multiplier 2.0, got %v", p.Multiplier)
	}
}

func TestNext_FirstAttemptEqualsInitial(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     1 * time.Minute,
		Multiplier:      2.0,
		Jitter:          0, // no jitter for deterministic test
	}
	got := p.Next(0)
	if got != 1*time.Second {
		t.Errorf("expected 1s, got %v", got)
	}
}

func TestNext_Doubles(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     1 * time.Minute,
		Multiplier:      2.0,
		Jitter:          0,
	}
	expected := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	for i, want := range expected {
		if got := p.Next(i); got != want {
			t.Errorf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}
	for i := 5; i < 20; i++ {
		if got := p.Next(i); got > p.MaxInterval {
			t.Errorf("attempt %d: %v exceeds max %v", i, got, p.MaxInterval)
		}
	}
}

func TestNext_NegativeAttemptTreatedAsZero(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 2 * time.Second,
		MaxInterval:     1 * time.Minute,
		Multiplier:      2.0,
		Jitter:          0,
	}
	if got := p.Next(-3); got != 2*time.Second {
		t.Errorf("expected 2s for negative attempt, got %v", got)
	}
}

func TestNext_JitterStaysWithinBounds(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     1 * time.Minute,
		Multiplier:      2.0,
		Jitter:          0.5,
	}
	for i := 0; i < 100; i++ {
		got := p.Next(0)
		if got < 0 {
			t.Errorf("negative delay: %v", got)
		}
		if got > 2*time.Second {
			t.Errorf("delay %v exceeds jitter upper bound", got)
		}
	}
}

func TestSequence_Length(t *testing.T) {
	p := backoff.DefaultPolicy()
	seq := p.Sequence(5)
	if len(seq) != 5 {
		t.Errorf("expected 5 delays, got %d", len(seq))
	}
}
