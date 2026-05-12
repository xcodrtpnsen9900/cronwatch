// Package tag provides job tagging and filtering for cronwatch.
// Tags allow grouping jobs by environment, team, or any custom label.
package tag

import (
	"sync"
)

// Store holds tag associations for jobs.
type Store struct {
	mu   sync.RWMutex
	tags map[string]map[string]struct{} // job -> set of tags
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		tags: make(map[string]map[string]struct{}),
	}
}

// Set replaces all tags for the given job.
func (s *Store) Set(job string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	set := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		set[t] = struct{}{}
	}
	s.tags[job] = set
}

// Add appends tags to a job without removing existing ones.
func (s *Store) Add(job string, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tags[job] == nil {
		s.tags[job] = make(map[string]struct{})
	}
	for _, t := range tags {
		s.tags[job][t] = struct{}{}
	}
}

// Has reports whether job carries all of the provided tags.
func (s *Store) Has(job string, tags ...string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set := s.tags[job]
	for _, t := range tags {
		if _, ok := set[t]; !ok {
			return false
		}
	}
	return true
}

// Get returns a copy of the tag slice for a job.
func (s *Store) Get(job string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set := s.tags[job]
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}

// JobsWithTag returns all job names that carry the given tag.
func (s *Store) JobsWithTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var jobs []string
	for job, set := range s.tags {
		if _, ok := set[tag]; ok {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

// Remove deletes specific tags from a job.
func (s *Store) Remove(job string, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tags {
		delete(s.tags[job], t)
	}
}
