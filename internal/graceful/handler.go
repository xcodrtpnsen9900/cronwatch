package graceful

import (
	"encoding/json"
	"net/http"
	"sort"
)

type statusResponse struct {
	ShuttingDown bool     `json:"shutting_down"`
	ActiveJobs   []string `json:"active_jobs"`
	ActiveCount  int      `json:"active_count"`
}

// Handler returns an http.HandlerFunc that reports the coordinator's current
// state — whether shutdown is in progress and which jobs are still running.
func Handler(c *Coordinator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		active := c.Active()
		sort.Strings(active)

		shuttingDown := false
		select {
		case <-c.done:
			shuttingDown = true
		default:
		}

		resp := statusResponse{
			ShuttingDown: shuttingDown,
			ActiveJobs:   active,
			ActiveCount:  len(active),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
