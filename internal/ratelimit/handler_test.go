package ratelimit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func serveHandler(t *testing.T, rl *RateLimit, query string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ratelimit?"+query, nil)
	Handler(rl)(w, req)
	return w
}

func TestHandler_ContentType(t *testing.T) {
	rl := New(time.Minute)
	w := serveHandler(t, rl, "")
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	rl := New(time.Minute)
	w := serveHandler(t, rl, "")
	var snaps []jobSnapshot
	if err := json.NewDecoder(w.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snaps) != 0 {
		t.Fatalf("expected empty, got %d entries", len(snaps))
	}
}

func TestHandler_ReturnsAllJobs(t *testing.T) {
	rl := New(time.Minute)
	rl.Allow("job-a")
	rl.Allow("job-b")

	w := serveHandler(t, rl, "")
	var snaps []jobSnapshot
	if err := json.NewDecoder(w.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snaps) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snaps))
	}
}

func TestHandler_FilterByJob(t *testing.T) {
	rl := New(time.Minute)
	rl.Allow("job-a")
	rl.Allow("job-b")

	w := serveHandler(t, rl, "job=job-a")
	var snaps []jobSnapshot
	if err := json.NewDecoder(w.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snaps) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snaps))
	}
	if snaps[0].Job != "job-a" {
		t.Fatalf("expected job-a, got %s", snaps[0].Job)
	}
}

func TestHandler_AllowedFalseAfterFirstCall(t *testing.T) {
	rl := New(time.Minute)
	rl.Allow("job-x") // first call records last alert

	w := serveHandler(t, rl, "job=job-x")
	var snaps []jobSnapshot
	if err := json.NewDecoder(w.Body).Decode(&snaps); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snaps) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snaps))
	}
	if snaps[0].Allowed {
		t.Fatal("expected allowed=false immediately after first call")
	}
}
