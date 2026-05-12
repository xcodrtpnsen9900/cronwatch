package triage

import (
	"encoding/json"
	"net/http"
	"time"
)

type jobLevel struct {
	Job   string `json:"job"`
	Level string `json:"level"`
}

type snapshot struct {
	Timestamp time.Time  `json:"timestamp"`
	Jobs      []jobLevel `json:"jobs"`
}

// Handler returns an HTTP handler that reports the current triage level
// for all tracked jobs.
func Handler(c *Classifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		jobs := make([]jobLevel, 0, len(c.states))
		for job := range c.states {
			jobs = append(jobs, jobLevel{
				Job:   job,
				Level: c.Level(job).String(),
			})
		}
		c.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapshot{
			Timestamp: time.Now().UTC(),
			Jobs:      jobs,
		})
	}
}
