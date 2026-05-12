package budget

import (
	"testing"
)

func newTracker(t *testing.T, threshold float64) *Tracker {
	t.Helper()
	tr, err := New(threshold)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestNew_InvalidThreshold(t *testing.T) {
	for _, v := range []float64{0, -0.1, 1.1} {
		_, err := New(v)
		if err == nil {
			t.Errorf("expected error for threshold %v", v)
		}
	}
}

func TestNew_ValidThreshold(t *testing.T) {
	if _, err := New(0.05); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSnapshot_NoRuns(t *testing.T) {
	tr := newTracker(t, 0.05)
	e := tr.Snapshot("job1")
	if e.Total != 0 || e.Failed != 0 {
		t.Errorf("expected zero counts, got total=%d failed=%d", e.Total, e.Failed)
	}
	if e.Exhausted {
		t.Error("expected budget not exhausted with no runs")
	}
	if e.Remaining != 1.0 {
		t.Errorf("expected remaining=1.0, got %v", e.Remaining)
	}
}

func TestBudgetNotExhausted_BelowThreshold(t *testing.T) {
	tr := newTracker(t, 0.10)
	for i := 0; i < 95; i++ {
		tr.RecordSuccess("job1")
	}
	for i := 0; i < 5; i++ {
		tr.RecordFailure("job1")
	}
	e := tr.Snapshot("job1")
	if e.Exhausted {
		t.Errorf("expected budget not exhausted: ratio=5/100=0.05 < threshold=0.10")
	}
	if e.Remaining <= 0 {
		t.Errorf("expected positive remaining, got %v", e.Remaining)
	}
}

func TestBudgetExhausted_AboveThreshold(t *testing.T) {
	tr := newTracker(t, 0.05)
	for i := 0; i < 9; i++ {
		tr.RecordSuccess("job1")
	}
	tr.RecordFailure("job1") // 1/10 = 10% > 5%
	e := tr.Snapshot("job1")
	if !e.Exhausted {
		t.Error("expected budget exhausted")
	}
	if e.Remaining != 0 {
		t.Errorf("expected remaining=0 when exhausted, got %v", e.Remaining)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	tr := newTracker(t, 0.05)
	tr.RecordFailure("job1")
	tr.Reset("job1")
	e := tr.Snapshot("job1")
	if e.Total != 0 || e.Failed != 0 {
		t.Errorf("expected zero after reset, got total=%d failed=%d", e.Total, e.Failed)
	}
}

func TestIndependentJobs(t *testing.T) {
	tr := newTracker(t, 0.05)
	tr.RecordFailure("jobA")
	tr.RecordFailure("jobA")
	tr.RecordSuccess("jobB")
	if tr.Snapshot("jobB").Failed != 0 {
		t.Error("jobB should have no failures")
	}
	if tr.Snapshot("jobA").Failed != 2 {
		t.Error("jobA should have 2 failures")
	}
}
