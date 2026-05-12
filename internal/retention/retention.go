// Package retention enforces data retention policies by evicting job
// history entries older than a configured maximum age.
package retention

import (
	"sync"
	"time"
)

// Entry represents a retainable record with a timestamp and job identifier.
type Entry struct {
	Job       string
	Timestamp time.Time
}

// Policy defines how long entries should be retained.
type Policy struct {
	MaxAge time.Duration
}

// DefaultPolicy retains entries for 30 days.
var DefaultPolicy = Policy{
	MaxAge: 30 * 24 * time.Hour,
}

// Store manages entries subject to a retention policy.
type Store struct {
	mu      sync.Mutex
	entries []Entry
	policy  Policy
	now     func() time.Time
}

// New creates a Store with the given retention policy.
func New(p Policy) *Store {
	return &Store{
		policy: p,
		now:    time.Now,
	}
}

// Add appends an entry and immediately evicts any entries that exceed MaxAge.
func (s *Store) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	s.evict()
}

// All returns a copy of all retained entries.
func (s *Store) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Evict removes entries older than MaxAge and returns the count removed.
func (s *Store) Evict() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	before := len(s.entries)
	s.evict()
	return before - len(s.entries)
}

// evict must be called with mu held.
func (s *Store) evict() {
	cutoff := s.now().Add(-s.policy.MaxAge)
	filtered := s.entries[:0]
	for _, e := range s.entries {
		if !e.Timestamp.Before(cutoff) {
			filtered = append(filtered, e)
		}
	}
	s.entries = filtered
}
