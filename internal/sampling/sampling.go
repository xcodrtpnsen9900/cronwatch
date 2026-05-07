// Package sampling provides adaptive sampling for alert events,
// allowing high-frequency alerts to be downsampled to reduce noise.
package sampling

import (
	"sync"
	"time"
)

// Policy defines how sampling is applied to a given job.
type Policy struct {
	// Every is the minimum interval between sampled events per key.
	Every time.Duration
	// MaxBurst is the number of events allowed before sampling kicks in.
	MaxBurst int
}

// DefaultPolicy returns a sensible default sampling policy.
func DefaultPolicy() Policy {
	return Policy{
		Every:    time.Minute,
		MaxBurst: 3,
	}
}

type state struct {
	count    int
	windowAt time.Time
}

// Sampler decides whether an event for a given key should be forwarded.
type Sampler struct {
	mu     sync.Mutex
	policy Policy
	states map[string]*state
	now    func() time.Time
}

// New creates a Sampler with the given policy.
func New(p Policy) *Sampler {
	return &Sampler{
		policy: p,
		states: make(map[string]*state),
		now:    time.Now,
	}
}

// Allow returns true if the event for key should be forwarded.
// It allows up to MaxBurst events, then one event per Every interval.
func (s *Sampler) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	st, ok := s.states[key]
	if !ok {
		s.states[key] = &state{count: 1, windowAt: now}
		return true
	}

	if now.Sub(st.windowAt) >= s.policy.Every {
		st.count = 1
		st.windowAt = now
		return true
	}

	if st.count < s.policy.MaxBurst {
		st.count++
		return true
	}

	return false
}

// Reset clears the sampling state for a key, allowing the next event through.
func (s *Sampler) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, key)
}
