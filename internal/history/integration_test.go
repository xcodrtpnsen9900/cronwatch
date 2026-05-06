package history_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/history"
)

// TestConcurrentRecord verifies that concurrent writes do not race or panic.
func TestConcurrentRecord(t *testing.T) {
	s := history.New(20)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Record(history.Entry{
				JobName:   "concurrent-job",
				Timestamp: time.Now(),
				Success:   i%2 == 0,
				Message:   "concurrent write",
			})
		}(i)
	}
	wg.Wait()

	entries := s.All("concurrent-job")
	if len(entries) > 20 {
		t.Errorf("expected at most 20 entries, got %d", len(entries))
	}
}

// TestDefaultMaxPer ensures a zero maxPer defaults to 50.
func TestDefaultMaxPer(t *testing.T) {
	s := history.New(0)

	for i := 0; i < 60; i++ {
		s.Record(history.Entry{
			JobName:   "default-job",
			Timestamp: time.Now(),
			Success:   true,
		})
	}

	entries := s.All("default-job")
	if len(entries) != 50 {
		t.Errorf("expected 50 entries with default max, got %d", len(entries))
	}
}
