package runbook

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler returns an http.HandlerFunc that exposes the runbook store
// over HTTP.
//
// GET /runbooks          — list all entries
// GET /runbooks?job=foo  — filter to a single job
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		job := strings.TrimSpace(r.URL.Query().Get("job"))
		var entries []Entry

		if job != "" {
			if e, ok := s.Get(job); ok {
				entries = []Entry{e}
			} else {
				entries = []Entry{}
			}
		} else {
			entries = s.All()
		}

		if entries == nil {
			entries = []Entry{}
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"runbooks": entries,
		})
	}
}
