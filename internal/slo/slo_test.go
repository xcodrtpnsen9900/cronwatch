package slo

import (
	"testing"
	"time"
)

func newTracker() *Tracker {
	return New(24*time.Hour, 100)
}

func TestRecord_MetCompliance(t *testing.T) {
	tr := newTracker()
	now := time.Now()
	tr.Record("job1", now, true)
	tr.Record("job1", now, true)
	tr.Record("job1", now, false)

	snap := tr.Snapshot("job1")
	if snap.Total != 3 {
		t.Fatalf("expected 3 total, got %d", snap.Total)
	}
	if snap.Met != 2 {
		t.Fatalf("expected 2 met, got %d", snap.Met)
	}
	want := 2.0 / 3.0 * 100
	if snap.Compliance != want {
		t.Fatalf("expected %.2f compliance, got %.2f", want, snap.Compliance)
	}
}

func TestSnapshot_UnknownJob_ZeroCompliance(t *testing.T) {
	tr := newTracker()
	snap := tr.Snapshot("unknown")
	if snap.Total != 0 || snap.Compliance != 0 {
		t.Fatalf("expected zero snapshot for unknown job, got %+v", snap)
	}
}

func TestRecord_PrunesExpiredEntries(t *testing.T) {
	tr := New(time.Minute, 100)
	old := time.Now().Add(-2 * time.Minute)
	tr.Record("job1", old, true)
	tr.Record("job1", time.Now(), false)

	snap := tr.Snapshot("job1")
	if snap.Total != 1 {
		t.Fatalf("expected 1 entry after pruning, got %d", snap.Total)
	}
	if snap.Met != 0 {
		t.Fatalf("expected 0 met after pruning old entry, got %d", snap.Met)
	}
}

func TestRecord_EnforcesMaxPer(t *testing.T) {
	tr := New(24*time.Hour, 3)
	now := time.Now()
	for i := 0; i < 5; i++ {
		tr.Record("job1", now, true)
	}
	snap := tr.Snapshot("job1")
	if snap.Total != 3 {
		t.Fatalf("expected max 3 entries, got %d", snap.Total)
	}
}

func TestAll_ReturnsAllJobs(t *testing.T) {
	tr := newTracker()
	now := time.Now()
	tr.Record("alpha", now, true)
	tr.Record("beta", now, false)

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(all))
	}
}

func TestAll_EmptyTracker(t *testing.T) {
	tr := newTracker()
	if got := tr.All(); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
