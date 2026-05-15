// Package incident tracks open incidents for cron jobs, grouping
// repeated failures under a single incident ID until the job recovers.
package incident

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the current state of an incident.
type Status string

const (
	StatusOpen     Status = "open"
	StatusResolved Status = "resolved"
)

// Incident holds metadata about an ongoing or resolved failure event.
type Incident struct {
	ID         string
	Job        string
	OpenedAt   time.Time
	ResolvedAt *time.Time
	Failures   int
	Status     Status
}

// Store manages incidents keyed by job name.
type Store struct {
	mu        sync.Mutex
	incidents map[string]*Incident
	counter   int
	now       func() time.Time
}

// New returns an initialised incident Store.
func New() *Store {
	return &Store{
		incidents: make(map[string]*Incident),
		now:       time.Now,
	}
}

// Open records a new failure for job. If no open incident exists one is
// created; otherwise the failure count is incremented.
func (s *Store) Open(job string) *Incident {
	s.mu.Lock()
	defer s.mu.Unlock()

	inc, ok := s.incidents[job]
	if !ok || inc.Status == StatusResolved {
		s.counter++
		inc = &Incident{
			ID:       fmt.Sprintf("INC-%04d", s.counter),
			Job:      job,
			OpenedAt: s.now(),
			Status:   StatusOpen,
		}
		s.incidents[job] = inc
	}
	inc.Failures++
	return copyOf(inc)
}

// Resolve marks the open incident for job as resolved. It is a no-op if
// there is no open incident.
func (s *Store) Resolve(job string) (*Incident, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	inc, ok := s.incidents[job]
	if !ok || inc.Status == StatusResolved {
		return nil, false
	}
	t := s.now()
	inc.ResolvedAt = &t
	inc.Status = StatusResolved
	return copyOf(inc), true
}

// Get returns the current incident for job, if any.
func (s *Store) Get(job string) (*Incident, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	inc, ok := s.incidents[job]
	if !ok {
		return nil, false
	}
	return copyOf(inc), true
}

// All returns a snapshot of every tracked incident.
func (s *Store) All() []*Incident {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Incident, 0, len(s.incidents))
	for _, inc := range s.incidents {
		out = append(out, copyOf(inc))
	}
	return out
}

func copyOf(inc *Incident) *Incident {
	c := *inc
	if inc.ResolvedAt != nil {
		t := *inc.ResolvedAt
		c.ResolvedAt = &t
	}
	return &c
}
