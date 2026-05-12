package deadletter

import (
	"encoding/json"
	"net/http"
)

type entryResponse struct {
	ID        string `json:"id"`
	JobName   string `json:"job_name"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
	Attempts  int    `json:"attempts"`
}

// Handler returns an http.HandlerFunc that exposes the dead-letter queue.
// DELETE /?id=<id> removes a single entry; GET returns all entries.
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "missing id", http.StatusBadRequest)
				return
			}
			if !s.Remove(id) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			entries := s.All()
			resp := make([]entryResponse, 0, len(entries))
			for _, e := range entries {
				resp = append(resp, entryResponse{
					ID:        e.ID,
					JobName:   e.JobName,
					Reason:    e.Reason,
					CreatedAt: e.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
					Attempts:  e.Attempts,
				})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"count":   len(resp),
				"entries": resp,
			})
		}
	}
}
