// Package jitter provides utilities for adding randomised jitter to
// durations, which helps spread out concurrent webhook calls and retry
// attempts so that downstream systems are not hit with a thundering herd.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is the interface satisfied by a random-number source so that
// tests can supply a deterministic replacement.
type Source interface {
	Int63n(n int64) int64
}

// defaultSource wraps the global math/rand functions behind the Source
// interface using a mutex so it is safe for concurrent use.
type defaultSource struct{ mu sync.Mutex }

func (d *defaultSource) Int63n(n int64) int64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rand.Int63n(n) //nolint:gosec // non-crypto jitter is intentional
}

// Jitter adds a uniformly-distributed random offset in the range
// [0, maxJitter) to base and returns the result.  If maxJitter is zero
// or negative the original duration is returned unchanged.
func Jitter(base, maxJitter time.Duration) time.Duration {
	return JitterWith(base, maxJitter, &defaultSource{})
}

// JitterWith is the same as Jitter but accepts a caller-supplied Source
// so that unit tests can produce deterministic output.
func JitterWith(base, maxJitter time.Duration, src Source) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	offset := time.Duration(src.Int63n(int64(maxJitter)))
	return base + offset
}

// Full returns a duration chosen uniformly at random from [0, max).
// It is useful when a caller wants pure jitter with no fixed base.
func Full(max time.Duration) time.Duration {
	return Jitter(0, max)
}
