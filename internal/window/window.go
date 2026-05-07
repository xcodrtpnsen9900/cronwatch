// Package window provides a sliding-window counter used to track
// event frequencies over a rolling time period.
package window

import (
	"sync"
	"time"
)

// Counter tracks how many events occurred within a sliding window.
type Counter struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  int
	counts   []int
	times    []time.Time
	current  int
	clock    func() time.Time
}

// New creates a Counter that partitions [size] into [buckets] sub-intervals.
// More buckets give finer resolution at the cost of memory.
func New(size time.Duration, buckets int) *Counter {
	if buckets < 1 {
		buckets = 1
	}
	c := &Counter{
		size:    size,
		buckets: buckets,
		counts:  make([]int, buckets),
		times:   make([]time.Time, buckets),
		clock:   time.Now,
	}
	now := c.clock()
	for i := range c.times {
		c.times[i] = now
	}
	return c
}

// Add records n events at the current time.
func (c *Counter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.advance()
	c.counts[c.current] += n
}

// Total returns the sum of all events within the sliding window.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.advance()
	total := 0
	now := c.clock()
	for i, t := range c.times {
		if now.Sub(t) <= c.size {
			total += c.counts[i]
		}
	}
	return total
}

// Reset zeroes all buckets.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	for i := range c.counts {
		c.counts[i] = 0
		c.times[i] = now
	}
}

// advance moves to a new bucket when the per-bucket interval has elapsed,
// evicting stale data.
func (c *Counter) advance() {
	now := c.clock()
	bucketDur := c.size / time.Duration(c.buckets)
	if now.Sub(c.times[c.current]) >= bucketDur {
		c.current = (c.current + 1) % c.buckets
		c.counts[c.current] = 0
		c.times[c.current] = now
	}
}
