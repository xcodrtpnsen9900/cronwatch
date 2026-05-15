package runbook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveHandler(s *Store, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	Handler(s)(rec, req)
	return rec
}

func TestHandler_ContentType(t *testing.T) {
	rec := serveHandler(New(), "/runbooks")
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("unexpected Content-Type: %q", ct)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	rec := serveHandler(New(), "/runbooks")
	var body map[string][]Entry
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(body["runbooks"]) != 0 {
		t.Errorf("expected empty list")
	}
}

func TestHandler_ReturnsAllEntries(t *testing.T) {
	s := newStore()
	_ = s.Set("alpha", "https://wiki.example.com/alpha", "Alpha job")
	_ = s.Set("beta", "https://wiki.example.com/beta", "")
	rec := serveHandler(s, "/runbooks")
	var body map[string][]Entry
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if len(body["runbooks"]) != 2 {
		t.Errorf("expected 2 entries, got %d", len(body["runbooks"]))
	}
}

func TestHandler_FilterByJob(t *testing.T) {
	s := newStore()
	_ = s.Set("alpha", "https://wiki.example.com/alpha", "Alpha job")
	_ = s.Set("beta", "https://wiki.example.com/beta", "")
	rec := serveHandler(s, "/runbooks?job=alpha")
	var body map[string][]Entry
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if len(body["runbooks"]) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(body["runbooks"]))
	}
	if body["runbooks"][0].Job != "alpha" {
		t.Errorf("unexpected job: %q", body["runbooks"][0].Job)
	}
}

func TestHandler_FilterByJob_NotFound(t *testing.T) {
	s := newStore()
	rec := serveHandler(s, "/runbooks?job=missing")
	var body map[string][]Entry
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if len(body["runbooks"]) != 0 {
		t.Errorf("expected empty list for missing job")
	}
}
