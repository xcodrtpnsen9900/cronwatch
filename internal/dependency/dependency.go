// Package dependency tracks inter-job dependencies, ensuring a job is only
// considered eligible to run after its declared prerequisites have succeeded.
package dependency

import (
	"fmt"
	"sync"
	"time"
)

// State represents the last known outcome of a job.
type State int

const (
	StateUnknown State = iota
	StateSuccess
	StateFailure
)

// Entry records the last run state and time for a job.
type Entry struct {
	Job       string
	State     State
	UpdatedAt time.Time
}

// Store tracks job states and their declared dependencies.
type Store struct {
	mu      sync.RWMutex
	states  map[string]Entry
	deps    map[string][]string // job -> required jobs
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		states: make(map[string]Entry),
		deps:   make(map[string][]string),
	}
}

// Declare registers the dependencies for a job. Calling Declare again
// replaces the previous dependency list.
func (s *Store) Declare(job string, requires []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deps[job] = append([]string(nil), requires...)
}

// Record updates the last known state for a job.
func (s *Store) Record(job string, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[job] = Entry{Job: job, State: state, UpdatedAt: time.Now()}
}

// Ready returns true when all declared dependencies for job have last
// completed with StateSuccess. If the job has no declared dependencies it is
// always considered ready. An error is returned if a required job has never
// been recorded.
func (s *Store) Ready(job string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reqs, ok := s.deps[job]
	if !ok || len(reqs) == 0 {
		return true, nil
	}

	for _, req := range reqs {
		entry, exists := s.states[req]
		if !exists {
			return false, fmt.Errorf("dependency %q has never run", req)
		}
		if entry.State != StateSuccess {
			return false, nil
		}
	}
	return true, nil
}

// Snapshot returns a copy of all recorded states.
func (s *Store) Snapshot() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.states))
	for _, e := range s.states {
		out = append(out, e)
	}
	return out
}
