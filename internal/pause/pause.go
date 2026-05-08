// Package pause provides a mechanism to temporarily pause alerting
// for specific cron jobs, suppressing notifications during maintenance windows
// or known downtime periods.
package pause

import (
	"sync"
	"time"
)

// Store tracks paused jobs and their resume times.
type Store struct {
	mu    sync.RWMutex
	entries map[string]time.Time
	now   func() time.Time
}

// New returns a new pause Store.
func New() *Store {
	return &Store{
		entries: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Pause marks a job as paused until the given duration has elapsed.
// Calling Pause on an already-paused job extends the pause window.
func (s *Store) Pause(jobName string, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[jobName] = s.now().Add(d)
}

// Resume immediately lifts the pause for a job.
func (s *Store) Resume(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, jobName)
}

// IsPaused reports whether the given job is currently paused.
func (s *Store) IsPaused(jobName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resumeAt, ok := s.entries[jobName]
	if !ok {
		return false
	}
	if s.now().After(resumeAt) {
		return false
	}
	return true
}

// PausedUntil returns the time at which the job's pause expires.
// The second return value is false if the job is not currently paused.
func (s *Store) PausedUntil(jobName string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resumeAt, ok := s.entries[jobName]
	if !ok {
		return time.Time{}, false
	}
	if s.now().After(resumeAt) {
		return time.Time{}, false
	}
	return resumeAt, true
}

// Evict removes all expired pause entries to keep the map from growing
// unboundedly. It is safe to call periodically in a background goroutine.
func (s *Store) Evict() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	for job, resumeAt := range s.entries {
		if now.After(resumeAt) {
			delete(s.entries, job)
		}
	}
}
