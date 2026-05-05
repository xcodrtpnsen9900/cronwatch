// Package scheduler provides cron expression parsing and schedule
// evaluation utilities used by the cronwatch monitor.
//
// It wraps the robfig/cron parser to expose a minimal API:
//
//	// Parse a standard 5-field cron expression.
//	sched, err := scheduler.Parse("*/5 * * * *")
//
//	// Check whether the job was expected in a time window.
//	if sched.WasExpected(since, until) {
//	    // alert if no heartbeat was received
//	}
//
// This package is intentionally stateless; timing state is managed
// by the monitor package.
package scheduler
