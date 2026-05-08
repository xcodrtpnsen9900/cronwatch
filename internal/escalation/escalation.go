// Package escalation tracks consecutive failures per job and escalates
// alert severity after a configurable threshold is crossed.
package escalation

import (
	"sync"
	"time"
)

// Level represents the severity of an alert.
type Level int

const (
	LevelWarn  Level = iota // first breach
	LevelError              // threshold crossed
	LevelCrit               // sustained failure
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelCrit:
		return "critical"
	default:
		return "unknown"
	}
}

// Policy defines when to escalate.
type Policy struct {
	// ErrorAfter is the number of consecutive failures before LevelError.
	ErrorAfter int
	// CritAfter is the number of consecutive failures before LevelCrit.
	CritAfter int
}

// DefaultPolicy returns a sensible default escalation policy.
func DefaultPolicy() Policy {
	return Policy{ErrorAfter: 3, CritAfter: 10}
}

type state struct {
	consecutive int
	lastFailure time.Time
}

// Tracker maintains per-job failure counts and derives alert levels.
type Tracker struct {
	mu     sync.Mutex
	states map[string]*state
	policy Policy
}

// New creates a Tracker with the given policy.
func New(p Policy) *Tracker {
	return &Tracker{
		states: make(map[string]*state),
		policy: p,
	}
}

// RecordFailure increments the consecutive failure count for job and
// returns the current escalation Level.
func (t *Tracker) RecordFailure(job string) Level {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(job)
	s.consecutive++
	s.lastFailure = time.Now()
	return t.level(s.consecutive)
}

// RecordSuccess resets the consecutive failure count for job.
func (t *Tracker) RecordSuccess(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.states, job)
}

// Level returns the current escalation level for job without mutating state.
func (t *Tracker) Level(job string) Level {
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.states[job]
	if !ok {
		return LevelWarn
	}
	return t.level(s.consecutive)
}

func (t *Tracker) getOrCreate(job string) *state {
	if s, ok := t.states[job]; ok {
		return s
	}
	s := &state{}
	t.states[job] = s
	return s
}

func (t *Tracker) level(n int) Level {
	switch {
	case n >= t.policy.CritAfter:
		return LevelCrit
	case n >= t.policy.ErrorAfter:
		return LevelError
	default:
		return LevelWarn
	}
}
