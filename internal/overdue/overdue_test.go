package overdue_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/overdue"
)

func newTracker() *overdue.Tracker { return overdue.New() }

func TestMark_AddsEntry(t *testing.T) {
	tr := newTracker()
	expected := time.Now().Add(-5 * time.Minute)
	tr.Mark("backup", expected)

	if !tr.IsOverdue("backup") {
		t.Fatal("expected 'backup' to be overdue")
	}
}

func TestMark_Idempotent(t *testing.T) {
	tr := newTracker()
	expected := time.Now().Add(-5 * time.Minute)
	tr.Mark("backup", expected)

	all := tr.All()
	if len(all) != 1 {
		t.Fatalf("want 1 entry, got %d", len(all))
	}
	first := all[0].Detected

	tr.Mark("backup", expected.Add(-time.Hour))
	all = tr.All()
	if !all[0].Detected.Equal(first) {
		t.Error("second Mark should not overwrite detection timestamp")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	tr := newTracker()
	tr.Mark("sync", time.Now().Add(-time.Minute))
	tr.Clear("sync")

	if tr.IsOverdue("sync") {
		t.Fatal("expected 'sync' to be cleared")
	}
}

func TestClear_UnknownJob_NoError(t *testing.T) {
	tr := newTracker()
	tr.Clear("nonexistent") // must not panic
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	tr := newTracker()
	tr.Mark("jobA", time.Now().Add(-time.Minute))
	tr.Mark("jobB", time.Now().Add(-2*time.Minute))

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("want 2 entries, got %d", len(all))
	}
}

func TestAll_ReturnsSnapshot_Isolation(t *testing.T) {
	// Mutating the returned slice must not affect the tracker's internal state.
	tr := newTracker()
	tr.Mark("jobA", time.Now().Add(-time.Minute))

	snap := tr.All()
	snap[0].Job = "tampered"

	all := tr.All()
	if all[0].Job != "jobA" {
		t.Error("All() should return an independent snapshot, not a reference to internal state")
	}
}

func TestIsOverdue_False_ForUnknown(t *testing.T) {
	tr := newTracker()
	if tr.IsOverdue("ghost") {
		t.Fatal("unknown job should not be overdue")
	}
}
