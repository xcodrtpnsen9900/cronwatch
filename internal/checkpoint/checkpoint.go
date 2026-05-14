// Package checkpoint tracks the last successful completion time for each
// cron job, enabling cronwatch to detect missed runs across restarts.
package checkpoint

import (
	"sync"
	"time"
)

// Entry holds the last successful run time and optional metadata for a job.
type Entry struct {
	Job       string
	LastOK    time.Time
	RunCount  int64
}

// Store persists checkpoint entries in memory.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Entry)}
}

// Record marks a successful completion for the given job at the given time.
func (s *Store) Record(job string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.entries[job]
	e.Job = job
	e.LastOK = at
	e.RunCount++
	s.entries[job] = e
}

// Get returns the checkpoint entry for a job and whether it exists.
func (s *Store) Get(job string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[job]
	return e, ok
}

// All returns a snapshot of all checkpoint entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// Reset removes the checkpoint for a job, causing the next check to treat
// the job as never having run.
func (s *Store) Reset(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, job)
}
