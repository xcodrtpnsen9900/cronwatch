// Package concurrency limits the number of jobs that may run simultaneously.
// A Limiter tracks active slots per job and globally, rejecting new runs when
// the configured ceiling is reached.
package concurrency

import (
	"errors"
	"sync"
)

// ErrLimitReached is returned when no slot is available.
var ErrLimitReached = errors.New("concurrency limit reached")

// Limiter enforces per-job and global concurrency ceilings.
type Limiter struct {
	mu        sync.Mutex
	globalMax int
	jobMax    int
	global    int
	jobs      map[string]int
}

// New creates a Limiter. globalMax is the total simultaneous runs allowed
// across all jobs; jobMax is the per-job ceiling. Zero means unlimited.
func New(globalMax, jobMax int) *Limiter {
	return &Limiter{
		globalMax: globalMax,
		jobMax:    jobMax,
		jobs:      make(map[string]int),
	}
}

// Acquire attempts to claim a slot for job. Returns ErrLimitReached if either
// ceiling would be exceeded.
func (l *Limiter) Acquire(job string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.globalMax > 0 && l.global >= l.globalMax {
		return ErrLimitReached
	}
	if l.jobMax > 0 && l.jobs[job] >= l.jobMax {
		return ErrLimitReached
	}
	l.global++
	l.jobs[job]++
	return nil
}

// Release frees a previously acquired slot for job.
func (l *Limiter) Release(job string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.jobs[job] > 0 {
		l.jobs[job]--
		if l.jobs[job] == 0 {
			delete(l.jobs, job)
		}
	}
	if l.global > 0 {
		l.global--
	}
}

// Snapshot returns the current active count per job and the global total.
func (l *Limiter) Snapshot() (global int, perJob map[string]int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	copy := make(map[string]int, len(l.jobs))
	for k, v := range l.jobs {
		copy[k] = v
	}
	return l.global, copy
}
