// Package runlock provides a concurrency guard that prevents a cron job from
// being executed more than once simultaneously.
//
// Typical usage:
//
//	lock := runlock.New(5 * time.Minute) // locks expire after 5 min
//
//	if err := lock.Acquire(jobName); err != nil {
//		// job is already running — skip this invocation
//		return
//	}
//	defer lock.Release(jobName)
//	// ... execute job ...
//
// The optional TTL ensures that a lock left behind by a crashed process is
// automatically released after the specified duration.
package runlock
