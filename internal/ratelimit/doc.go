// Package ratelimit implements per-job cooldown-based alert suppression for
// cronwatch. It prevents notification floods by tracking when the last alert
// was dispatched for each job and rejecting subsequent alerts that arrive
// before the configured cooldown period has elapsed.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//
//	if limiter.Allow(jobName) {
//		// send alert
//	}
//
// Calling Reset clears the suppression state for a single job, while
// ResetAll clears all tracked jobs. Both are useful when a job recovers
// and the next failure should always produce a fresh alert regardless of
// how recently the previous alert was sent.
//
// The limiter is safe for concurrent use by multiple goroutines.
package ratelimit
