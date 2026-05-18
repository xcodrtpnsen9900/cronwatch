package ratelimit

import (
	"encoding/json"
	"net/http"
	"time"
)

type jobSnapshot struct {
	Job       string    `json:"job"`
	Allowed   bool      `json:"allowed"`
	Cooldown  string    `json:"cooldown"`
	LastAlert time.Time `json:"last_alert,omitempty"`
}

// Handler returns an HTTP handler that exposes the current rate-limit state
// for all tracked jobs as a JSON array.
func Handler(rl *RateLimit) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		job := r.URL.Query().Get("job")

		var snaps []jobSnapshot
		for k, entry := range rl.entries {
			if job != "" && k != job {
				continue
			}
			allowed := time.Since(entry.lastAlert) >= rl.cooldown
			snaps = append(snaps, jobSnapshot{
				Job:       k,
				Allowed:   allowed,
				Cooldown:  rl.cooldown.String(),
				LastAlert: entry.lastAlert,
			})
		}
		if snaps == nil {
			snaps = []jobSnapshot{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snaps)
	}
}
