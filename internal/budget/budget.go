// Package budget implements an error-budget tracker for cron jobs.
// It tracks the ratio of failed runs to total runs within a rolling
// window and reports when the budget has been exhausted.
package budget

import (
	"fmt"
	"sync"
	"time"
)

// Entry records a single budget snapshot for a job.
type Entry struct {
	Job       string
	Total     int
	Failed    int
	Remaining float64 // fraction of budget remaining, 0.0–1.0
	Exhausted bool
	At        time.Time
}

// Tracker maintains per-job error-budget state.
type Tracker struct {
	mu       sync.Mutex
	threshold float64 // maximum allowed failure ratio, e.g. 0.05 = 5%
	jobs     map[string]*state
}

type state struct {
	total  int
	failed int
}

// New creates a Tracker with the given failure-ratio threshold (0 < threshold <= 1).
// For example, threshold=0.05 means at most 5% of runs may fail.
func New(threshold float64) (*Tracker, error) {
	if threshold <= 0 || threshold > 1 {
		return nil, fmt.Errorf("budget: threshold must be in (0, 1], got %v", threshold)
	}
	return &Tracker{
		threshold: threshold,
		jobs:      make(map[string]*state),
	}, nil
}

// RecordSuccess increments the total run count for job.
func (t *Tracker) RecordSuccess(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.ensure(job)
	s.total++
}

// RecordFailure increments both the total and failed run counts for job.
func (t *Tracker) RecordFailure(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.ensure(job)
	s.total++
	s.failed++
}

// Snapshot returns the current budget entry for job.
func (t *Tracker) Snapshot(job string) Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.ensure(job)
	var ratio float64
	if s.total > 0 {
		ratio = float64(s.failed) / float64(s.total)
	}
	remaining := 1.0 - (ratio / t.threshold)
	if remaining < 0 {
		remaining = 0
	}
	return Entry{
		Job:       job,
		Total:     s.total,
		Failed:    s.failed,
		Remaining: remaining,
		Exhausted: ratio > t.threshold,
		At:        time.Now(),
	}
}

// Reset clears the counters for job.
func (t *Tracker) Reset(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.jobs, job)
}

func (t *Tracker) ensure(job string) *state {
	if _, ok := t.jobs[job]; !ok {
		t.jobs[job] = &state{}
	}
	return t.jobs[job]
}
