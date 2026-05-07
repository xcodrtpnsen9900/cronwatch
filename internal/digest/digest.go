// Package digest provides a periodic summary report of cron job
// activity, aggregating missed, failed, and recovered events over
// a configurable window and emitting them as a single webhook payload.
package digest

import (
	"fmt"
	"sync"
	"time"
)

// Entry records a single job event for inclusion in a digest.
type Entry struct {
	Job       string
	Status    string // "missed", "failed", "recovered"
	OccuredAt time.Time
}

// Report is the payload sent by Flush.
type Report struct {
	Window    time.Duration
	GeneratedAt time.Time
	Entries   []Entry
	Missed    int
	Failed    int
	Recovered int
}

// Sender is the interface used to deliver a digest report.
type Sender interface {
	SendDigest(r Report) error
}

// Digest accumulates job events and periodically flushes them.
type Digest struct {
	mu      sync.Mutex
	entries []Entry
	window  time.Duration
	sender  Sender
	stop    chan struct{}
	wg      sync.WaitGroup
}

// New creates a Digest that flushes to sender every window duration.
// Call Start to begin the flush ticker.
func New(window time.Duration, sender Sender) *Digest {
	return &Digest{
		window: window,
		sender: sender,
		stop:   make(chan struct{}),
	}
}

// Record adds a job event to the current accumulation window.
func (d *Digest) Record(job, status string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = append(d.entries, Entry{
		Job:       job,
		Status:    status,
		OccuredAt: time.Now(),
	})
}

// Flush drains accumulated entries into a Report and sends it.
// It is a no-op when there are no entries.
func (d *Digest) Flush() error {
	d.mu.Lock()
	entries := d.entries
	d.entries = nil
	d.mu.Unlock()

	if len(entries) == 0 {
		return nil
	}

	r := Report{
		Window:      d.window,
		GeneratedAt: time.Now(),
		Entries:     entries,
	}
	for _, e := range entries {
		switch e.Status {
		case "missed":
			r.Missed++
		case "failed":
			r.Failed++
		case "recovered":
			r.Recovered++
		}
	}
	return d.sender.SendDigest(r)
}

// Start begins the background flush ticker.
func (d *Digest) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		ticker := time.NewTicker(d.window)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := d.Flush(); err != nil {
					_ = fmt.Errorf("digest flush: %w", err)
				}
			case <-d.stop:
				return
			}
		}
	}()
}

// Stop halts the background ticker and performs a final flush.
func (d *Digest) Stop() error {
	close(d.stop)
	d.wg.Wait()
	return d.Flush()
}
