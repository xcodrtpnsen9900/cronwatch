package incident

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveHandler(t *testing.T, s *Store, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	Handler(s).ServeHTTP(rec, req)
	return rec
}

func TestHandler_ContentType(t *testing.T) {
	s := newStore()
	rec := serveHandler(t, s, "/incidents")
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	s := newStore()
	rec := serveHandler(t, s, "/incidents")
	var out []response
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty list, got %d entries", len(out))
	}
}

func TestHandler_ReturnsAllIncidents(t *testing.T) {
	s := newStore()
	s.Open("alpha")
	s.Open("beta")
	rec := serveHandler(t, s, "/incidents")
	var out []response
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 incidents, got %d", len(out))
	}
}

func TestHandler_FilterByJob(t *testing.T) {
	s := newStore()
	s.Open("alpha")
	s.Open("beta")
	rec := serveHandler(t, s, "/incidents?job=alpha")
	var out []response
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 incident, got %d", len(out))
	}
	if out[0].Job != "alpha" {
		t.Fatalf("expected job=alpha, got %s", out[0].Job)
	}
}

func TestHandler_ResolvedAtPresent(t *testing.T) {
	s := newStore()
	s.Open("gamma")
	s.Resolve("gamma")
	rec := serveHandler(t, s, "/incidents?job=gamma")
	var out []response
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatal("expected 1 result")
	}
	if out[0].ResolvedAt == nil {
		t.Fatal("expected resolved_at to be set")
	}
	if out[0].Status != StatusResolved {
		t.Fatalf("expected status=resolved, got %s", out[0].Status)
	}
}
