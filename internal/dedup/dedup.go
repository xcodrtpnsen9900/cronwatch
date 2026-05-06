// Package dedup provides alert deduplication to prevent sending
// identical alerts for the same job within a configurable time window.
package dedup

import (
	"sync"
	"time"
)

// Key uniquely identifies an alert event.
type Key struct {
	Job  string
	Kind string // "missed", "failed", "recovered"
}

// entry tracks the last time an alert was emitted for a given key.
type entry struct {
	lastSent time.Time
}

// Deduplicator suppresses duplicate alerts within a TTL window.
type Deduplicator struct {
	mu      sync.Mutex
	records map[Key]entry
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Deduplicator with the given TTL.
// Alerts with the same Key within TTL are suppressed.
func New(ttl time.Duration) *Deduplicator {
	return &Deduplicator{
		records: make(map[Key]entry),
		ttl:     ttl,
		now:     time.Now,
	}
}

// IsDuplicate returns true if an alert for key was already sent within the TTL.
// If not a duplicate, it records the current time and returns false.
func (d *Deduplicator) IsDuplicate(k Key) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if e, ok := d.records[k]; ok {
		if now.Sub(e.lastSent) < d.ttl {
			return true
		}
	}
	d.records[k] = entry{lastSent: now}
	return false
}

// Reset clears the deduplication state for a specific key,
// allowing the next alert to be sent immediately.
func (d *Deduplicator) Reset(k Key) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.records, k)
}

// Purge removes all entries whose TTL has expired, freeing memory.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, e := range d.records {
		if now.Sub(e.lastSent) >= d.ttl {
			delete(d.records, k)
		}
	}
}
