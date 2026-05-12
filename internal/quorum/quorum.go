// Package quorum tracks whether a minimum number of job instances
// have reported success within a time window, useful for distributed
// cron setups where multiple workers may run the same job.
package quorum

import (
	"sync"
	"time"
)

// Entry records a single success report for a job instance.
type Entry struct {
	Instance  string
	ReportedAt time.Time
}

// Status summarises the quorum state for a job.
type Status struct {
	Job       string
	Required  int
	Reported  int
	Met       bool
	Instances []string
}

// Store tracks quorum state for named jobs.
type Store struct {
	mu       sync.Mutex
	policies map[string]int            // job -> required count
	reports  map[string][]Entry        // job -> recent reports
	window   time.Duration
	now      func() time.Time
}

// New creates a Store with the given reporting window.
func New(window time.Duration) *Store {
	return &Store{
		policies: make(map[string]int),
		reports:  make(map[string][]Entry),
		window:   window,
		now:      time.Now,
	}
}

// Require sets the minimum number of instances that must report
// success for the named job within the window.
func (s *Store) Require(job string, n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[job] = n
}

// Report records a success report from the given instance for job.
func (s *Store) Report(job, instance string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.reports[job] = append(s.evict(job, now), Entry{Instance: instance, ReportedAt: now})
}

// Check returns the current quorum Status for job.
func (s *Store) Check(job string) Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	required := s.policies[job]
	active := s.evict(job, s.now())
	s.reports[job] = active
	instances := make([]string, 0, len(active))
	for _, e := range active {
		instances = append(instances, e.Instance)
	}
	return Status{
		Job:       job,
		Required:  required,
		Reported:  len(active),
		Met:       len(active) >= required && required > 0,
		Instances: instances,
	}
}

// evict removes entries outside the window; caller must hold mu.
func (s *Store) evict(job string, now time.Time) []Entry {
	cutoff := now.Add(-s.window)
	var kept []Entry
	for _, e := range s.reports[job] {
		if e.ReportedAt.After(cutoff) {
			kept = append(kept, e)
		}
	}
	return kept
}
