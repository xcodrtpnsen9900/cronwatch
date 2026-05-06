// Package runner executes shell commands for monitored cron jobs,
// capturing stdout/stderr output and exit codes. It supports context
// cancellation for timeout enforcement and returns structured results
// that can be recorded in history and evaluated for alerting.
package runner
