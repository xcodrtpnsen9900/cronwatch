// Package eventlog provides a structured, in-memory ring buffer for recording
// significant cronwatch lifecycle events (alerts sent, jobs started, errors).
package eventlog

import (
	"sync"
	"time"
)

// Level represents the severity of a log event.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Event is a single structured entry in the event log.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Job       string            `json:"job,omitempty"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Log is a bounded ring buffer of Events.
type Log struct {
	mu       sync.Mutex
	events   []Event
	maxSize  int
	now      func() time.Time
}

// New creates a new Log that retains at most maxSize events.
func New(maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &Log{
		events:  make([]Event, 0, maxSize),
		maxSize: maxSize,
		now:     time.Now,
	}
}

// Add appends a new event to the log, evicting the oldest if full.
func (l *Log) Add(level Level, job, message string, fields map[string]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e := Event{
		Timestamp: l.now().UTC(),
		Level:     level,
		Job:       job,
		Message:   message,
		Fields:    fields,
	}

	if len(l.events) >= l.maxSize {
		l.events = l.events[1:]
	}
	l.events = append(l.events, e)
}

// All returns a snapshot of all events in chronological order.
func (l *Log) All() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Filter returns events matching the given level and/or job name.
// Pass empty strings to skip that filter.
func (l *Log) Filter(level Level, job string) []Event {
	l.mu.Lock()
	defer l.mu.Unlock()

	var out []Event
	for _, e := range l.events {
		if level != "" && e.Level != level {
			continue
		}
		if job != "" && e.Job != job {
			continue
		}
		out = append(out, e)
	}
	return out
}
