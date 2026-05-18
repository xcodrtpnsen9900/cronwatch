package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	lastAlert time.Time
}

// RateLimit enforces a per-job cooldown between successive alerts.
type RateLimit struct {
	mu       sync.Mutex
	entries  map[string]*entry
	cooldown time.Duration
	now      func() time.Time
}

// New creates a RateLimit with the given cooldown duration.
func New(cooldown time.Duration) *RateLimit {
	return &RateLimit{
		entries:  make(map[string]*entry),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true if an alert for the given job is permitted under the
// current cooldown policy, and records the attempt time. Subsequent calls
// within the cooldown window return false.
func (rl *RateLimit) Allow(job string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()
	e, ok := rl.entries[job]
	if !ok {
		rl.entries[job] = &entry{lastAlert: now}
		return true
	}
	if now.Sub(e.lastAlert) >= rl.cooldown {
		e.lastAlert = now
		return true
	}
	return false
}

// Reset clears the cooldown state for the given job, allowing the next call
// to Allow to succeed immediately.
func (rl *RateLimit) Reset(job string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.entries, job)
}

// ResetAll clears cooldown state for every tracked job.
func (rl *RateLimit) ResetAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.entries = make(map[string]*entry)
}

// Remaining returns the duration until the next alert is permitted for job.
// Returns zero if the job is not throttled or the cooldown has elapsed.
func (rl *RateLimit) Remaining(job string) time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	e, ok := rl.entries[job]
	if !ok {
		return 0
	}
	remaining := rl.cooldown - rl.now().Sub(e.lastAlert)
	if remaining < 0 {
		return 0
	}
	return remaining
}
