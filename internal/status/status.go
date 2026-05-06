// Package status provides an HTTP handler that exposes the current
// state of all monitored cron jobs for health-check and dashboard use.
package status

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/cronwatch/internal/history"
)

// JobStatus summarises the last known state of a single cron job.
type JobStatus struct {
	Name      string     `json:"name"`
	LastRun   *time.Time `json:"last_run,omitempty"`
	LastExit  *int       `json:"last_exit_code,omitempty"`
	Healthy   bool       `json:"healthy"`
	Missed    bool       `json:"missed"`
}

// Provider is the minimal interface the handler needs from the history store.
type Provider interface {
	All(name string) []history.Entry
}

// Handler returns an http.Handler that renders JSON job statuses.
func Handler(jobs []string, h Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statuses := make([]JobStatus, 0, len(jobs))
		for _, name := range jobs {
			entries := h.All(name)
			js := JobStatus{Name: name, Healthy: true}
			if len(entries) > 0 {
				latest := entries[len(entries)-1]
				js.LastRun = &latest.StartedAt
				js.LastExit = &latest.ExitCode
				js.Healthy = latest.ExitCode == 0
				js.Missed = latest.Missed
			}
			statuses = append(statuses, js)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"jobs": statuses,
			"generated_at": time.Now().UTC(),
		})
	})
}
