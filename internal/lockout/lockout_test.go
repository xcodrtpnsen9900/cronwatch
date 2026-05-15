package lockout

import (
	"testing"
	"time"
)

func newStore(t *testing.T, threshold int, window time.Duration) *Store {
	t.Helper()
	s, err := New(threshold, window)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for threshold=0")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(1, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestIsLockedOut_NotRegistered(t *testing.T) {
	s := newStore(t, 3, time.Minute)
	if s.IsLockedOut("job") {
		t.Fatal("expected not locked out")
	}
}

func TestRecordFailure_BelowThreshold(t *testing.T) {
	s := newStore(t, 3, time.Minute)
	s.RecordFailure("job")
	s.RecordFailure("job")
	if s.IsLockedOut("job") {
		t.Fatal("should not be locked out below threshold")
	}
}

func TestRecordFailure_AtThreshold(t *testing.T) {
	s := newStore(t, 3, time.Minute)
	s.RecordFailure("job")
	s.RecordFailure("job")
	s.RecordFailure("job")
	if !s.IsLockedOut("job") {
		t.Fatal("expected lockout at threshold")
	}
}

func TestRecordSuccess_ClearsLockout(t *testing.T) {
	s := newStore(t, 2, time.Minute)
	s.RecordFailure("job")
	s.RecordFailure("job")
	s.RecordSuccess("job")
	if s.IsLockedOut("job") {
		t.Fatal("expected lockout cleared after success")
	}
}

func TestIsLockedOut_ExpiredWindow(t *testing.T) {
	s := newStore(t, 1, time.Millisecond)
	s.RecordFailure("job")
	time.Sleep(5 * time.Millisecond)
	if s.IsLockedOut("job") {
		t.Fatal("expected lockout expired")
	}
}

func TestLift_RemovesLockout(t *testing.T) {
	s := newStore(t, 1, time.Minute)
	s.RecordFailure("job")
	s.Lift("job")
	if s.IsLockedOut("job") {
		t.Fatal("expected lockout lifted")
	}
}

func TestAll_ReturnsActiveOnly(t *testing.T) {
	s := newStore(t, 1, time.Minute)
	s.RecordFailure("a")
	s.RecordFailure("b")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 locked-out jobs, got %d", len(all))
	}
	if _, ok := all["a"]; !ok {
		t.Error("expected job a in All()")
	}
}

func TestAll_IndependentJobs(t *testing.T) {
	s := newStore(t, 2, time.Minute)
	s.RecordFailure("x")
	s.RecordFailure("x")
	s.RecordFailure("y") // below threshold
	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 locked-out job, got %d", len(all))
	}
}
