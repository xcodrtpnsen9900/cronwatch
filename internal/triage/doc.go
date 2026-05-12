// Package triage classifies alert severity for cron jobs based on their
// recent failure history within a configurable rolling window.
//
// Levels progress from Warn (first failure) through Error (repeated failures)
// to Critical (sustained failures). A successful run resets the counter via
// Reset, returning the job to LevelOK.
//
// Example:
//
//	c := triage.New(triage.DefaultPolicy())
//	level := c.Record("nightly-backup", time.Now())
//	fmt.Println(level) // "warn", "error", or "critical"
package triage
