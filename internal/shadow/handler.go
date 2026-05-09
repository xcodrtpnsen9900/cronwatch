package shadow

import (
	"encoding/json"
	"net/http"
	"time"
)

// resultJSON is the JSON representation of a shadow Result.
type resultJSON struct {
	Job      string  `json:"job"`
	DurationMs int64  `json:"duration_ms"`
	Error    *string `json:"error,omitempty"`
}

// Handler returns an http.HandlerFunc that exposes collected shadow results
// as a JSON array. Useful for debugging and observability dashboards.
func Handler(r *Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		results := r.Results()
		out := make([]resultJSON, 0, len(results))
		for _, res := range results {
			rj := resultJSON{
				Job:        res.Job,
				DurationMs: res.Duration.Truncate(time.Millisecond).Milliseconds(),
			}
			if res.Err != nil {
				s := res.Err.Error()
				rj.Error = &s
			}
			out = append(out, rj)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"count":   len(out),
			"results": out,
		})
	}
}
