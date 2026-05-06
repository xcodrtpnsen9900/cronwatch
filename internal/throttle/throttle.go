// Package throttle provides a token-bucket style throttle for limiting
// the rate at which alert notifications are dispatched per job.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits calls to at most Burst events per Window per key.
type Throttle struct {
	mu     sync.Mutex
	window time.Duration
	burst  int
	buckets map[string]*bucket
}

type bucket struct {
	tokens int
	reset  time.Time
}

// New creates a Throttle that allows up to burst events per window for each key.
func New(window time.Duration, burst int) *Throttle {
	if burst < 1 {
		burst = 1
	}
	return &Throttle{
		window:  window,
		burst:   burst,
		buckets: make(map[string]*bucket),
	}
}

// Allow reports whether an event for key is permitted under the current quota.
// It consumes one token if permitted.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	b, ok := t.buckets[key]
	if !ok || now.After(b.reset) {
		t.buckets[key] = &bucket{tokens: t.burst - 1, reset: now.Add(t.window)}
		return true
	}
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// Remaining returns the number of tokens left for key in the current window.
func (t *Throttle) Remaining(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	b, ok := t.buckets[key]
	if !ok || now.After(b.reset) {
		return t.burst
	}
	return b.tokens
}

// Reset clears the quota for key, allowing a fresh burst immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, key)
}
