package replay

import (
	"encoding/json"
	"net/http"
)

type entryJSON struct {
	JobName   string `json:"job_name"`
	Reason    string `json:"reason"`
	Scheduled string `json:"scheduled"`
	QueuedAt  string `json:"queued_at"`
}

// Handler returns an http.HandlerFunc that exposes the current replay queue
// as JSON. Supports an optional ?job= query parameter to filter by job name.
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		job := r.URL.Query().Get("job")

		var entries []Entry
		if job != "" {
			// Return a read-only view without draining.
			s.mu.Lock()
			entries = make([]Entry, len(s.entries[job]))
			copy(entries, s.entries[job])
			s.mu.Unlock()
		} else {
			entries = s.All()
		}

		out := make([]entryJSON, len(entries))
		for i, e := range entries {
			out[i] = entryJSON{
				JobName:   e.JobName,
				Reason:    e.Reason,
				Scheduled: e.Scheduled.UTC().Format("2006-01-02T15:04:05Z"),
				QueuedAt:  e.QueuedAt.UTC().Format("2006-01-02T15:04:05Z"),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}
