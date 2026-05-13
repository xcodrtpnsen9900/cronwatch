// Package runlock prevents concurrent execution of the same cron job
// by maintaining a set of active job locks with optional TTL expiry.
package runlock

import (
	"fmt"
	"sync"
	"time"
)

// Store tracks which jobs are currently running.
type Store struct {
	mu      sync.Mutex
	locks   map[string]time.Time // job name -> acquired at
	ttl     time.Duration
	nowFunc func() time.Time
}

// New returns a Store where locks expire after ttl if not explicitly released.
// A zero ttl means locks never expire automatically.
func New(ttl time.Duration) *Store {
	return &Store{
		locks:   make(map[string]time.Time),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Acquire attempts to acquire a lock for job. It returns an error if the job
// is already running (and the lock has not expired).
func (s *Store) Acquire(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if at, ok := s.locks[job]; ok {
		if s.ttl == 0 || s.nowFunc().Sub(at) < s.ttl {
			return fmt.Errorf("runlock: job %q is already running (acquired at %s)", job, at.Format(time.RFC3339))
		}
		// lock has expired — allow re-acquisition
	}
	s.locks[job] = s.nowFunc()
	return nil
}

// Release removes the lock for job. It is a no-op if no lock exists.
func (s *Store) Release(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.locks, job)
}

// IsLocked reports whether job currently holds an active (non-expired) lock.
func (s *Store) IsLocked(job string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	at, ok := s.locks[job]
	if !ok {
		return false
	}
	if s.ttl > 0 && s.nowFunc().Sub(at) >= s.ttl {
		return false
	}
	return true
}

// Active returns the names of all jobs that currently hold active locks.
func (s *Store) Active() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nowFunc()
	out := make([]string, 0, len(s.locks))
	for job, at := range s.locks {
		if s.ttl == 0 || now.Sub(at) < s.ttl {
			out = append(out, job)
		}
	}
	return out
}
