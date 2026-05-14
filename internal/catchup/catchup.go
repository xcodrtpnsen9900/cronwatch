// Package catchup detects and records missed cron job executions
// that should have fired while cronwatch was offline or restarting.
package catchup

import (
	"sync"
	"time"
)

// MissedRun describes a single execution window that was skipped.
type MissedRun struct {
	Job       string
	Scheduled time.Time
	DetectedAt time.Time
}

// Detector compares the last known checkpoint against the current time
// and enumerates any cron windows that were missed in between.
type Detector struct {
	mu      sync.Mutex
	missed  []MissedRun
	now     func() time.Time
}

// New returns a Detector that uses real wall-clock time.
func New() *Detector {
	return &Detector{now: time.Now}
}

// newWithClock returns a Detector using a custom clock (for testing).
func newWithClock(clock func() time.Time) *Detector {
	return &Detector{now: clock}
}

// Scan evaluates the gap between lastRun and the current time using the
// supplied nextFn (e.g. scheduler.Next) to walk expected fire times.
// Any fire time that falls strictly before now is recorded as missed.
func (d *Detector) Scan(job string, lastRun time.Time, nextFn func(time.Time) time.Time) []MissedRun {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	var found []MissedRun

	t := nextFn(lastRun)
	for t.Before(now) {
		mr := MissedRun{
			Job:        job,
			Scheduled:  t,
			DetectedAt: now,
		}
		found = append(found, mr)
		d.missed = append(d.missed, mr)
		t = nextFn(t)
	}
	return found
}

// All returns a snapshot of every missed run recorded so far.
func (d *Detector) All() []MissedRun {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]MissedRun, len(d.missed))
	copy(out, d.missed)
	return out
}

// Clear removes all recorded missed runs.
func (d *Detector) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.missed = d.missed[:0]
}
