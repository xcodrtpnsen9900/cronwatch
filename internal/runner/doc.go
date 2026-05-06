// Package runner provides a thin wrapper around os/exec for running
// arbitrary shell commands on behalf of monitored cron jobs.
//
// A Runner is constructed with New, specifying the shell interpreter to use.
// The Run method executes a command string, captures combined output, and
// returns a Result that summarises the exit code, duration, and any error.
//
// Example:
//
//	r := runner.New("/bin/sh")
//	res := r.Run(ctx, "backup", "/usr/local/bin/backup.sh")
//	if !res.Succeeded() {
//		// report failure
//	}
package runner
