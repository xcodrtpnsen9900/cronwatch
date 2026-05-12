// Package triage classifies alert severity based on failure history.
package triage

import (
	"sync"
	"time"
)

// Level represents an alert severity level.
type Level int

const (
	LevelOK       Level = iota
	LevelWarn            // first failure or single miss
	LevelError           // repeated failures within window
	LevelCritical        // sustained failures exceeding threshold
)

func (l Level) String() string {
	switch l {
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	default:
		return "ok"
	}
}

// Policy controls thresholds for level promotion.
type Policy struct {
	ErrorAfter    int           // failures before Error
	CritAfter     int           // failures before Critical
	Window        time.Duration // rolling window for counts
}

// DefaultPolicy returns sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		ErrorAfter: 2,
		CritAfter:  5,
		Window:     30 * time.Minute,
	}
}

type entry struct {
	count int
	first time.Time
}

// Classifier assigns severity levels to jobs.
type Classifier struct {
	mu     sync.Mutex
	policy Policy
	states map[string]*entry
}

// New creates a Classifier with the given policy.
func New(p Policy) *Classifier {
	return &Classifier{
		policy: p,
		states: make(map[string]*entry),
	}
}

// Record registers a failure for job and returns the resulting Level.
func (c *Classifier) Record(job string, now time.Time) Level {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.states[job]
	if !ok || now.Sub(e.first) > c.policy.Window {
		e = &entry{first: now}
		c.states[job] = e
	}
	e.count++

	switch {
	case e.count >= c.policy.CritAfter:
		return LevelCritical
	case e.count >= c.policy.ErrorAfter:
		return LevelError
	default:
		return LevelWarn
	}
}

// Reset clears the failure state for job (e.g. on success).
func (c *Classifier) Reset(job string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.states, job)
}

// Level returns the current level for job without recording a new failure.
func (c *Classifier) Level(job string) Level {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.states[job]
	if !ok {
		return LevelOK
	}
	switch {
	case e.count >= c.policy.CritAfter:
		return LevelCritical
	case e.count >= c.policy.ErrorAfter:
		return LevelError
	default:
		return LevelWarn
	}
}
