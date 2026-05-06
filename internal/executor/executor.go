// Package executor ties together the runner, history, and notifier
// packages to execute a cron job command, record the result, and
// send alerts when the job fails or recovers.
package executor

import (
	"context"
	"log"
	"time"

	"github.com/example/cronwatch/internal/alert"
	"github.com/example/cronwatch/internal/history"
	"github.com/example/cronwatch/internal/notifier"
	"github.com/example/cronwatch/internal/runner"
)

// Executor runs a named cron job command and handles alerting.
type Executor struct {
	name     string
	command  string
	args     []string
	runner   *runner.Runner
	history  *history.History
	notifier *notifier.Notifier
	alert    *alert.Builder
}

// New creates an Executor for the given job name and shell command.
func New(name, command string, args []string, h *history.History, n *notifier.Notifier, ab *alert.Builder) *Executor {
	return &Executor{
		name:     name,
		command:  command,
		args:     args,
		runner:   runner.New(name),
		history:  h,
		notifier: n,
		alert:    ab,
	}
}

// Run executes the job, records the result, and fires alerts as needed.
func (e *Executor) Run(ctx context.Context) error {
	result := e.runner.Run(ctx, e.command, e.args...)

	entry := history.Entry{
		JobName:   e.name,
		StartedAt: result.StartedAt,
		Finished:  true,
		ExitCode:  result.ExitCode,
		Output:    result.Output,
	}
	e.history.Record(entry)

	if result.ExitCode != 0 {
		payload := e.alert.Failed(e.name, result.ExitCode, result.Output)
		if err := e.notifier.Send(ctx, payload); err != nil {
			log.Printf("executor: failed to send alert for job %q: %v", e.name, err)
		}
		return result.Err
	}

	// Check if the previous run had failed; if so, send a recovery alert.
	if prev, ok := e.history.Latest(e.name); ok && prev.ExitCode != 0 {
		payload := e.alert.Recovered(e.name, time.Since(prev.StartedAt))
		if err := e.notifier.Send(ctx, payload); err != nil {
			log.Printf("executor: failed to send recovery alert for job %q: %v", e.name, err)
		}
	}

	return nil
}
