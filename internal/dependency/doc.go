// Package dependency tracks inter-job dependencies for cronwatch.
//
// Jobs may declare that they require one or more other jobs to have completed
// successfully before they are eligible to run. The Store records outcomes
// reported by the executor and exposes a Ready query used by the scheduler
// or monitor to gate job execution.
//
// Example usage:
//
//	store := dependency.New()
//	store.Declare("report", []string{"fetch", "transform"})
//	store.Record("fetch", dependency.StateSuccess)
//	store.Record("transform", dependency.StateSuccess)
//	ok, err := store.Ready("report") // true, nil
package dependency
