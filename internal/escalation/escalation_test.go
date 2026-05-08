package escalation_test

import (
	"testing"

	"github.com/cronwatch/cronwatch/internal/escalation"
)

func newTracker() *escalation.Tracker {
	return escalation.New(escalation.Policy{ErrorAfter: 3, CritAfter: 6})
}

func TestRecordFailure_WarnBelowThreshold(t *testing.T) {
	tr := newTracker()
	for i := 0; i < 2; i++ {
		lvl := tr.RecordFailure("job-a")
		if lvl != escalation.LevelWarn {
			t.Fatalf("iteration %d: expected Warn, got %s", i, lvl)
		}
	}
}

func TestRecordFailure_ErrorAtThreshold(t *testing.T) {
	tr := newTracker()
	var lvl escalation.Level
	for i := 0; i < 3; i++ {
		lvl = tr.RecordFailure("job-b")
	}
	if lvl != escalation.LevelError {
		t.Fatalf("expected Error at 3 failures, got %s", lvl)
	}
}

func TestRecordFailure_CritAtHighThreshold(t *testing.T) {
	tr := newTracker()
	var lvl escalation.Level
	for i := 0; i < 6; i++ {
		lvl = tr.RecordFailure("job-c")
	}
	if lvl != escalation.LevelCrit {
		t.Fatalf("expected Crit at 6 failures, got %s", lvl)
	}
}

func TestRecordSuccess_ResetsCount(t *testing.T) {
	tr := newTracker()
	for i := 0; i < 5; i++ {
		tr.RecordFailure("job-d")
	}
	tr.RecordSuccess("job-d")
	lvl := tr.Level("job-d")
	if lvl != escalation.LevelWarn {
		t.Fatalf("expected Warn after reset, got %s", lvl)
	}
}

func TestLevel_IndependentJobs(t *testing.T) {
	tr := newTracker()
	for i := 0; i < 5; i++ {
		tr.RecordFailure("job-x")
	}
	lvl := tr.Level("job-y")
	if lvl != escalation.LevelWarn {
		t.Fatalf("job-y should be unaffected by job-x, got %s", lvl)
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		lvl  escalation.Level
		want string
	}{
		{escalation.LevelWarn, "warn"},
		{escalation.LevelError, "error"},
		{escalation.LevelCrit, "critical"},
	}
	for _, c := range cases {
		if got := c.lvl.String(); got != c.want {
			t.Errorf("Level(%d).String() = %q, want %q", c.lvl, got, c.want)
		}
	}
}

func TestDefaultPolicy(t *testing.T) {
	p := escalation.DefaultPolicy()
	if p.ErrorAfter <= 0 || p.CritAfter <= p.ErrorAfter {
		t.Fatalf("unexpected default policy: %+v", p)
	}
}
