package suppress

import (
	"testing"
	"time"
)

func newSuppressor(now time.Time) *Suppressor {
	s := New()
	s.now = func() time.Time { return now }
	return s
}

func TestIsSuppressed_NotRegistered(t *testing.T) {
	s := newSuppressor(time.Now())
	if s.IsSuppressed("job1") {
		t.Fatal("expected job1 to not be suppressed")
	}
}

func TestIsSuppressed_ActiveWindow(t *testing.T) {
	base := time.Now()
	s := newSuppressor(base)
	s.Suppress("job1", 10*time.Minute)

	s.now = func() time.Time { return base.Add(5 * time.Minute) }
	if !s.IsSuppressed("job1") {
		t.Fatal("expected job1 to be suppressed within window")
	}
}

func TestIsSuppressed_ExpiredWindow(t *testing.T) {
	base := time.Now()
	s := newSuppressor(base)
	s.Suppress("job1", 5*time.Minute)

	s.now = func() time.Time { return base.Add(10 * time.Minute) }
	if s.IsSuppressed("job1") {
		t.Fatal("expected job1 to not be suppressed after window expires")
	}
}

func TestLift_RemovesWindow(t *testing.T) {
	base := time.Now()
	s := newSuppressor(base)
	s.Suppress("job1", 10*time.Minute)
	s.Lift("job1")

	if s.IsSuppressed("job1") {
		t.Fatal("expected job1 to not be suppressed after lift")
	}
}

func TestActive_ReturnsOnlyLiveWindows(t *testing.T) {
	base := time.Now()
	s := newSuppressor(base)
	s.Suppress("job1", 10*time.Minute)
	s.Suppress("job2", 1*time.Minute)

	// Advance past job2's window.
	s.now = func() time.Time { return base.Add(2 * time.Minute) }

	active := s.Active()
	if _, ok := active["job1"]; !ok {
		t.Error("expected job1 in active windows")
	}
	if _, ok := active["job2"]; ok {
		t.Error("expected job2 to be absent from active windows")
	}
}

func TestSuppress_ReplacesExistingWindow(t *testing.T) {
	base := time.Now()
	s := newSuppressor(base)
	s.Suppress("job1", 1*time.Minute)
	s.Suppress("job1", 30*time.Minute)

	w := s.windows["job1"]
	expected := base.Add(30 * time.Minute)
	if !w.End.Equal(expected) {
		t.Errorf("expected end %v, got %v", expected, w.End)
	}
}
