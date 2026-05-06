// Package circuitbreaker implements a simple circuit breaker that stops
// forwarding webhook alerts when the downstream endpoint is repeatedly
// failing, and automatically re-tries after a configurable cooldown.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current circuit state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // blocking calls
	StateHalfOpen              // testing recovery
)

// ErrOpen is returned when the circuit is open and calls are blocked.
var ErrOpen = errors.New("circuit breaker is open")

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures      int
	maxFailures   int
	cooldown      time.Duration
	openedAt      time.Time
	now           func() time.Time
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts recovery after cooldown.
func New(maxFailures int, cooldown time.Duration) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
		now:         time.Now,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if blocked.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the breaker to closed state.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit if
// the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current circuit state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
