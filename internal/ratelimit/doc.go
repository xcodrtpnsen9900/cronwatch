// Package ratelimit provides per-job alert rate limiting to prevent
// alert storms when a cron job repeatedly fails. Each job has an
// independent cooldown window; once an alert is sent, subsequent
// alerts for the same job are suppressed until the window expires.
//
// Usage:
//
//	rl := ratelimit.New(5 * time.Minute)
//	if rl.Allow("my-job") {
//		// send alert
//	}
//
// Reset clears the cooldown for a job, useful when a job recovers.
package ratelimit
