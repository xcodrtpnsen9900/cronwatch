package window

import (
	"testing"
	"time"
)

func newCounter(size time.Duration, buckets int) *Counter {
	c := New(size, buckets)
	return c
}

func TestAdd_And_Total(t *testing.T) {
	c := newCounter(time.Second, 10)
	c.Add(3)
	c.Add(2)
	if got := c.Total(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestTotal_ZeroInitially(t *testing.T) {
	c := newCounter(time.Second, 5)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	c := newCounter(time.Second, 5)
	c.Add(10)
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestEviction_AfterWindow(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	c := New(100*time.Millisecond, 10)
	c.clock = func() time.Time { return now }

	c.Add(7)

	// advance time beyond the window
	now = now.Add(200 * time.Millisecond)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after window expired, got %d", got)
	}
}

func TestBucketAdvance_WritesNewBucket(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	c := New(100*time.Millisecond, 10)
	c.clock = func() time.Time { return now }

	c.Add(4)

	// advance one bucket interval
	now = now.Add(15 * time.Millisecond)
	c.Add(6)

	if got := c.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestNew_MinOneBucket(t *testing.T) {
	c := New(time.Second, 0)
	if c.buckets != 1 {
		t.Fatalf("expected buckets to be clamped to 1, got %d", c.buckets)
	}
}

func TestConcurrent_Add(t *testing.T) {
	c := newCounter(time.Second, 10)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.Add(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := c.Total(); got < 1 {
		t.Fatal("expected at least one count after concurrent adds")
	}
}
