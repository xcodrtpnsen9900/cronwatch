// Package runcounter tracks the total number of executions per job,
// providing a lightweight counter store with snapshot support.
package runcounter

import (
	"sync"
)

// Store holds per-job execution counters.
type Store struct {
	mu       sync.RWMutex
	counters map[string]int64
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		counters: make(map[string]int64),
	}
}

// Increment adds one to the counter for the given job and returns the new value.
func (s *Store) Increment(job string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[job]++
	return s.counters[job]
}

// Get returns the current counter value for job. Returns 0 if unknown.
func (s *Store) Get(job string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.counters[job]
}

// Reset sets the counter for job back to zero.
func (s *Store) Reset(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counters, job)
}

// Snapshot returns a copy of all counters keyed by job name.
func (s *Store) Snapshot() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int64, len(s.counters))
	for k, v := range s.counters {
		out[k] = v
	}
	return out
}

// Jobs returns the list of job names that have at least one recorded run.
func (s *Store) Jobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.counters))
	for k := range s.counters {
		names = append(names, k)
	}
	return names
}
