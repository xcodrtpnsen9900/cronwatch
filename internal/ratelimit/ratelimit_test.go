package ratelimit_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/ratelimit"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Minute)
	if !l.Allow("job1") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_SecondCallSuppressed(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	if l.Allow("job1") {
		t.Fatal("second call within cooldown should be suppressed")
	}
}

func TestAllow_AfterCooldownPermitted(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("job1")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("job1") {
		t.Fatal("call after cooldown should be allowed")
	}
}

func TestAllow_IndependentJobs(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	if !l.Allow("job2") {
		t.Fatal("different job should be independent")
	}
}

func TestReset_AllowsImmediateAlert(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	l.Reset("job1")
	if !l.Allow("job1") {
		t.Fatal("after Reset, next call should be allowed")
	}
}

func TestResetAll_ClearsAllJobs(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("job1")
	l.Allow("job2")
	l.ResetAll()
	if !l.Allow("job1") || !l.Allow("job2") {
		t.Fatal("after ResetAll, all jobs should be allowed")
	}
}

func TestAllow_ZeroCooldownAlwaysPermits(t *testing.T) {
	l := ratelimit.New(0)
	for i := 0; i < 5; i++ {
		if !l.Allow("job1") {
			t.Fatalf("zero cooldown: call %d should be allowed", i)
		}
	}
}

func TestLastAlert_RecordsTime(t *testing.T) {
	l := ratelimit.New(time.Minute)
	before := time.Now()
	l.Allow("job1")
	after := time.Now()

	t2, ok := l.LastAlert("job1")
	if !ok {
		t.Fatal("expected last alert entry")
	}
	if t2.Before(before) || t2.After(after) {
		t.Errorf("last alert time %v not in expected range [%v, %v]", t2, before, after)
	}
}

func TestLastAlert_MissingJob(t *testing.T) {
	l := ratelimit.New(time.Minute)
	_, ok := l.LastAlert("missing")
	if ok {
		t.Fatal("expected no entry for unseen job")
	}
}
