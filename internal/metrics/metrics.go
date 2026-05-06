// Package metrics exposes lightweight in-process counters for cronwatch
// operational visibility.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counters holds atomic counters for key cronwatch events.
type Counters struct {
	AlertsTotal   atomic.Int64
	ChecksTotal   atomic.Int64
	MissedTotal   atomic.Int64
	FailedTotal   atomic.Int64
	RecoveredTotal atomic.Int64
}

// Registry is a thread-safe collection of named job counters plus global
// operational counters.
type Registry struct {
	mu      sync.RWMutex
	global  Counters
	jobs    map[string]*Counters
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{jobs: make(map[string]*Counters)}
}

// Global returns the global Counters (alerts, checks, etc.).
func (r *Registry) Global() *Counters { return &r.global }

// Job returns the Counters for a named job, creating them on first access.
func (r *Registry) Job(name string) *Counters {
	r.mu.RLock()
	c, ok := r.jobs[name]
	r.mu.RUnlock()
	if ok {
		return c
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok = r.jobs[name]; ok {
		return c
	}
	c = &Counters{}
	r.jobs[name] = c
	return c
}

// Snapshot returns a point-in-time copy of all counters keyed by job name.
// The special key "__global__" holds the global counters.
func (r *Registry) Snapshot() map[string]map[string]int64 {
	out := make(map[string]map[string]int64)
	out["__global__"] = countersToMap(&r.global)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name, c := range r.jobs {
		out[name] = countersToMap(c)
	}
	return out
}

func countersToMap(c *Counters) map[string]int64 {
	return map[string]int64{
		"alerts_total":    c.AlertsTotal.Load(),
		"checks_total":    c.ChecksTotal.Load(),
		"missed_total":    c.MissedTotal.Load(),
		"failed_total":    c.FailedTotal.Load(),
		"recovered_total": c.RecoveredTotal.Load(),
	}
}
