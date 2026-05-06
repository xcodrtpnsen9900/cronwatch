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
// Calling Reset or ResetAll clears the suppression state, which is useful
// when a job recovers and the next failure should always produce a fresh alert.
package ratelimit
