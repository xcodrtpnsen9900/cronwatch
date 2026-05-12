package budget

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ContentType(t *testing.T) {
	tr := newTracker(t, 0.05)
	h := Handler(tr, []string{})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/budget", nil))
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHandler_EmptyJobs(t *testing.T) {
	tr := newTracker(t, 0.05)
	h := Handler(tr, []string{})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/budget", nil))
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(resp.Jobs))
	}
}

func TestHandler_ReportsThreshold(t *testing.T) {
	tr := newTracker(t, 0.10)
	h := Handler(tr, []string{"job1"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/budget", nil))
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Threshold != 0.10 {
		t.Errorf("expected threshold=0.10, got %v", resp.Threshold)
	}
}

func TestHandler_ExhaustedJob(t *testing.T) {
	tr := newTracker(t, 0.05)
	for i := 0; i < 9; i++ {
		tr.RecordSuccess("jobX")
	}
	tr.RecordFailure("jobX") // 10% > 5%
	h := Handler(tr, []string{"jobX"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/budget", nil))
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(resp.Jobs))
	}
	if !resp.Jobs[0].Exhausted {
		t.Error("expected job to be exhausted")
	}
	if resp.Jobs[0].Failed != 1 {
		t.Errorf("expected 1 failure, got %d", resp.Jobs[0].Failed)
	}
}

func TestHandler_HealthyJob(t *testing.T) {
	tr := newTracker(t, 0.10)
	for i := 0; i < 100; i++ {
		tr.RecordSuccess("jobY")
	}
	h := Handler(tr, []string{"jobY"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/budget", nil))
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Jobs[0].Exhausted {
		t.Error("expected healthy job not exhausted")
	}
	if resp.Jobs[0].Remaining != 1.0 {
		t.Errorf("expected remaining=1.0, got %v", resp.Jobs[0].Remaining)
	}
}
