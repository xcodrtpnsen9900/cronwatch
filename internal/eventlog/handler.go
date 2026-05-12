package eventlog

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that exposes the event log as JSON.
// Optional query params: level=info|warn|error, job=<name>.
//
// Responses:
//
//	200 OK            – JSON array of matching events (may be empty)
//	405 Method Not Allowed – non-GET request
//	500 Internal Server Error – JSON encoding failure
func Handler(l *Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		level := Level(q.Get("level"))
		job := q.Get("job")

		var events []Event
		if level != "" || job != "" {
			events = l.Filter(level, job)
		} else {
			events = l.All()
		}

		if events == nil {
			events = []Event{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	})
}
