// Package cooldown provides per-job cooldown tracking to prevent
// alert fatigue by enforcing a minimum quiet period between successive
// alerts for the same job.
package cooldown

import (
	"sync"
	"time"
)

// Store tracks the last alert time for each job and enforces a
// minimum cooldown duration before allowing another alert.
type Store struct {
	mu       sync.Mutex
	entries  map[string]time.Time
	cooldown time.Duration
	now      func() time.Time
}

// New creates a Store with the given cooldown duration.
// Alerts for a job are suppressed until at least d has elapsed
// since the previous alert was recorded.
func New(d time.Duration) *Store {
	return &Store{
		entries:  make(map[string]time.Time),
		cooldown: d,
		now:      time.Now,
	}
}

// Allow reports whether an alert for the given job key is permitted.
// It returns true if no prior alert has been recorded, or if the
// cooldown period has fully elapsed since the last alert.
func (s *Store) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	last, ok := s.entries[key]
	if !ok {
		return true
	}
	return s.now().After(last.Add(s.cooldown))
}

// Record marks the current time as the most recent alert for key.
// Callers should invoke Record immediately after sending an alert.
func (s *Store) Record(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = s.now()
}

// Reset removes the cooldown entry for key, allowing the next call
// to Allow to return true regardless of the cooldown duration.
func (s *Store) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Remaining returns the duration left in the cooldown window for key.
// It returns zero if the key is unknown or the cooldown has elapsed.
func (s *Store) Remaining(key string) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	last, ok := s.entries[key]
	if !ok {
		return 0
	}
	remaining := last.Add(s.cooldown).Sub(s.now())
	if remaining < 0 {
		return 0
	}
	return remaining
}
