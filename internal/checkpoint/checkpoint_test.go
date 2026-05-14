package checkpoint_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/checkpoint"
)

func newStore() *checkpoint.Store { return checkpoint.New() }

func TestRecord_And_Get(t *testing.T) {
	s := newStore()
	now := time.Now()
	s.Record("backup", now)

	e, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.LastOK.Equal(now) {
		t.Errorf("LastOK = %v, want %v", e.LastOK, now)
	}
	if e.RunCount != 1 {
		t.Errorf("RunCount = %d, want 1", e.RunCount)
	}
}

func TestRecord_IncrementsRunCount(t *testing.T) {
	s := newStore()
	for i := 0; i < 5; i++ {
		s.Record("sync", time.Now())
	}
	e, _ := s.Get("sync")
	if e.RunCount != 5 {
		t.Errorf("RunCount = %d, want 5", e.RunCount)
	}
}

func TestGet_Missing(t *testing.T) {
	s := newStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected missing entry")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := newStore()
	s.Record("a", time.Now())
	s.Record("b", time.Now())

	all := s.All()
	if len(all) != 2 {
		t.Errorf("len(All) = %d, want 2", len(all))
	}
}

func TestReset_RemovesEntry(t *testing.T) {
	s := newStore()
	s.Record("cleanup", time.Now())
	s.Reset("cleanup")

	_, ok := s.Get("cleanup")
	if ok {
		t.Error("expected entry to be removed after Reset")
	}
}

func TestReset_UnknownJob_NoError(t *testing.T) {
	s := newStore()
	// Should not panic
	s.Reset("ghost")
}

func TestAll_IndependentOfInternalMap(t *testing.T) {
	s := newStore()
	s.Record("job", time.Now())

	all := s.All()
	all[0].RunCount = 999 // mutate snapshot

	e, _ := s.Get("job")
	if e.RunCount == 999 {
		t.Error("All() snapshot should not alias internal state")
	}
}
