package status_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/history"
	"github.com/example/cronwatch/internal/status"
)

type fakeHistory struct {
	data map[string][]history.Entry
}

func (f *fakeHistory) All(name string) []history.Entry {
	return f.data[name]
}

func TestHandler_NoEntries(t *testing.T) {
	h := &fakeHistory{data: map[string][]history.Entry{}}
	rec := httptest.NewRecorder()
	status.Handler([]string{"backup"}, h).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	jobs := body["jobs"].([]interface{})
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	job := jobs[0].(map[string]interface{})
	if job["name"] != "backup" {
		t.Errorf("unexpected name: %v", job["name"])
	}
	if job["healthy"] != true {
		t.Errorf("expected healthy=true for job with no history")
	}
}

func TestHandler_FailedJob(t *testing.T) {
	now := time.Now().UTC()
	exit := 1
	h := &fakeHistory{
		data: map[string][]history.Entry{
			"nightly": {{StartedAt: now, ExitCode: exit, Missed: false}},
		},
	}
	rec := httptest.NewRecorder()
	status.Handler([]string{"nightly"}, h).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)
	job := body["jobs"].([]interface{})[0].(map[string]interface{})
	if job["healthy"] != false {
		t.Errorf("expected healthy=false for exit code 1")
	}
}

func TestHandler_ContentType(t *testing.T) {
	h := &fakeHistory{data: map[string][]history.Entry{}}
	rec := httptest.NewRecorder()
	status.Handler([]string{}, h).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
