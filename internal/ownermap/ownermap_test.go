package ownermap_test

import (
	"testing"

	"github.com/cronwatch/cronwatch/internal/ownermap"
)

func newStore() *ownermap.Store {
	return ownermap.New()
}

func TestSet_And_Get(t *testing.T) {
	s := newStore()
	o := ownermap.Owner{Name: "Alice", Email: "alice@example.com", Team: "platform"}
	if err := s.Set("backup", o); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected owner to be found")
	}
	if got != o {
		t.Errorf("got %+v, want %+v", got, o)
	}
}

func TestGet_Missing(t *testing.T) {
	s := newStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Fatal("expected no owner for unknown job")
	}
}

func TestSet_EmptyJob_Error(t *testing.T) {
	s := newStore()
	err := s.Set("", ownermap.Owner{Name: "Bob"})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestSet_EmptyOwnerName_Error(t *testing.T) {
	s := newStore()
	err := s.Set("backup", ownermap.Owner{Email: "bob@example.com"})
	if err == nil {
		t.Fatal("expected error for empty owner name")
	}
}

func TestSet_Replaces(t *testing.T) {
	s := newStore()
	_ = s.Set("backup", ownermap.Owner{Name: "Alice", Team: "platform"})
	_ = s.Set("backup", ownermap.Owner{Name: "Bob", Team: "infra"})
	got, _ := s.Get("backup")
	if got.Name != "Bob" {
		t.Errorf("expected Bob, got %s", got.Name)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	s := newStore()
	_ = s.Set("backup", ownermap.Owner{Name: "Alice"})
	s.Remove("backup")
	_, ok := s.Get("backup")
	if ok {
		t.Fatal("expected owner to be removed")
	}
}

func TestRemove_UnknownJob_NoError(t *testing.T) {
	s := newStore()
	s.Remove("ghost") // must not panic
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := newStore()
	_ = s.Set("job-a", ownermap.Owner{Name: "Alice"})
	_ = s.Set("job-b", ownermap.Owner{Name: "Bob"})
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the snapshot must not affect the store.
	delete(all, "job-a")
	if _, ok := s.Get("job-a"); !ok {
		t.Fatal("store was mutated via snapshot")
	}
}
