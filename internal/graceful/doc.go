// Package graceful provides a shutdown coordinator for cronwatch.
//
// When the process receives a termination signal, the coordinator prevents new
// cron jobs from starting and waits up to a configurable timeout for any
// in-flight executions to complete before the process exits.
//
// Usage:
//
//	coord := graceful.New(30 * time.Second)
//
//	// Before running a job:
//	if !coord.Acquire(jobID) {
//		// shutdown in progress — skip this run
//		return
//	}
//	defer coord.Release(jobID)
//
//	// On SIGTERM / SIGINT:
//	if err := coord.Shutdown(ctx); err != nil {
//		log.Printf("graceful shutdown incomplete: %v", err)
//	}
package graceful
