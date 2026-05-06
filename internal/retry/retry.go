// Package retry provides configurable retry logic for transient failures
// such as webhook delivery or external command execution.
package retry

import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Policy defines how retries are performed.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Delay is the wait duration between attempts.
	Delay time.Duration
	// Multiplier scales the delay after each failure (1.0 = constant).
	Multiplier float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
		Multiplier:  2.0,
	}
}

// Do executes fn according to p, retrying on non-nil errors.
// The context is checked before each attempt; cancellation stops retries.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.Multiplier <= 0 {
		p.Multiplier = 1.0
	}

	delay := p.Delay
	var lastErr error

	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if lastErr = fn(); lastErr == nil {
			return nil
		}

		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}

	return errors.Join(ErrMaxAttempts, lastErr)
}
