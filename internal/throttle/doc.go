// Package throttle implements a per-key token-bucket throttle used to cap
// the number of alert notifications dispatched within a sliding time window.
//
// # Usage
//
//	th := throttle.New(time.Minute, 5) // allow 5 alerts per minute per job
//	if th.Allow(jobName) {
//		// send alert
//	}
//
// Tokens are replenished automatically when the window expires. Calling
// Reset clears the quota for a specific key, useful after a job recovers.
package throttle
