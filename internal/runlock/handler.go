package runlock

import (
	"encoding/json"
	"net/http"
	"sort"
)

type response struct {
	Active []string `json:"active"`
	Count  int      `json:"count"`
}

// Handler returns an HTTP handler that reports currently locked (running) jobs.
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		active := s.Active()
		sort.Strings(active)

		resp := response{
			Active: active,
			Count:  len(active),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
