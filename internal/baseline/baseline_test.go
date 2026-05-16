package baseline

import (
	"testing"
	"time"
)

func newTracker(t *testing.T) *Tracker {
	t.Helper()
	tr, err := New(100, 2.0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestNew_InvalidMaxPer(t *testing.T) {
	_, err := New(0, 2.0)
	if err == nil {
		t.Fatal("expected error for maxPer=0")
	}
}

func TestNew_InvalidSigma(t *testing.T) {
	_, err := New(10, 0)
	if err == nil {
		t.Fatal("expected error for sigma=0")
	}
}

func TestIsAnomaly_InsufficientData(t *testing.T) {
	tr := newTracker(t)
	tr.Record("job1", 500*time.Millisecond)
	// Only 1 sample — should never flag.
	if tr.IsAnomaly("job1", 10*time.Second) {
		t.Fatal("expected false with only 1 sample")
	}
}

func TestIsAnomaly_NormalDuration(t *testing.T) {
	tr := newTracker(t)
	for i := 0; i < 10; i++ {
		tr.Record("job1", 100*time.Millisecond)
	}
	if tr.IsAnomaly("job1", 110*time.Millisecond) {
		t.Fatal("expected false for duration near mean")
	}
}

func TestIsAnomaly_ExtremeOutlier(t *testing.T) {
	tr := newTracker(t)
	for i := 0; i < 20; i++ {
		tr.Record("job1", 100*time.Millisecond)
	}
	// 10 seconds is far beyond mean+2σ when mean≈100ms, σ≈0.
	if !tr.IsAnomaly("job1", 10*time.Second) {
		t.Fatal("expected true for extreme outlier")
	}
}

func TestSnapshot_UnknownJob(t *testing.T) {
	tr := newTracker(t)
	s := tr.Snapshot("missing")
	if s.Samples != 0 {
		t.Fatalf("expected 0 samples, got %d", s.Samples)
	}
}

func TestSnapshot_AccumulatesSamples(t *testing.T) {
	tr := newTracker(t)
	tr.Record("job2", 200*time.Millisecond)
	tr.Record("job2", 400*time.Millisecond)
	s := tr.Snapshot("job2")
	if s.Samples != 2 {
		t.Fatalf("expected 2 samples, got %d", s.Samples)
	}
	if s.MeanMs != 300 {
		t.Fatalf("expected mean 300ms, got %v", s.MeanMs)
	}
}

func TestRecord_PrunesOldSamples(t *testing.T) {
	tr, _ := New(5, 2.0)
	for i := 0; i < 10; i++ {
		tr.Record("job3", time.Duration(i)*time.Millisecond)
	}
	s := tr.Snapshot("job3")
	if s.Samples != 5 {
		t.Fatalf("expected 5 samples after pruning, got %d", s.Samples)
	}
}
