// Package graceful provides a shutdown coordinator that waits for in-flight
// cron jobs to complete before the process exits.
package graceful

import (
	"context"
	"sync"
	"time"
)

// Coordinator tracks active jobs and provides a clean shutdown mechanism.
type Coordinator struct {
	mu      sync.Mutex
	active  map[string]struct{}
	wg      sync.WaitGroup
	done    chan struct{}
	timeout time.Duration
}

// New creates a Coordinator with the given shutdown timeout.
func New(timeout time.Duration) *Coordinator {
	return &Coordinator{
		active:  make(map[string]struct{}),
		done:    make(chan struct{}),
		timeout: timeout,
	}
}

// Acquire registers a job as in-flight. Returns false if shutdown has begun.
func (c *Coordinator) Acquire(jobID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return false
	default:
	}
	c.active[jobID] = struct{}{}
	c.wg.Add(1)
	return true
}

// Release marks a job as completed.
func (c *Coordinator) Release(jobID string) {
	c.mu.Lock()
	delete(c.active, jobID)
	c.mu.Unlock()
	c.wg.Done()
}

// Active returns the set of currently running job IDs.
func (c *Coordinator) Active() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	ids := make([]string, 0, len(c.active))
	for id := range c.active {
		ids = append(ids, id)
	}
	return ids
}

// Shutdown signals no new jobs should start and waits up to the configured
// timeout for in-flight jobs to finish. Returns context.DeadlineExceeded if
// the timeout is reached before all jobs complete.
func (c *Coordinator) Shutdown(ctx context.Context) error {
	close(c.done)

	finished := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(finished)
	}()

	timer := time.NewTimer(c.timeout)
	defer timer.Stop()

	select {
	case <-finished:
		return nil
	case <-timer.C:
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
