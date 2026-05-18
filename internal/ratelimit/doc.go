// Package ratelimit provides per-job alert rate limiting for cronwatch.
//
// A RateLimit enforces a minimum cooldown period between successive alerts
// for the same job, preventing alert storms when a job fails repeatedly in
// a short window.
//
// Usage:
//
//	rl := ratelimit.New(5 * time.Minute)
//	if rl.Allow("my-job") {
//		// send alert
//	}
package ratelimit
