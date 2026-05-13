package runlock

import (
	"testing"
	"time"
)

func newStore(ttl time.Duration) *Store {
	s := New(ttl)
	return s
}

func TestAcquire_FirstCallSucceeds(t *testing.T) {
	s := newStore(0)
	if err := s.Acquire("job-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAcquire_SecondCallBlocked(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("job-a")
	if err := s.Acquire("job-a"); err == nil {
		t.Fatal("expected error for duplicate acquire, got nil")
	}
}

func TestRelease_AllowsReacquire(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("job-a")
	s.Release("job-a")
	if err := s.Acquire("job-a"); err != nil {
		t.Fatalf("expected acquire after release to succeed, got: %v", err)
	}
}

func TestRelease_UnknownJob_NoError(t *testing.T) {
	s := newStore(0)
	s.Release("nonexistent") // must not panic
}

func TestIsLocked_True(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("job-b")
	if !s.IsLocked("job-b") {
		t.Fatal("expected job-b to be locked")
	}
}

func TestIsLocked_False_AfterRelease(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("job-b")
	s.Release("job-b")
	if s.IsLocked("job-b") {
		t.Fatal("expected job-b to be unlocked after release")
	}
}

func TestTTL_ExpiresLock(t *testing.T) {
	s := New(50 * time.Millisecond)
	past := time.Now().Add(-100 * time.Millisecond)
	s.nowFunc = func() time.Time { return past }
	_ = s.Acquire("job-c")

	s.nowFunc = time.Now // advance time past TTL
	if s.IsLocked("job-c") {
		t.Fatal("expected lock to have expired")
	}
	if err := s.Acquire("job-c"); err != nil {
		t.Fatalf("expected re-acquire after TTL expiry, got: %v", err)
	}
}

func TestActive_ReturnsLockedJobs(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("job-x")
	_ = s.Acquire("job-y")

	active := s.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active jobs, got %d", len(active))
	}
}

func TestActive_ExcludesExpired(t *testing.T) {
	s := New(50 * time.Millisecond)
	past := time.Now().Add(-100 * time.Millisecond)
	s.nowFunc = func() time.Time { return past }
	_ = s.Acquire("old-job")
	s.nowFunc = time.Now

	if len(s.Active()) != 0 {
		t.Fatal("expected no active jobs after TTL expiry")
	}
}
