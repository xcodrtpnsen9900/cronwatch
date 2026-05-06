// Package ratelimit provides per-job alert suppression to prevent
// notification floods when a job repeatedly fails within a cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per job and suppresses duplicate
// alerts that arrive within the configured cooldown period.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// A zero or negative cooldown disables suppression (every call is allowed).
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if an alert for jobName should be sent.
// It returns false when a previous alert was sent within the cooldown window.
func (l *Limiter) Allow(jobName string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cooldown <= 0 {
		return true
	}

	now := time.Now()
	if t, ok := l.last[jobName]; ok && now.Sub(t) < l.cooldown {
		return false
	}

	l.last[jobName] = now
	return true
}

// Reset clears the suppression state for jobName, allowing the next
// alert to pass through immediately regardless of the cooldown.
func (l *Limiter) Reset(jobName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, jobName)
}

// ResetAll clears suppression state for every tracked job.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// LastAlert returns the time of the most recent allowed alert for jobName
// and whether an entry exists.
func (l *Limiter) LastAlert(jobName string) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.last[jobName]
	return t, ok
}
