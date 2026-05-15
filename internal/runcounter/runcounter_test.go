package runcounter_test

import (
	"testing"

	"github.com/cronwatch/cronwatch/internal/runcounter"
)

func newStore() *runcounter.Store {
	return runcounter.New()
}

func TestIncrement_StartsAtOne(t *testing.T) {
	s := newStore()
	if got := s.Increment("job-a"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestIncrement_Accumulates(t *testing.T) {
	s := newStore()
	s.Increment("job-a")
	s.Increment("job-a")
	if got := s.Increment("job-a"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestGet_UnknownJobReturnsZero(t *testing.T) {
	s := newStore()
	if got := s.Get("missing"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestGet_ReturnsCurrentValue(t *testing.T) {
	s := newStore()
	s.Increment("job-b")
	s.Increment("job-b")
	if got := s.Get("job-b"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestReset_ZeroesCounter(t *testing.T) {
	s := newStore()
	s.Increment("job-c")
	s.Reset("job-c")
	if got := s.Get("job-c"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestReset_UnknownJob_NoError(t *testing.T) {
	s := newStore()
	s.Reset("nonexistent") // must not panic
}

func TestSnapshot_ContainsAllJobs(t *testing.T) {
	s := newStore()
	s.Increment("alpha")
	s.Increment("alpha")
	s.Increment("beta")

	snap := s.Snapshot()
	if snap["alpha"] != 2 {
		t.Errorf("alpha: expected 2, got %d", snap["alpha"])
	}
	if snap["beta"] != 1 {
		t.Errorf("beta: expected 1, got %d", snap["beta"])
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	s := newStore()
	s.Increment("job-d")
	snap := s.Snapshot()
	snap["job-d"] = 999
	if got := s.Get("job-d"); got != 1 {
		t.Fatalf("snapshot mutation affected store: got %d", got)
	}
}

func TestJobs_ListsKnownJobs(t *testing.T) {
	s := newStore()
	s.Increment("x")
	s.Increment("y")

	jobs := s.Jobs()
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}
