// Package history provides a simple in-memory store for tracking
// the execution history of monitored cron jobs.
package history

import (
	"sync"
	"time"
)

// Entry records a single execution event for a cron job.
type Entry struct {
	JobName   string
	Timestamp time.Time
	Success   bool
	Message   string
}

// Store holds recent execution history for all jobs.
type Store struct {
	mu      sync.RWMutex
	entries map[string][]Entry
	maxPer  int
}

// New creates a new Store that retains up to maxPer entries per job.
func New(maxPer int) *Store {
	if maxPer <= 0 {
		maxPer = 50
	}
	return &Store{
		entries: make(map[string][]Entry),
		maxPer:  maxPer,
	}
}

// Record appends an Entry for the given job, pruning old entries if needed.
func (s *Store) Record(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := s.entries[e.JobName]
	list = append(list, e)
	if len(list) > s.maxPer {
		list = list[len(list)-s.maxPer:]
	}
	s.entries[e.JobName] = list
}

// Latest returns the most recent Entry for a job and whether one exists.
func (s *Store) Latest(jobName string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := s.entries[jobName]
	if len(list) == 0 {
		return Entry{}, false
	}
	return list[len(list)-1], true
}

// All returns a copy of all entries for a job in chronological order.
func (s *Store) All(jobName string) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := s.entries[jobName]
	out := make([]Entry, len(list))
	copy(out, list)
	return out
}
