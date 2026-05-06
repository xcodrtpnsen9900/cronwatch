// Package ratelimit provides a simple per-job alert rate limiter to prevent
// alert storms when a job repeatedly fails within a short window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per job and suppresses duplicate alerts
// within the configured cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// Alerts for the same job will be suppressed until cooldown elapses.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether an alert for jobName should be sent.
// It returns true the first time a job is seen, and again only after
// the cooldown window has elapsed since the last allowed alert.
func (l *Limiter) Allow(jobName string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[jobName]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[jobName] = now
	return true
}

// Reset clears the rate-limit state for jobName, allowing the next alert
// through immediately. Useful when a job recovers successfully.
func (l *Limiter) Reset(jobName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, jobName)
}

// Snapshot returns a copy of the current last-alert times keyed by job name.
// Intended for diagnostics and status endpoints.
func (l *Limiter) Snapshot() map[string]time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[string]time.Time, len(l.last))
	for k, v := range l.last {
		out[k] = v
	}
	return out
}
