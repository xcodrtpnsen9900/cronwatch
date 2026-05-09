// Package overdue tracks jobs that have not reported a heartbeat within
// their expected schedule window, exposing a simple query interface.
package overdue

import (
	"sync"
	"time"
)

// Entry describes a single overdue job.
type Entry struct {
	Job      string
	Expected time.Time
	Detected time.Time
}

// Tracker records and surfaces overdue jobs.
type Tracker struct {
	mu      sync.RWMutex
	overdue map[string]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{overdue: make(map[string]Entry)}
}

// Mark records job as overdue with the given expected execution time.
// Calling Mark on an already-overdue job is idempotent; the original
// detection timestamp is preserved.
func (t *Tracker) Mark(job string, expected time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.overdue[job]; exists {
		return
	}
	t.overdue[job] = Entry{
		Job:      job,
		Expected: expected,
		Detected: time.Now().UTC(),
	}
}

// Clear removes job from the overdue set (e.g. after a successful run).
func (t *Tracker) Clear(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.overdue, job)
}

// IsOverdue reports whether job is currently tracked as overdue.
func (t *Tracker) IsOverdue(job string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.overdue[job]
	return ok
}

// All returns a snapshot of every currently overdue job.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.overdue))
	for _, e := range t.overdue {
		out = append(out, e)
	}
	return out
}
