// Package slo tracks Service Level Objective compliance for cron jobs.
// It records whether each run met its expected deadline and computes
// a rolling compliance ratio over a configurable window.
package slo

import (
	"sync"
	"time"
)

// Entry records a single run's SLO outcome.
type Entry struct {
	Job       string
	At        time.Time
	Met       bool // true if the run completed within the SLO deadline
}

// Snapshot holds the current SLO state for a job.
type Snapshot struct {
	Job        string  `json:"job"`
	Total      int     `json:"total"`
	Met        int     `json:"met"`
	Compliance float64 `json:"compliance_pct"`
}

// Tracker records SLO outcomes and computes compliance per job.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	maxPer  int
	entries map[string][]Entry
}

// New creates a Tracker that retains entries within the given window.
// maxPer caps the number of entries stored per job.
func New(window time.Duration, maxPer int) *Tracker {
	if maxPer <= 0 {
		maxPer = 500
	}
	return &Tracker{
		window:  window,
		maxPer:  maxPer,
		entries: make(map[string][]Entry),
	}
}

// Record adds an outcome for the given job.
func (t *Tracker) Record(job string, at time.Time, met bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e := Entry{Job: job, At: at, Met: met}
	t.entries[job] = append(t.entries[job], e)
	t.prune(job, at)
}

// Snapshot returns the current compliance statistics for a job.
func (t *Tracker) Snapshot(job string) Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(job, time.Now())
	entries := t.entries[job]
	snap := Snapshot{Job: job, Total: len(entries)}
	for _, e := range entries {
		if e.Met {
			snap.Met++
		}
	}
	if snap.Total > 0 {
		snap.Compliance = float64(snap.Met) / float64(snap.Total) * 100
	}
	return snap
}

// All returns snapshots for every tracked job.
func (t *Tracker) All() []Snapshot {
	t.mu.Lock()
	jobs := make([]string, 0, len(t.entries))
	for j := range t.entries {
		jobs = append(jobs, j)
	}
	t.mu.Unlock()

	out := make([]Snapshot, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, t.Snapshot(j))
	}
	return out
}

// prune removes entries outside the window and enforces maxPer. Must be called with lock held.
func (t *Tracker) prune(job string, now time.Time) {
	cutoff := now.Add(-t.window)
	list := t.entries[job]
	start := 0
	for start < len(list) && list[start].At.Before(cutoff) {
		start++
	}
	list = list[start:]
	if len(list) > t.maxPer {
		list = list[len(list)-t.maxPer:]
	}
	t.entries[job] = list
}
