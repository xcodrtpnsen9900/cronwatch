// Package backoff provides exponential backoff with optional jitter for
// computing retry delays in cronwatch.
//
// Usage:
//
//	p := backoff.DefaultPolicy()
//	for attempt := 0; attempt < maxAttempts; attempt++ {
//		delay := p.Next(attempt)
//		time.Sleep(delay)
//		// ... retry operation
//	}
//
// The Policy.Jitter field (0–1) adds randomness to prevent thundering-herd
// problems when multiple jobs retry simultaneously.
package backoff
