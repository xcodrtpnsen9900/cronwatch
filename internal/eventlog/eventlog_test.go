package eventlog

import (
	"testing"
	"time"
)

func newLog(max int) *Log {
	l := New(max)
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	l.now = func() time.Time { return fixed }
	return l
}

func TestAdd_And_All(t *testing.T) {
	l := newLog(10)
	l.Add(LevelInfo, "job1", "started", nil)
	l.Add(LevelError, "job2", "failed", map[string]string{"code": "1"})

	events := l.All()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Message != "started" {
		t.Errorf("unexpected message: %s", events[0].Message)
	}
}

func TestAdd_Evicts_Oldest(t *testing.T) {
	l := newLog(3)
	for i := 0; i < 5; i++ {
		l.Add(LevelInfo, "", "msg", nil)
	}
	if len(l.All()) != 3 {
		t.Fatalf("expected ring buffer capped at 3")
	}
}

func TestFilter_ByLevel(t *testing.T) {
	l := newLog(10)
	l.Add(LevelInfo, "a", "ok", nil)
	l.Add(LevelError, "b", "bad", nil)
	l.Add(LevelWarn, "c", "meh", nil)

	errors := l.Filter(LevelError, "")
	if len(errors) != 1 || errors[0].Job != "b" {
		t.Errorf("expected 1 error event for job b, got %+v", errors)
	}
}

func TestFilter_ByJob(t *testing.T) {
	l := newLog(10)
	l.Add(LevelInfo, "alpha", "run", nil)
	l.Add(LevelError, "beta", "fail", nil)
	l.Add(LevelInfo, "alpha", "done", nil)

	results := l.Filter("", "alpha")
	if len(results) != 2 {
		t.Errorf("expected 2 events for alpha, got %d", len(results))
	}
}

func TestFilter_Combined(t *testing.T) {
	l := newLog(10)
	l.Add(LevelError, "alpha", "fail", nil)
	l.Add(LevelError, "beta", "fail", nil)
	l.Add(LevelInfo, "alpha", "ok", nil)

	results := l.Filter(LevelError, "alpha")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := newLog(10)
	l.Add(LevelInfo, "", "first", nil)
	snap := l.All()
	l.Add(LevelInfo, "", "second", nil)
	if len(snap) != 1 {
		t.Error("snapshot should not be affected by later writes")
	}
}
