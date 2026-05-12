package triage_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/triage"
)

func newClassifier() *triage.Classifier {
	return triage.New(triage.Policy{
		ErrorAfter: 2,
		CritAfter:  4,
		Window:     10 * time.Minute,
	})
}

func TestRecord_FirstFailureIsWarn(t *testing.T) {
	c := newClassifier()
	level := c.Record("job1", time.Now())
	if level != triage.LevelWarn {
		t.Fatalf("expected Warn, got %s", level)
	}
}

func TestRecord_RepeatedFailureIsError(t *testing.T) {
	c := newClassifier()
	now := time.Now()
	c.Record("job1", now)
	level := c.Record("job1", now.Add(time.Minute))
	if level != triage.LevelError {
		t.Fatalf("expected Error, got %s", level)
	}
}

func TestRecord_SustainedFailureIsCritical(t *testing.T) {
	c := newClassifier()
	now := time.Now()
	for i := 0; i < 4; i++ {
		c.Record("job1", now.Add(time.Duration(i)*time.Minute))
	}
	level := c.Record("job1", now.Add(5*time.Minute))
	if level != triage.LevelCritical {
		t.Fatalf("expected Critical, got %s", level)
	}
}

func TestRecord_WindowResetStartsOver(t *testing.T) {
	c := newClassifier()
	now := time.Now()
	c.Record("job1", now)
	c.Record("job1", now.Add(time.Minute))
	// advance beyond window
	level := c.Record("job1", now.Add(20*time.Minute))
	if level != triage.LevelWarn {
		t.Fatalf("expected Warn after window reset, got %s", level)
	}
}

func TestReset_ClearsState(t *testing.T) {
	c := newClassifier()
	now := time.Now()
	c.Record("job1", now)
	c.Record("job1", now.Add(time.Minute))
	c.Reset("job1")
	if got := c.Level("job1"); got != triage.LevelOK {
		t.Fatalf("expected OK after reset, got %s", got)
	}
}

func TestLevel_OKForUnknownJob(t *testing.T) {
	c := newClassifier()
	if got := c.Level("unknown"); got != triage.LevelOK {
		t.Fatalf("expected OK for unknown job, got %s", got)
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level triage.Level
		want  string
	}{
		{triage.LevelOK, "ok"},
		{triage.LevelWarn, "warn"},
		{triage.LevelError, "error"},
		{triage.LevelCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
