// Package ratelimit provides per-job alert rate limiting to prevent
// notification storms when a cron job repeatedly fails.
//
// A Limiter tracks the last alert time for each job and suppresses
// duplicate alerts within a configurable cooldown window. Once the
// cooldown expires the next alert is permitted and the window resets.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//	if limiter.Allow(jobName) {
//		// send alert
//	}
package ratelimit
