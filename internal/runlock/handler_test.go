package runlock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveHandler(s *Store) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	Handler(s).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/runlock", nil))
	return rec
}

func TestHandler_ContentType(t *testing.T) {
	rec := serveHandler(newStore(0))
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	rec := serveHandler(newStore(0))
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Count != 0 {
		t.Fatalf("expected count 0, got %d", resp.Count)
	}
	if len(resp.Active) != 0 {
		t.Fatalf("expected empty active list, got %v", resp.Active)
	}
}

func TestHandler_ReportsLockedJobs(t *testing.T) {
	s := newStore(0)
	_ = s.Acquire("alpha")
	_ = s.Acquire("beta")

	rec := serveHandler(s)
	var resp response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Count != 2 {
		t.Fatalf("expected count 2, got %d", resp.Count)
	}
	if resp.Active[0] != "alpha" || resp.Active[1] != "beta" {
		t.Fatalf("unexpected active list: %v", resp.Active)
	}
}

func TestHandler_StatusOK(t *testing.T) {
	rec := serveHandler(newStore(0))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
