package pause

import (
	"testing"
	"time"
)

func newStore(t *testing.T) (*Store, *time.Time) {
	t.Helper()
	s := New()
	now := time.Now()
	s.now = func() time.Time { return now }
	return s, &now
}

func TestIsPaused_NotRegistered(t *testing.T) {
	s, _ := newStore(t)
	if s.IsPaused("job-a") {
		t.Fatal("expected job-a to not be paused")
	}
}

func TestIsPaused_ActiveWindow(t *testing.T) {
	s, _ := newStore(t)
	s.Pause("job-a", 10*time.Minute)
	if !s.IsPaused("job-a") {
		t.Fatal("expected job-a to be paused")
	}
}

func TestIsPaused_ExpiredWindow(t *testing.T) {
	s, now := newStore(t)
	s.Pause("job-a", 5*time.Minute)
	*now = now.Add(10 * time.Minute)
	if s.IsPaused("job-a") {
		t.Fatal("expected pause to have expired")
	}
}

func TestResume_LiftsPause(t *testing.T) {
	s, _ := newStore(t)
	s.Pause("job-a", 10*time.Minute)
	s.Resume("job-a")
	if s.IsPaused("job-a") {
		t.Fatal("expected job-a to no longer be paused after Resume")
	}
}

func TestPausedUntil_ReturnsFalseWhenNotPaused(t *testing.T) {
	s, _ := newStore(t)
	_, ok := s.PausedUntil("job-a")
	if ok {
		t.Fatal("expected ok=false for unregistered job")
	}
}

func TestPausedUntil_ReturnsResumeTime(t *testing.T) {
	s, now := newStore(t)
	d := 30 * time.Minute
	s.Pause("job-a", d)
	want := now.Add(d)
	got, ok := s.PausedUntil("job-a")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if !got.Equal(want) {
		t.Fatalf("PausedUntil = %v, want %v", got, want)
	}
}

func TestPause_ExtendsWindow(t *testing.T) {
	s, now := newStore(t)
	s.Pause("job-a", 10*time.Minute)
	s.Pause("job-a", 30*time.Minute)
	want := now.Add(30 * time.Minute)
	got, _ := s.PausedUntil("job-a")
	if !got.Equal(want) {
		t.Fatalf("extended pause = %v, want %v", got, want)
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	s, now := newStore(t)
	s.Pause("job-a", 5*time.Minute)
	s.Pause("job-b", 20*time.Minute)
	*now = now.Add(10 * time.Minute)
	s.Evict()
	if s.IsPaused("job-a") {
		t.Error("job-a should have been evicted")
	}
	if !s.IsPaused("job-b") {
		t.Error("job-b should still be paused")
	}
}
