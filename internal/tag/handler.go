package tag

import (
	"encoding/json"
	"net/http"
	"sort"
)

type tagResponse struct {
	Job  string   `json:"job"`
	Tags []string `json:"tags"`
}

// Handler returns an HTTP handler that lists all tags for every known job.
// GET /tags          — returns all job→tag mappings
// GET /tags?tag=prod — filters to jobs carrying the given tag
func Handler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("tag")

		var jobs []string
		if filter != "" {
			jobs = s.JobsWithTag(filter)
		} else {
			// collect all jobs that have at least one tag
			s.mu.RLock()
			for job := range s.tags {
				jobs = append(jobs, job)
			}
			s.mu.RUnlock()
		}

		sort.Strings(jobs)

		results := make([]tagResponse, 0, len(jobs))
		for _, job := range jobs {
			tags := s.Get(job)
			sort.Strings(tags)
			results = append(results, tagResponse{Job: job, Tags: tags})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results) //nolint:errcheck
	}
}
