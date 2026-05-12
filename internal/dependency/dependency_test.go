package dependency

import (
	"testing"
)

func newStore() *Store { return New() }

func TestReady_NoDependencies(t *testing.T) {
	s := newStore()
	ok, err := s.Ready("jobA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ready=true for job with no dependencies")
	}
}

func TestReady_DependencyNeverRun(t *testing.T) {
	s := newStore()
	s.Declare("jobB", []string{"jobA"})
	ok, err := s.Ready("jobB")
	if err == nil {
		t.Fatal("expected error for unrecorded dependency")
	}
	if ok {
		t.Fatal("expected ready=false")
	}
}

func TestReady_DependencyFailed(t *testing.T) {
	s := newStore()
	s.Declare("jobB", []string{"jobA"})
	s.Record("jobA", StateFailure)
	ok, err := s.Ready("jobB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected ready=false when dependency failed")
	}
}

func TestReady_DependencySucceeded(t *testing.T) {
	s := newStore()
	s.Declare("jobB", []string{"jobA"})
	s.Record("jobA", StateSuccess)
	ok, err := s.Ready("jobB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ready=true when dependency succeeded")
	}
}

func TestReady_MultipleDependencies_AllMustSucceed(t *testing.T) {
	s := newStore()
	s.Declare("jobC", []string{"jobA", "jobB"})
	s.Record("jobA", StateSuccess)
	s.Record("jobB", StateFailure)
	ok, err := s.Ready("jobC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected ready=false when one dependency failed")
	}
}

func TestDeclare_Replaces(t *testing.T) {
	s := newStore()
	s.Declare("jobB", []string{"jobA"})
	s.Declare("jobB", []string{}) // replace with no deps
	ok, err := s.Ready("jobB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ready=true after replacing deps with empty list")
	}
}

func TestSnapshot_ReturnsAllEntries(t *testing.T) {
	s := newStore()
	s.Record("jobA", StateSuccess)
	s.Record("jobB", StateFailure)
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}
