package slo

import (
	"encoding/json"
	"net/http"
)

// Handler returns an HTTP handler that serves SLO snapshots for all tracked jobs.
// GET /slo          — returns all job snapshots
// GET /slo?job=name — returns the snapshot for a single job
func Handler(t *Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var payload interface{}
		if job := r.URL.Query().Get("job"); job != "" {
			payload = t.Snapshot(job)
		} else {
			snaps := t.All()
			if snaps == nil {
				snaps = []Snapshot{}
			}
			payload = snaps
		}

		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}
}
