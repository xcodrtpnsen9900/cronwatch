// Package backoff provides exponential backoff calculation for retry delays.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines the parameters for exponential backoff.
type Policy struct {
	// InitialInterval is the delay before the first retry.
	InitialInterval time.Duration
	// MaxInterval caps the computed delay.
	MaxInterval time.Duration
	// Multiplier is applied to the interval on each attempt.
	Multiplier float64
	// Jitter adds randomness as a fraction of the computed delay (0–1).
	Jitter float64
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.2,
	}
}

// Next returns the delay for the given attempt number (0-indexed).
// The delay is capped at MaxInterval and optionally jittered.
func (p Policy) Next(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	base := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt))
	if base > float64(p.MaxInterval) {
		base = float64(p.MaxInterval)
	}
	if p.Jitter > 0 {
		// jitter in range [-jitter*base, +jitter*base]
		delta := p.Jitter * base
		base += (rand.Float64()*2 - 1) * delta //nolint:gosec
		if base < 0 {
			base = 0
		}
	}
	return time.Duration(base)
}

// Sequence returns a slice of delays for n attempts.
func (p Policy) Sequence(n int) []time.Duration {
	delays := make([]time.Duration, n)
	for i := range delays {
		delays[i] = p.Next(i)
	}
	return delays
}
