package budget

import (
	"encoding/json"
	"net/http"
	"time"
)

type responseEntry struct {
	Job       string    `json:"job"`
	Total     int       `json:"total"`
	Failed    int       `json:"failed"`
	Remaining float64   `json:"remaining"`
	Exhausted bool      `json:"exhausted"`
	At        time.Time `json:"at"`
}

type response struct {
	Threshold float64          `json:"threshold"`
	Jobs      []responseEntry  `json:"jobs"`
}

// Handler returns an http.HandlerFunc that reports the current error-budget
// state for all known jobs. Jobs are supplied via the jobs parameter so the
// handler does not need to maintain an additional index.
func Handler(t *Tracker, jobs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries := make([]responseEntry, 0, len(jobs))
		for _, j := range jobs {
			e := t.Snapshot(j)
			entries = append(entries, responseEntry{
				Job:       e.Job,
				Total:     e.Total,
				Failed:    e.Failed,
				Remaining: e.Remaining,
				Exhausted: e.Exhausted,
				At:        e.At,
			})
		}
		resp := response{
			Threshold: t.threshold,
			Jobs:      entries,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
