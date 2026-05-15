package runbook

import (
	"testing"
)

func newStore() *Store { return New() }

func TestSet_And_Get(t *testing.T) {
	s := newStore()
	if err := s.Set("backup", "https://wiki.example.com/backup", "Backup job"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.URL != "https://wiki.example.com/backup" {
		t.Errorf("url mismatch: got %q", e.URL)
	}
	if e.Summary != "Backup job" {
		t.Errorf("summary mismatch: got %q", e.Summary)
	}
}

func TestGet_Missing(t *testing.T) {
	s := newStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSet_EmptyJob_Error(t *testing.T) {
	s := newStore()
	if err := s.Set("", "https://example.com", ""); err == nil {
		t.Fatal("expected error for empty job")
	}
}

func TestSet_EmptyURL_Error(t *testing.T) {
	s := newStore()
	if err := s.Set("myjob", "", ""); err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestSet_Replaces(t *testing.T) {
	s := newStore()
	_ = s.Set("job", "https://old.example.com", "old")
	_ = s.Set("job", "https://new.example.com", "new")
	e, _ := s.Get("job")
	if e.URL != "https://new.example.com" {
		t.Errorf("expected updated url, got %q", e.URL)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	s := newStore()
	_ = s.Set("job", "https://example.com", "")
	s.Remove("job")
	_, ok := s.Get("job")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestRemove_Unknown_NoError(t *testing.T) {
	s := newStore()
	s.Remove("ghost") // must not panic
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := newStore()
	_ = s.Set("a", "https://a.example.com", "")
	_ = s.Set("b", "https://b.example.com", "")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_Empty(t *testing.T) {
	s := newStore()
	all := s.All()
	if len(all) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(all))
	}
}

func TestAll_IsolatedFromStore(t *testing.T) {
	s := newStore()
	_ = s.Set("job", "https://example.com", "desc")
	all := s.All()
	// Mutating the returned map must not affect the store.
	delete(all, "job")
	_, ok := s.Get("job")
	if !ok {
		t.Fatal("store was mutated via All() return value")
	}
}
