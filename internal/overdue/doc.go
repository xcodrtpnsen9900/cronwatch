// Package overdue provides a thread-safe tracker for cron jobs that have
// exceeded their expected execution window without reporting a heartbeat.
//
// Typical usage:
//
//	tr := overdue.New()
//
//	// When the scheduler detects a missed window:
//	tr.Mark(jobName, expectedAt)
//
//	// When a heartbeat or successful run is received:
//	tr.Clear(jobName)
//
//	// Expose current state over HTTP:
//	http.Handle("/overdue", overdue.Handler(tr))
package overdue
