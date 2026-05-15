package graceful_test

import (
	"context"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/graceful"
)

func newCoordinator(t time.Duration) *graceful.Coordinator {
	return graceful.New(t)
}

func TestAcquire_AllowsNewJob(t *testing.T) {
	c := newCoordinator(time.Second)
	if !c.Acquire("job-1") {
		t.Fatal("expected Acquire to return true before shutdown")
	}
	c.Release("job-1")
}

func TestAcquire_BlockedAfterShutdown(t *testing.T) {
	c := newCoordinator(time.Second)
	c.Acquire("job-1")
	go func() {
		time.Sleep(10 * time.Millisecond)
		c.Release("job-1")
	}()
	_ = c.Shutdown(context.Background())

	if c.Acquire("job-2") {
		t.Fatal("expected Acquire to return false after shutdown")
	}
}

func TestShutdown_WaitsForJobs(t *testing.T) {
	c := newCoordinator(time.Second)
	c.Acquire("job-1")

	start := time.Now()
	go func() {
		time.Sleep(30 * time.Millisecond)
		c.Release("job-1")
	}()

	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) < 20*time.Millisecond {
		t.Error("shutdown returned too quickly")
	}
}

func TestShutdown_Timeout(t *testing.T) {
	c := newCoordinator(20 * time.Millisecond)
	c.Acquire("job-1") // never released

	err := c.Shutdown(context.Background())
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestActive_ReturnsRunningJobs(t *testing.T) {
	c := newCoordinator(time.Second)
	c.Acquire("job-a")
	c.Acquire("job-b")

	active := c.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active jobs, got %d", len(active))
	}
	c.Release("job-a")
	c.Release("job-b")
}

func TestShutdown_NoJobs_ReturnsImmediately(t *testing.T) {
	c := newCoordinator(time.Second)
	start := time.Now()
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Error("shutdown took too long with no active jobs")
	}
}
