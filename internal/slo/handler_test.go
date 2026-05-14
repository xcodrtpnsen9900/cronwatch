package slo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func serveHandler(t *testing.T, tracker *Tracker, url string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	Handler(tracker)(rr, req)
	return rr
}

func TestHandler_ContentType(t *testing.T) {
	rr := serveHandler(t, newTracker(), "/slo")
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_EmptyTracker(t *testing.T) {
	rr := serveHandler(t, newTracker(), "/slo")
	var snaps []Snapshot
	if err := json.NewDecoder(rr.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(snaps) != 0 {
		t.Fatalf("expected empty, got %v", snaps)
	}
}

func TestHandler_ReturnsAllJobs(t *testing.T) {
	tr := newTracker()
	now := time.Now()
	tr.Record("alpha", now, true)
	tr.Record("beta", now, false)

	rr := serveHandler(t, tr, "/slo")
	var snaps []Snapshot
	if err := json.NewDecoder(rr.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(snaps) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snaps))
	}
}

func TestHandler_FilterByJob(t *testing.T) {
	tr := newTracker()
	now := time.Now()
	tr.Record("alpha", now, true)
	tr.Record("alpha", now, true)
	tr.Record("beta", now, false)

	rr := serveHandler(t, tr, "/slo?job=alpha")
	var snap Snapshot
	if err := json.NewDecoder(rr.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if snap.Job != "alpha" {
		t.Fatalf("expected job alpha, got %s", snap.Job)
	}
	if snap.Total != 2 || snap.Met != 2 {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}
}

func TestHandler_StatusOK(t *testing.T) {
	rr := serveHandler(t, newTracker(), "/slo")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
