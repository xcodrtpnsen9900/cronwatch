package throttle_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/throttle"
)

func newThrottle(burst int) *throttle.Throttle {
	return throttle.New(100*time.Millisecond, burst)
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := newThrottle(3)
	if !th.Allow("job1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_BurstExhausted(t *testing.T) {
	th := newThrottle(2)
	if !th.Allow("job1") {
		t.Fatal("first call should be allowed")
	}
	if !th.Allow("job1") {
		t.Fatal("second call should be allowed")
	}
	if th.Allow("job1") {
		t.Fatal("third call should be denied after burst=2")
	}
}

func TestAllow_AfterWindowResets(t *testing.T) {
	th := throttle.New(30*time.Millisecond, 1)
	if !th.Allow("job1") {
		t.Fatal("first call should be allowed")
	}
	if th.Allow("job1") {
		t.Fatal("second call should be denied")
	}
	time.Sleep(40 * time.Millisecond)
	if !th.Allow("job1") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := newThrottle(1)
	if !th.Allow("a") {
		t.Fatal("key a should be allowed")
	}
	if !th.Allow("b") {
		t.Fatal("key b should be independent and allowed")
	}
}

func TestRemaining(t *testing.T) {
	th := newThrottle(3)
	if r := th.Remaining("job1"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
	th.Allow("job1")
	if r := th.Remaining("job1"); r != 2 {
		t.Fatalf("expected 2 remaining after one allow, got %d", r)
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	th := newThrottle(1)
	th.Allow("job1")
	if th.Allow("job1") {
		t.Fatal("should be denied before reset")
	}
	th.Reset("job1")
	if !th.Allow("job1") {
		t.Fatal("should be allowed after reset")
	}
}
