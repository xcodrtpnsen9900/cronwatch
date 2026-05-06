package ratelimit_test

import (
	"testing"
	"time"

	"github.com/example/cronwatch/internal/ratelimit"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Minute)
	if !l.Allow("job1") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestAllow_SecondCallSuppressed(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	if l.Allow("job1") {
		t.Fatal("expected second alert within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldownPermitted(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("job1")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("job1") {
		t.Fatal("expected alert to be allowed after cooldown elapsed")
	}
}

func TestAllow_IndependentJobs(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	if !l.Allow("job2") {
		t.Fatal("expected different job to be allowed independently")
	}
}

func TestReset_AllowsImmediateAlert(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	l.Reset("job1")
	if !l.Allow("job1") {
		t.Fatal("expected alert after reset to be allowed")
	}
}

func TestReset_UnknownJobNoOp(t *testing.T) {
	l := ratelimit.New(time.Minute)
	// should not panic
	l.Reset("nonexistent")
}

func TestSnapshot_ReflectsState(t *testing.T) {
	l := ratelimit.New(time.Minute)
	before := time.Now()
	l.Allow("jobA")
	l.Allow("jobB")

	snap := l.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries in snapshot, got %d", len(snap))
	}
	for _, name := range []string{"jobA", "jobB"} {
		if ts, ok := snap[name]; !ok || ts.Before(before) {
			t.Errorf("snapshot entry for %q missing or has wrong timestamp", name)
		}
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	snap := l.Snapshot()
	delete(snap, "job1")

	snap2 := l.Snapshot()
	if _, ok := snap2["job1"]; !ok {
		t.Fatal("modifying snapshot should not affect limiter state")
	}
}
