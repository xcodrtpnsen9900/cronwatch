// Package lockout disables a cron job after a configurable number of
// consecutive failures, preventing further scheduling until the lockout
// window expires or is manually lifted.
//
// Usage:
//
//	store, err := lockout.New(5, 30*time.Minute)
//	if err != nil { ... }
//
//	// On each failure:
//	store.RecordFailure(jobName)
//
//	// Before scheduling:
//	if store.IsLockedOut(jobName) {
//	    // skip this run
//	}
//
//	// On success, reset the counter:
//	store.RecordSuccess(jobName)
package lockout
