// Package concurrency provides a Limiter that enforces per-job and global
// ceilings on simultaneously executing cron jobs.
//
// Usage:
//
//	lim := concurrency.New(globalMax, jobMax)
//
//	if err := lim.Acquire(jobName); err != nil {
//		// too many concurrent runs — skip or queue
//	}
//	defer lim.Release(jobName)
//	// ... run the job ...
//
// A zero value for either ceiling means that dimension is unlimited.
package concurrency
