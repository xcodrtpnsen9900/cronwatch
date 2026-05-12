package retention

import (
	"testing"
	"time"
)

func newStore(maxAge time.Duration) *Store {
	s := New(Policy{MaxAge: maxAge})
	return s
}

func TestAdd_RetainsRecentEntry(t *testing.T) {
	s := newStore(time.Hour)
	s.Add(Entry{Job: "backup", Timestamp: time.Now()})
	if got := len(s.All()); got != 1 {
		t.Fatalf("expected 1 entry, got %d", got)
	}
}

func TestAdd_EvictsExpiredEntry(t *testing.T) {
	s := newStore(time.Hour)
	old := time.Now().Add(-2 * time.Hour)
	s.Add(Entry{Job: "backup", Timestamp: old})
	if got := len(s.All()); got != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", got)
	}
}

func TestEvict_ReturnsRemovedCount(t *testing.T) {
	s := newStore(time.Hour)
	now := time.Now()
	s.mu.Lock()
	s.entries = []Entry{
		{Job: "a", Timestamp: now.Add(-2 * time.Hour)},
		{Job: "b", Timestamp: now.Add(-30 * time.Minute)},
		{Job: "c", Timestamp: now},
	}
	s.mu.Unlock()

	removed := s.Evict()
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if got := len(s.All()); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := newStore(time.Hour)
	s.Add(Entry{Job: "sync", Timestamp: time.Now()})
	a := s.All()
	a[0].Job = "mutated"
	b := s.All()
	if b[0].Job == "mutated" {
		t.Fatal("All() returned a reference to internal slice")
	}
}

func TestDefaultPolicy_MaxAge(t *testing.T) {
	expected := 30 * 24 * time.Hour
	if DefaultPolicy.MaxAge != expected {
		t.Fatalf("expected MaxAge %v, got %v", expected, DefaultPolicy.MaxAge)
	}
}

func TestAdd_MixedExpiry(t *testing.T) {
	s := newStore(time.Hour)
	now := time.Now()
	s.Add(Entry{Job: "old", Timestamp: now.Add(-90 * time.Minute)})
	s.Add(Entry{Job: "new", Timestamp: now})

	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if all[0].Job != "new" {
		t.Fatalf("expected job 'new', got %q", all[0].Job)
	}
}
