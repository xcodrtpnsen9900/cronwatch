package cooldown

import (
	"testing"
	"time"
)

func newStore(d time.Duration) *Store {
	s := New(d)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return base }
	return s
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	s := newStore(5 * time.Minute)
	if !s.Allow("job-a") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SuppressedDuringCooldown(t *testing.T) {
	s := newStore(5 * time.Minute)
	s.Record("job-a")
	if s.Allow("job-a") {
		t.Fatal("expected Allow to return false during cooldown")
	}
}

func TestAllow_PermittedAfterCooldown(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	s := New(5 * time.Minute)
	s.now = func() time.Time { return current }

	s.Record("job-a")
	current = base.Add(6 * time.Minute)

	if !s.Allow("job-a") {
		t.Fatal("expected Allow to return true after cooldown elapsed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	s := newStore(5 * time.Minute)
	s.Record("job-a")

	if s.Allow("job-a") {
		t.Fatal("expected job-a to be suppressed")
	}
	if !s.Allow("job-b") {
		t.Fatal("expected job-b to be permitted independently")
	}
}

func TestReset_AllowsImmediateAlert(t *testing.T) {
	s := newStore(5 * time.Minute)
	s.Record("job-a")
	s.Reset("job-a")

	if !s.Allow("job-a") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestRemaining_ZeroWhenUnknown(t *testing.T) {
	s := newStore(5 * time.Minute)
	if r := s.Remaining("job-x"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown key, got %v", r)
	}
}

func TestRemaining_PositiveDuringCooldown(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	s := New(5 * time.Minute)
	s.now = func() time.Time { return current }

	s.Record("job-a")
	current = base.Add(2 * time.Minute)

	r := s.Remaining("job-a")
	if r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestRemaining_ZeroAfterCooldown(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	s := New(5 * time.Minute)
	s.now = func() time.Time { return current }

	s.Record("job-a")
	current = base.Add(10 * time.Minute)

	if r := s.Remaining("job-a"); r != 0 {
		t.Fatalf("expected 0 remaining after cooldown, got %v", r)
	}
}
