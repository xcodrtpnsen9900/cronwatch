// Package audit provides a lightweight, in-memory audit log that records
// significant cronwatch events (alert sent, circuit opened, job recovered)
// for inspection via the status endpoint or tests.
package audit

import (
	"sync"
	"time"
)

// EventKind classifies an audit event.
type EventKind string

const (
	EventAlertSent      EventKind = "alert_sent"
	EventCircuitOpen    EventKind = "circuit_open"
	EventCircuitClose   EventKind = "circuit_close"
	EventJobRecovered   EventKind = "job_recovered"
	EventJobFailed      EventKind = "job_failed"
	EventJobMissed      EventKind = "job_missed"
)

// Entry is a single audit log record.
type Entry struct {
	Time    time.Time `json:"time"`
	Job     string    `json:"job"`
	Kind    EventKind `json:"kind"`
	Message string    `json:"message,omitempty"`
}

// Log is a bounded, thread-safe audit log.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a Log that retains at most maxSize entries.
// If maxSize <= 0 it defaults to 200.
func New(maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &Log{maxSize: maxSize}
}

// Record appends an entry to the log, pruning the oldest entry when
// the capacity is exceeded.
func (l *Log) Record(job string, kind EventKind, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, Entry{
		Time:    time.Now().UTC(),
		Job:     job,
		Kind:    kind,
		Message: msg,
	})
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// All returns a copy of all entries, oldest first.
func (l *Log) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// ForJob returns entries for a specific job, oldest first.
func (l *Log) ForJob(job string) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	var out []Entry
	for _, e := range l.entries {
		if e.Job == job {
			out = append(out, e)
		}
	}
	return out
}
