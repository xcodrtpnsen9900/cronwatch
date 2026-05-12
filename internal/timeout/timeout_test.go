package timeout_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/timeout"
)

func newTracker() *timeout.Tracker {
	return timeout.New()
}

func TestMark_AddsEntry(t *testing.T) {
	tr := newTracker()
	deadline := time.Now().Add(-time.Minute)
	tr.Mark("job-a", deadline)

	if !tr.IsTimedOut("job-a") {
		t.Fatal("expected job-a to be timed out")
	}
}

func TestMark_Idempotent(t *testing.T) {
	tr := newTracker()
	deadline := time.Now().Add(-time.Minute)
	tr.Mark("job-a", deadline)
	tr.Mark("job-a", deadline.Add(-time.Hour)) // second call should be ignored

	all := tr.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if !all[0].Deadline.Equal(deadline) {
		t.Fatal("deadline should not have changed on second Mark")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	tr := newTracker()
	tr.Mark("job-b", time.Now())
	tr.Clear("job-b")

	if tr.IsTimedOut("job-b") {
		t.Fatal("expected job-b to be cleared")
	}
}

func TestClear_UnknownJob_NoError(t *testing.T) {
	tr := newTracker()
	// should not panic
	tr.Clear("nonexistent")
}

func TestIsTimedOut_NotRegistered(t *testing.T) {
	tr := newTracker()
	if tr.IsTimedOut("missing") {
		t.Fatal("expected missing job to not be timed out")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	tr := newTracker()
	tr.Mark("job-1", time.Now())
	tr.Mark("job-2", time.Now())

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_EmptyTracker(t *testing.T) {
	tr := newTracker()
	if got := tr.All(); len(got) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(got))
	}
}
