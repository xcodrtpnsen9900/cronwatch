package skiplist

import (
	"encoding/json"
	"net/http"
	"time"
)

type skipRequest struct {
	Job    string `json:"job"`
	Reason string `json:"reason"`
	// DurationSeconds is how long to skip the job.
	DurationSeconds int `json:"duration_seconds"`
}

// Handler returns an http.Handler that exposes the skip-list over HTTP.
//
//   GET  /skiplist        — list all current skip entries
//   POST /skiplist        — register a new skip window
//   DELETE /skiplist?job= — lift a skip window immediately
func Handler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/skiplist", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			entries := s.All()
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(entries); err != nil {
				http.Error(w, "encode error", http.StatusInternalServerError)
			}

		case http.MethodPost:
			var req skipRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if req.DurationSeconds <= 0 {
				http.Error(w, "duration_seconds must be positive", http.StatusBadRequest)
				return
			}
			until := time.Now().Add(time.Duration(req.DurationSeconds) * time.Second)
			if err := s.Skip(req.Job, req.Reason, until); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			job := r.URL.Query().Get("job")
			if job == "" {
				http.Error(w, "job query param required", http.StatusBadRequest)
				return
			}
			s.Lift(job)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
