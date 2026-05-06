// Package runner executes shell commands on behalf of monitored cron jobs
// and reports success or failure back to the monitor.
package runner

import (
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a single job execution.
type Result struct {
	JobName  string
	ExitCode int
	Duration time.Duration
	Err      error
	Output   []byte
}

// Runner executes a shell command and returns a Result.
type Runner struct {
	shell string
}

// New returns a Runner that uses the given shell (e.g. "/bin/sh").
// If shell is empty, "/bin/sh" is used.
func New(shell string) *Runner {
	if shell == "" {
		shell = "/bin/sh"
	}
	return &Runner{shell: shell}
}

// Run executes cmd under the configured shell with the provided context.
// It captures combined stdout+stderr output and records the elapsed time.
func (r *Runner) Run(ctx context.Context, jobName, cmd string) Result {
	start := time.Now()

	c := exec.CommandContext(ctx, r.shell, "-c", cmd) //nolint:gosec
	out, err := c.CombinedOutput()

	result := Result{
		JobName:  jobName,
		Duration: time.Since(start),
		Output:   out,
		Err:      err,
	}

	if c.ProcessState != nil {
		result.ExitCode = c.ProcessState.ExitCode()
	}

	return result
}

// Succeeded returns true when the command exited with code 0 and no error.
func (res Result) Succeeded() bool {
	return res.Err == nil && res.ExitCode == 0
}
