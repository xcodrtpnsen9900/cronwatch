package overdue

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

type overdueResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
	Jobs      []Entry   `json:"jobs"`
}

// Handler returns an http.HandlerFunc that serialises the current overdue
// set as JSON. The response is sorted by job name for deterministic output.
func Handler(tr *Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries := tr.All()
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Job < entries[j].Job
		})

		resp := overdueResponse{
			Timestamp: time.Now().UTC(),
			Count:     len(entries),
			Jobs:      entries,
		}

		w.Header().Set("Content-Type", "application/json")
		if len(entries) > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}
