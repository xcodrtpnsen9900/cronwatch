package replay

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func newStore() *Store { return New(3) }

func TestEnqueue_AppearsInAll(t *testing.T) {
	s := newStore()
	s.Enqueue("backup", "missed", t0)

	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if all[0].JobName != "backup" || all[0].Reason != "missed" {
		t.Errorf("unexpected entry: %+v", all[0])
	}
}

func TestEnqueue_EvictsOldestWhenFull(t *testing.T) {
	s := newStore() // maxPer = 3
	for i := 0; i < 5; i++ {
		s.Enqueue("job", "failed", t0.Add(time.Duration(i)*time.Minute))
	}

	s.mu.Lock()
	got := len(s.entries["job"])
	s.mu.Unlock()

	if got != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", got)
	}
}

func TestDrain_RemovesEntries(t *testing.T) {
	s := newStore()
	s.Enqueue("sync", "missed", t0)
	s.Enqueue("sync", "missed", t0.Add(time.Minute))

	drained := s.Drain("sync")
	if len(drained) != 2 {
		t.Fatalf("expected 2 drained entries, got %d", len(drained))
	}
	if s.Len() != 0 {
		t.Errorf("expected store to be empty after drain")
	}
}

func TestDrain_UnknownJob_ReturnsNil(t *testing.T) {
	s := newStore()
	out := s.Drain("nonexistent")
	if out != nil {
		t.Errorf("expected nil for unknown job, got %v", out)
	}
}

func TestLen_MultipleJobs(t *testing.T) {
	s := newStore()
	s.Enqueue("a", "missed", t0)
	s.Enqueue("a", "missed", t0)
	s.Enqueue("b", "failed", t0)

	if s.Len() != 3 {
		t.Errorf("expected Len 3, got %d", s.Len())
	}
}

func TestNew_DefaultMaxPer(t *testing.T) {
	s := New(0)
	if s.maxPer != 10 {
		t.Errorf("expected default maxPer 10, got %d", s.maxPer)
	}
}
