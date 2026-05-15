// Package skiplist tracks jobs that have been explicitly skipped for one or
// more upcoming runs. A skipped job will not trigger a missed-run alert for
// the suppressed occurrence.
package skiplist

import (
	"errors"
	"sync"
	"time"
)

// Entry records a single skip window for a job.
type Entry struct {
	Job       string
	Reason    string
	SkipUntil time.Time
	CreatedAt time.Time
}

// Store holds skip windows keyed by job name.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Skip registers a skip window for job until the given time.
func (s *Store) Skip(job, reason string, until time.Time) error {
	if job == "" {
		return errors.New("skiplist: job name must not be empty")
	}
	if until.IsZero() {
		return errors.New("skiplist: until time must not be zero")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[job] = Entry{
		Job:       job,
		Reason:    reason,
		SkipUntil: until,
		CreatedAt: s.now(),
	}
	return nil
}

// IsSkipped reports whether job currently has an active skip window.
func (s *Store) IsSkipped(job string) bool {
	s.mu.RLock()
	e, ok := s.entries[job]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	if s.now().After(e.SkipUntil) {
		s.mu.Lock()
		delete(s.entries, job)
		s.mu.Unlock()
		return false
	}
	return true
}

// Lift removes a skip window for job regardless of its expiry.
func (s *Store) Lift(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, job)
}

// All returns a snapshot of all current (including expired) entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
