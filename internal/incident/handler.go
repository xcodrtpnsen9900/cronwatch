package incident

import (
	"encoding/json"
	"net/http"
)

type response struct {
	ID         string  `json:"id"`
	Job        string  `json:"job"`
	Status     Status  `json:"status"`
	Failures   int     `json:"failures"`
	OpenedAt   string  `json:"opened_at"`
	ResolvedAt *string `json:"resolved_at,omitempty"`
}

// Handler returns an http.Handler that exposes all tracked incidents.
// An optional ?job= query parameter filters results to a single job.
func Handler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var incidents []*Incident
		if job := r.URL.Query().Get("job"); job != "" {
			if inc, ok := s.Get(job); ok {
				incidents = []*Incident{inc}
			}
		} else {
			incidents = s.All()
		}

		out := make([]response, 0, len(incidents))
		for _, inc := range incidents {
			r := response{
				ID:       inc.ID,
				Job:      inc.Job,
				Status:   inc.Status,
				Failures: inc.Failures,
				OpenedAt: inc.OpenedAt.UTC().Format("2006-01-02T15:04:05Z"),
			}
			if inc.ResolvedAt != nil {
				s := inc.ResolvedAt.UTC().Format("2006-01-02T15:04:05Z")
				r.ResolvedAt = &s
			}
			out = append(out, r)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})
}
