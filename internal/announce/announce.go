// Package announce provides a broadcast mechanism for notifying multiple
// subscribers when a cron job changes state (started, succeeded, failed).
package announce

import (
	"sync"
	"time"
)

// EventKind classifies a job lifecycle event.
type EventKind string

const (
	KindStarted   EventKind = "started"
	KindSucceeded EventKind = "succeeded"
	KindFailed    EventKind = "failed"
)

// Event carries information about a single job lifecycle transition.
type Event struct {
	Job       string
	Kind      EventKind
	OccurredAt time.Time
	Message   string
}

// Handler is a callback invoked for each broadcast event.
type Handler func(Event)

// Broadcaster fans out events to all registered handlers.
type Broadcaster struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New returns an initialised Broadcaster.
func New() *Broadcaster {
	return &Broadcaster{
		handlers: make(map[string]Handler),
	}
}

// Subscribe registers a named handler. Registering with the same name
// replaces the previous handler.
func (b *Broadcaster) Subscribe(name string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
}

// Unsubscribe removes a previously registered handler.
func (b *Broadcaster) Unsubscribe(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, name)
}

// Publish delivers e to every registered handler. Handlers are invoked
// synchronously in an unspecified order.
func (b *Broadcaster) Publish(e Event) {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now().UTC()
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers {
		h(e)
	}
}

// Len returns the number of active subscribers.
func (b *Broadcaster) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}
