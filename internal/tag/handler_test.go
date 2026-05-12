package tag_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/cronwatch/internal/tag"
)

type tagResp struct {
	Job  string   `json:"job"`
	Tags []string `json:"tags"`
}

func TestHandler_ContentType(t *testing.T) {
	s := tag.New()
	rec := httptest.NewRecorder()
	tag.Handler(s)(rec, httptest.NewRequest(http.MethodGet, "/tags", nil))
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_ReturnsAllJobs(t *testing.T) {
	s := tag.New()
	s.Add("job1", "prod")
	s.Add("job2", "staging")

	rec := httptest.NewRecorder()
	tag.Handler(s)(rec, httptest.NewRequest(http.MethodGet, "/tags", nil))

	var results []tagResp
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(results))
	}
}

func TestHandler_FilterByTag(t *testing.T) {
	s := tag.New()
	s.Add("job1", "prod")
	s.Add("job2", "staging")
	s.Add("job3", "prod")

	rec := httptest.NewRecorder()
	tag.Handler(s)(rec, httptest.NewRequest(http.MethodGet, "/tags?tag=prod", nil))

	var results []tagResp
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 prod jobs, got %d", len(results))
	}
	for _, r := range results {
		found := false
		for _, tg := range r.Tags {
			if tg == "prod" {
				found = true
			}
		}
		if !found {
			t.Fatalf("job %s missing prod tag", r.Job)
		}
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	s := tag.New()
	rec := httptest.NewRecorder()
	tag.Handler(s)(rec, httptest.NewRequest(http.MethodGet, "/tags", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
