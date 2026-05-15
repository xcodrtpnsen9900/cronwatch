// Package runbook associates jobs with runbook URLs and metadata
// that operators can consult when an alert fires.
package runbook

import (
	"fmt"
	"sync"
)

// Entry holds the runbook information for a single job.
type Entry struct {
	Job     string `json:"job"`
	URL     string `json:"url"`
	Summary string `json:"summary,omitempty"`
}

// Store maps job names to their runbook entries.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Store.
func New() *Store {
	return &Store{entries: make(map[string]Entry)}
}

// Set registers or replaces the runbook entry for a job.
func (s *Store) Set(job, url, summary string) error {
	if job == "" {
		return fmt.Errorf("runbook: job name must not be empty")
	}
	if url == "" {
		return fmt.Errorf("runbook: url must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[job] = Entry{Job: job, URL: url, Summary: summary}
	return nil
}

// Get returns the runbook entry for a job and whether it was found.
func (s *Store) Get(job string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[job]
	return e, ok
}

// Remove deletes the runbook entry for a job. It is a no-op if the
// job is not registered.
func (s *Store) Remove(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, job)
}

// All returns a snapshot of all registered entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
