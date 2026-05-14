// Package replay provides a mechanism to re-enqueue missed or failed cron
// job runs so they can be retried without waiting for the next scheduled tick.
package replay

import (
	"sync"
	"time"
)

// Entry represents a single job run that is queued for replay.
type Entry struct {
	JobName   string
	Reason    string // "missed" or "failed"
	Scheduled time.Time
	QueuedAt  time.Time
}

// Store holds pending replay entries per job.
type Store struct {
	mu      sync.Mutex
	entries map[string][]Entry
	maxPer  int
}

// New creates a Store that retains at most maxPer entries per job.
// If maxPer is <= 0 it defaults to 10.
func New(maxPer int) *Store {
	if maxPer <= 0 {
		maxPer = 10
	}
	return &Store{
		entries: make(map[string][]Entry),
		maxPer:  maxPer,
	}
}

// Enqueue adds an entry for the given job. Oldest entries are evicted when
// the per-job cap is reached.
func (s *Store) Enqueue(job, reason string, scheduled time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := Entry{
		JobName:   job,
		Reason:    reason,
		Scheduled: scheduled,
		QueuedAt:  time.Now(),
	}
	s.entries[job] = append(s.entries[job], e)
	if len(s.entries[job]) > s.maxPer {
		s.entries[job] = s.entries[job][len(s.entries[job])-s.maxPer:]
	}
}

// Drain removes and returns all pending entries for a job.
func (s *Store) Drain(job string) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := s.entries[job]
	delete(s.entries, job)
	return out
}

// All returns a snapshot of every pending entry across all jobs.
func (s *Store) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []Entry
	for _, list := range s.entries {
		out = append(out, list...)
	}
	return out
}

// Len returns the total number of pending entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	for _, list := range s.entries {
		n += len(list)
	}
	return n
}
