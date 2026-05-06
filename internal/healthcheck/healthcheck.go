// Package healthcheck provides a simple liveness and readiness probe
// endpoint for cronwatch, reporting overall system health based on
// recent job execution state and circuit breaker status.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health state of the service.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
)

// Response is the JSON body returned by the health endpoint.
type Response struct {
	Status    Status            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Checker holds registered named check functions.
type Checker struct {
	mu     sync.RWMutex
	checks map[string]func() error
}

// New creates a new Checker with no registered checks.
func New() *Checker {
	return &Checker{
		checks: make(map[string]func() error),
	}
}

// Register adds a named check function. It is safe to call concurrently.
func (c *Checker) Register(name string, fn func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = fn
}

// Handler returns an http.HandlerFunc that runs all checks and responds
// with 200 OK when healthy or 503 Service Unavailable when degraded.
func (c *Checker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		resp := Response{
			Status:    StatusOK,
			Timestamp: time.Now().UTC(),
			Checks:    make(map[string]string),
		}

		for name, fn := range c.checks {
			if err := fn(); err != nil {
				resp.Checks[name] = err.Error()
				resp.Status = StatusDegraded
			} else {
				resp.Checks[name] = "ok"
			}
		}

		code := http.StatusOK
		if resp.Status == StatusDegraded {
			code = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
