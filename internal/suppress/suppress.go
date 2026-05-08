// Package suppress provides time-based alert suppression windows.
// During a suppression window, alerts for a given job are silently dropped,
// allowing operators to schedule maintenance without noise.
package suppress

import (
	"sync"
	"time"
)

// Window represents an active suppression period for a job.
type Window struct {
	Start time.Time
	End   time.Time
}

// Suppressor tracks suppression windows keyed by job name.
type Suppressor struct {
	mu      sync.RWMutex
	windows map[string]Window
	now     func() time.Time
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{
		windows: make(map[string]Window),
		now:     time.Now,
	}
}

// Suppress registers a suppression window for the given job.
// Any existing window for the job is replaced.
func (s *Suppressor) Suppress(job string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.windows[job] = Window{
		Start: now,
		End:   now.Add(duration),
	}
}

// Lift removes the suppression window for the given job immediately.
func (s *Suppressor) Lift(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.windows, job)
}

// IsSuppressed reports whether alerts for the given job are currently suppressed.
func (s *Suppressor) IsSuppressed(job string) bool {
	s.mu.RLock()
	w, ok := s.windows[job]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	now := s.now()
	if now.After(w.End) {
		// Lazily evict expired window.
		s.mu.Lock()
		delete(s.windows, job)
		s.mu.Unlock()
		return false
	}
	return true
}

// Active returns a snapshot of all currently active suppression windows.
func (s *Suppressor) Active() map[string]Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Window, len(s.windows))
	now := s.now()
	for job, w := range s.windows {
		if !now.After(w.End) {
			out[job] = w
		}
	}
	return out
}
