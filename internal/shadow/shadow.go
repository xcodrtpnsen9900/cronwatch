// Package shadow provides a shadow-mode runner that executes jobs in parallel
// with the primary executor and compares outcomes without affecting production alerts.
package shadow

import (
	"context"
	"log"
	"sync"
	"time"
)

// Result holds the outcome of a shadow execution.
type Result struct {
	Job      string
	Duration time.Duration
	Err      error
}

// Runner executes a function in shadow mode: the primary fn runs normally;
// the shadow fn runs concurrently and its result is only logged.
type Runner struct {
	mu      sync.Mutex
	results []Result
	log     *log.Logger
}

// New creates a new shadow Runner using the provided logger.
func New(logger *log.Logger) *Runner {
	return &Runner{log: logger}
}

// Run executes primary immediately and shadow concurrently.
// The shadow result is stored internally and never blocks the caller.
func (r *Runner) Run(ctx context.Context, job string, primary, shadow func(context.Context) error) error {
	go func() {
		start := time.Now()
		err := shadow(ctx)
		res := Result{
			Job:      job,
			Duration: time.Since(start),
			Err:      err,
		}
		r.mu.Lock()
		r.results = append(r.results, res)
		r.mu.Unlock()
		if err != nil {
			r.log.Printf("[shadow] job=%s err=%v duration=%s", job, err, res.Duration)
		} else {
			r.log.Printf("[shadow] job=%s ok duration=%s", job, res.Duration)
		}
	}()
	return primary(ctx)
}

// Results returns a snapshot of all shadow results collected so far.
func (r *Runner) Results() []Result {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Result, len(r.results))
	copy(out, r.results)
	return out
}

// Reset clears all stored shadow results.
func (r *Runner) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results = r.results[:0]
}
