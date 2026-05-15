package concurrency_test

import (
	"sync"
	"testing"

	"github.com/cronwatch/internal/concurrency"
)

func newLimiter(g, j int) *concurrency.Limiter {
	return concurrency.New(g, j)
}

func TestAcquire_FirstCallSucceeds(t *testing.T) {
	l := newLimiter(5, 2)
	if err := l.Acquire("job-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAcquire_JobLimitEnforced(t *testing.T) {
	l := newLimiter(10, 1)
	if err := l.Acquire("job-a"); err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	if err := l.Acquire("job-a"); err != concurrency.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestAcquire_GlobalLimitEnforced(t *testing.T) {
	l := newLimiter(2, 10)
	if err := l.Acquire("a"); err != nil {
		t.Fatal(err)
	}
	if err := l.Acquire("b"); err != nil {
		t.Fatal(err)
	}
	if err := l.Acquire("c"); err != concurrency.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestRelease_AllowsReacquire(t *testing.T) {
	l := newLimiter(1, 1)
	if err := l.Acquire("job-a"); err != nil {
		t.Fatal(err)
	}
	l.Release("job-a")
	if err := l.Acquire("job-a"); err != nil {
		t.Fatalf("after release: %v", err)
	}
}

func TestRelease_UnknownJob_NoError(t *testing.T) {
	l := newLimiter(5, 5)
	l.Release("ghost") // must not panic
}

func TestSnapshot_ReflectsState(t *testing.T) {
	l := newLimiter(10, 5)
	_ = l.Acquire("job-a")
	_ = l.Acquire("job-a")
	_ = l.Acquire("job-b")

	g, jobs := l.Snapshot()
	if g != 3 {
		t.Errorf("global want 3, got %d", g)
	}
	if jobs["job-a"] != 2 {
		t.Errorf("job-a want 2, got %d", jobs["job-a"])
	}
	if jobs["job-b"] != 1 {
		t.Errorf("job-b want 1, got %d", jobs["job-b"])
	}
}

func TestConcurrentAcquireRelease(t *testing.T) {
	l := newLimiter(50, 10)
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire("shared"); err == nil {
				l.Release("shared")
			}
		}()
	}
	wg.Wait()
	g, _ := l.Snapshot()
	if g != 0 {
		t.Errorf("expected global=0 after all releases, got %d", g)
	}
}
