package deadletter

import (
	"fmt"
	"testing"
)

func newStore(max int) *Store { return New(max) }

func TestAdd_AppearsInAll(t *testing.T) {
	s := newStore(10)
	s.Add("job1", []byte(`{}`), "timeout")
	s.Add("job2", []byte(`{}`), "500")

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[0].JobName != "job1" {
		t.Errorf("expected job1 first, got %s", all[0].JobName)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	s := newStore(3)
	for i := 0; i < 4; i++ {
		s.Add(fmt.Sprintf("job%d", i), nil, "err")
	}
	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	if all[0].JobName != "job1" {
		t.Errorf("expected job1 after eviction, got %s", all[0].JobName)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	s := newStore(10)
	e := s.Add("job1", nil, "err")

	ok := s.Remove(e.ID)
	if !ok {
		t.Fatal("expected Remove to return true")
	}
	if s.Size() != 0 {
		t.Errorf("expected size 0, got %d", s.Size())
	}
}

func TestRemove_UnknownID(t *testing.T) {
	s := newStore(10)
	ok := s.Remove("nonexistent")
	if ok {
		t.Error("expected Remove to return false for unknown ID")
	}
}

func TestDefaultMaxSize(t *testing.T) {
	s := New(0)
	for i := 0; i < 105; i++ {
		s.Add("job", nil, "err")
	}
	if s.Size() != 100 {
		t.Errorf("expected size 100, got %d", s.Size())
	}
}

func TestEntry_AttemptsSetToOne(t *testing.T) {
	s := newStore(10)
	e := s.Add("myjob", []byte("payload"), "connect refused")
	if e.Attempts != 1 {
		t.Errorf("expected Attempts=1, got %d", e.Attempts)
	}
	if e.Reason != "connect refused" {
		t.Errorf("unexpected reason: %s", e.Reason)
	}
}
