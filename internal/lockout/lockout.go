// Package lockout temporarily disables a job after repeated failures,
// preventing it from being scheduled until the lockout window expires.
package lockout

import (
	"fmt"
	"sync"
	"time"
)

// Store tracks lockout windows for jobs.
type Store struct {
	mu       sync.Mutex
	entries  map[string]time.Time
	threshold int
	window   time.Duration
	failures map[string]int
	now      func() time.Time
}

// New creates a Store that locks out a job after threshold consecutive
// failures for the given window duration.
func New(threshold int, window time.Duration) (*Store, error) {
	if threshold < 1 {
		return nil, fmt.Errorf("lockout: threshold must be >= 1, got %d", threshold)
	}
	if window <= 0 {
		return nil, fmt.Errorf("lockout: window must be positive")
	}
	return &Store{
		entries:   make(map[string]time.Time),
		failures:  make(map[string]int),
		threshold: threshold,
		window:    window,
		now:       time.Now,
	}, nil
}

// RecordFailure increments the failure count for job. If the count reaches
// the threshold, the job is locked out for the configured window.
func (s *Store) RecordFailure(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures[job]++
	if s.failures[job] >= s.threshold {
		s.entries[job] = s.now().Add(s.window)
	}
}

// RecordSuccess clears the failure count and any active lockout for job.
func (s *Store) RecordSuccess(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.failures, job)
	delete(s.entries, job)
}

// IsLockedOut reports whether job is currently locked out.
func (s *Store) IsLockedOut(job string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, ok := s.entries[job]
	if !ok {
		return false
	}
	if s.now().After(expiry) {
		delete(s.entries, job)
		delete(s.failures, job)
		return false
	}
	return true
}

// Lift removes any active lockout for job immediately.
func (s *Store) Lift(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, job)
	delete(s.failures, job)
}

// All returns a snapshot of currently locked-out job names and their expiry times.
func (s *Store) All() map[string]time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]time.Time, len(s.entries))
	now := s.now()
	for job, expiry := range s.entries {
		if now.Before(expiry) {
			out[job] = expiry
		}
	}
	return out
}
