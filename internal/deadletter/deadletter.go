// Package deadletter stores alerts that could not be delivered so they can
// be retried or inspected later.
package deadletter

import (
	"sync"
	"time"
)

// Entry holds a single undelivered alert payload.
type Entry struct {
	ID        string
	JobName   string
	Payload   []byte
	Reason    string
	CreatedAt time.Time
	Attempts  int
}

// Store is a bounded in-memory dead-letter queue.
type Store struct {
	mu      sync.Mutex
	entries []*Entry
	maxSize int
	nextID  int
}

// New creates a Store that holds at most maxSize entries.
// If maxSize is <= 0 it defaults to 100.
func New(maxSize int) *Store {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Store{maxSize: maxSize}
}

// Add appends an entry to the queue, evicting the oldest if the queue is full.
func (s *Store) Add(jobName string, payload []byte, reason string) *Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	e := &Entry{
		ID:        fmt.Sprintf("%d", s.nextID),
		JobName:   jobName,
		Payload:   payload,
		Reason:    reason,
		CreatedAt: time.Now(),
		Attempts:  1,
	}

	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, e)
	return e
}

// All returns a shallow copy of all queued entries, oldest first.
func (s *Store) All() []*Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Remove deletes the entry with the given ID.
func (s *Store) Remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.entries {
		if e.ID == id {
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			return true
		}
	}
	return false
}

// Size returns the current number of queued entries.
func (s *Store) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}
