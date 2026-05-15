package skiplist

import (
	"testing"
	"time"
)

func newStore(now func() time.Time) *Store {
	s := New()
	s.now = now
	return s
}

func TestIsSkipped_NotRegistered(t *testing.T) {
	s := New()
	if s.IsSkipped("backup") {
		t.Fatal("expected false for unknown job")
	}
}

func TestIsSkipped_ActiveWindow(t *testing.T) {
	now := time.Now()
	s := newStore(func() time.Time { return now })
	_ = s.Skip("backup", "maintenance", now.Add(time.Hour))
	if !s.IsSkipped("backup") {
		t.Fatal("expected job to be skipped")
	}
}

func TestIsSkipped_ExpiredWindow(t *testing.T) {
	now := time.Now()
	s := newStore(func() time.Time { return now })
	_ = s.Skip("backup", "maintenance", now.Add(-time.Second))
	if s.IsSkipped("backup") {
		t.Fatal("expected false for expired skip window")
	}
}

func TestLift_RemovesWindow(t *testing.T) {
	now := time.Now()
	s := newStore(func() time.Time { return now })
	_ = s.Skip("backup", "planned", now.Add(time.Hour))
	s.Lift("backup")
	if s.IsSkipped("backup") {
		t.Fatal("expected skip to be lifted")
	}
}

func TestSkip_EmptyJob_Error(t *testing.T) {
	s := New()
	if err := s.Skip("", "reason", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected error for empty job")
	}
}

func TestSkip_ZeroUntil_Error(t *testing.T) {
	s := New()
	if err := s.Skip("backup", "reason", time.Time{}); err == nil {
		t.Fatal("expected error for zero until time")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	now := time.Now()
	s := newStore(func() time.Time { return now })
	_ = s.Skip("job-a", "r", now.Add(time.Hour))
	_ = s.Skip("job-b", "r", now.Add(time.Hour))
	entries := s.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_IndependentJobs(t *testing.T) {
	now := time.Now()
	s := newStore(func() time.Time { return now })
	_ = s.Skip("job-a", "reason-a", now.Add(time.Hour))
	_ = s.Skip("job-b", "reason-b", now.Add(2*time.Hour))
	if !s.IsSkipped("job-a") || !s.IsSkipped("job-b") {
		t.Fatal("both jobs should be independently skipped")
	}
	s.Lift("job-a")
	if s.IsSkipped("job-a") {
		t.Fatal("job-a should no longer be skipped")
	}
	if !s.IsSkipped("job-b") {
		t.Fatal("job-b should still be skipped")
	}
}
