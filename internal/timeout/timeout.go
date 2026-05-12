// Package timeout tracks jobs that have exceeded their expected duration
// and provides a simple store for querying timed-out job state.
package timeout

import (
	"sync"
	"time"
)

// Entry records when a job was marked as timed out.
type Entry struct {
	Job       string
	MarkedAt  time.Time
	Deadline  time.Time
}

// Tracker holds timeout state for monitored jobs.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
	}
}

// Mark records that job has exceeded its deadline.
// Calling Mark for an already-timed-out job is a no-op.
func (t *Tracker) Mark(job string, deadline time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.entries[job]; exists {
		return
	}
	t.entries[job] = Entry{
		Job:      job,
		MarkedAt: time.Now(),
		Deadline: deadline,
	}
}

// Clear removes a job from the timed-out set (e.g. after it completes).
func (t *Tracker) Clear(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, job)
}

// IsTimedOut reports whether job is currently marked as timed out.
func (t *Tracker) IsTimedOut(job string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.entries[job]
	return ok
}

// All returns a snapshot of all currently timed-out jobs.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
