package overdue_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/overdue"
)

func TestHandler_EmptyTracker(t *testing.T) {
	tr := overdue.New()
	h := overdue.Handler(tr)

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/overdue", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}

	var resp struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Count != 0 {
		t.Fatalf("want count=0, got %d", resp.Count)
	}
}

func TestHandler_ReturnsOverdueJobs(t *testing.T) {
	tr := overdue.New()
	tr.Mark("alpha", time.Now().Add(-10*time.Minute))
	tr.Mark("beta", time.Now().Add(-5*time.Minute))

	h := overdue.Handler(tr)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/overdue", nil))

	var resp struct {
		Count int              `json:"count"`
		Jobs  []overdue.Entry  `json:"jobs"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Count != 2 {
		t.Fatalf("want 2, got %d", resp.Count)
	}
	if resp.Jobs[0].Job != "alpha" {
		t.Errorf("expected sorted order; first job = %q", resp.Jobs[0].Job)
	}
}

func TestHandler_ContentType(t *testing.T) {
	tr := overdue.New()
	h := overdue.Handler(tr)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/overdue", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("want application/json, got %q", ct)
	}
}
