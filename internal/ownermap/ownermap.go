// Package ownermap maps cron jobs to their owners (teams or individuals)
// for use in alert routing and escalation.
package ownermap

import (
	"errors"
	"fmt"
	"sync"
)

// Owner holds contact information for a job owner.
type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Team  string `json:"team"`
}

// Store maps job names to their owners.
type Store struct {
	mu     sync.RWMutex
	owners map[string]Owner
}

// New returns an empty owner Store.
func New() *Store {
	return &Store{owners: make(map[string]Owner)}
}

// Set registers or replaces the owner for the given job.
func (s *Store) Set(job string, o Owner) error {
	if job == "" {
		return errors.New("ownermap: job name must not be empty")
	}
	if o.Name == "" {
		return errors.New("ownermap: owner name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.owners[job] = o
	return nil
}

// Get returns the owner for the given job, and a boolean indicating whether
// an entry was found.
func (s *Store) Get(job string) (Owner, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.owners[job]
	return o, ok
}

// Remove deletes the owner mapping for the given job. It is a no-op if the
// job is not registered.
func (s *Store) Remove(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.owners, job)
}

// All returns a snapshot of all job-to-owner mappings.
func (s *Store) All() map[string]Owner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Owner, len(s.owners))
	for k, v := range s.owners {
		out[k] = v
	}
	return out
}

// MustGet returns the owner for the given job or panics. Useful in tests.
func (s *Store) MustGet(job string) Owner {
	o, ok := s.Get(job)
	if !ok {
		panic(fmt.Sprintf("ownermap: no owner registered for job %q", job))
	}
	return o
}
