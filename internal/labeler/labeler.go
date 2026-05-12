// Package labeler attaches arbitrary key-value metadata labels to jobs.
// Labels can be used by other subsystems (e.g. tag, routing) to filter or
// annotate alert payloads.
package labeler

import (
	"fmt"
	"sync"
)

// Store holds labels for each job.
type Store struct {
	mu     sync.RWMutex
	labels map[string]map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{labels: make(map[string]map[string]string)}
}

// Set replaces all labels for the given job with the provided map.
// Passing a nil or empty map clears the job's labels.
func (s *Store) Set(job string, labels map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	s.labels[job] = copy
}

// Put adds or updates a single label for the given job.
func (s *Store) Put(job, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.labels[job] == nil {
		s.labels[job] = make(map[string]string)
	}
	s.labels[job][key] = value
}

// Get returns the value for a label key on a job.
// The second return value is false when the key is absent.
func (s *Store) Get(job, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.labels[job][key]
	return v, ok
}

// All returns a copy of all labels for the given job.
// Returns an empty map when the job has no labels.
func (s *Store) All(job string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.labels[job]))
	for k, v := range s.labels[job] {
		out[k] = v
	}
	return out
}

// Delete removes a single label key from a job.
func (s *Store) Delete(job, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.labels[job]; !ok {
		return fmt.Errorf("labeler: job %q not found", job)
	}
	delete(s.labels[job], key)
	return nil
}

// Jobs returns the list of job names that have at least one label.
func (s *Store) Jobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.labels))
	for job := range s.labels {
		out = append(out, job)
	}
	return out
}
