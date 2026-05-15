package incident

import (
	"testing"
	"time"
)

func newStore() *Store {
	s := New()
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return fixed }
	return s
}

func TestOpen_CreatesIncident(t *testing.T) {
	s := newStore()
	inc := s.Open("backup")
	if inc.Job != "backup" {
		t.Fatalf("expected job=backup, got %s", inc.Job)
	}
	if inc.Status != StatusOpen {
		t.Fatalf("expected open, got %s", inc.Status)
	}
	if inc.Failures != 1 {
		t.Fatalf("expected 1 failure, got %d", inc.Failures)
	}
	if inc.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestOpen_IncrementsFailures(t *testing.T) {
	s := newStore()
	s.Open("backup")
	inc := s.Open("backup")
	if inc.Failures != 2 {
		t.Fatalf("expected 2 failures, got %d", inc.Failures)
	}
}

func TestOpen_SameIDWhileOpen(t *testing.T) {
	s := newStore()
	a := s.Open("backup")
	b := s.Open("backup")
	if a.ID != b.ID {
		t.Fatalf("expected same incident ID, got %s vs %s", a.ID, b.ID)
	}
}

func TestResolve_MarksResolved(t *testing.T) {
	s := newStore()
	s.Open("backup")
	inc, ok := s.Resolve("backup")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if inc.Status != StatusResolved {
		t.Fatalf("expected resolved, got %s", inc.Status)
	}
	if inc.ResolvedAt == nil {
		t.Fatal("expected ResolvedAt to be set")
	}
}

func TestResolve_NoOpenIncident(t *testing.T) {
	s := newStore()
	_, ok := s.Resolve("backup")
	if ok {
		t.Fatal("expected ok=false for missing incident")
	}
}

func TestOpen_AfterResolve_CreatesNewIncident(t *testing.T) {
	s := newStore()
	a := s.Open("backup")
	s.Resolve("backup")
	b := s.Open("backup")
	if a.ID == b.ID {
		t.Fatal("expected a new incident ID after resolution")
	}
	if b.Failures != 1 {
		t.Fatalf("expected failures reset to 1, got %d", b.Failures)
	}
}

func TestGet_ReturnsCurrentIncident(t *testing.T) {
	s := newStore()
	s.Open("sync")
	inc, ok := s.Get("sync")
	if !ok || inc.Job != "sync" {
		t.Fatal("expected incident for sync")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	s := newStore()
	s.Open("a")
	s.Open("b")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 incidents, got %d", len(all))
	}
}
