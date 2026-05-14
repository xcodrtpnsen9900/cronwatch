package replay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func serveHandler(s *Store, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	Handler(s)(rec, req)
	return rec
}

func TestHandler_ContentType(t *testing.T) {
	s := newStore()
	rec := serveHandler(s, "/replay")
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	s := newStore()
	rec := serveHandler(s, "/replay")

	var out []entryJSON
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d entries", len(out))
	}
}

func TestHandler_ReturnsAllEntries(t *testing.T) {
	s := newStore()
	s.Enqueue("backup", "missed", time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))
	s.Enqueue("sync", "failed", time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))

	rec := serveHandler(s, "/replay")

	var out []entryJSON
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestHandler_FilterByJob(t *testing.T) {
	s := newStore()
	s.Enqueue("backup", "missed", t0)
	s.Enqueue("sync", "failed", t0)

	rec := serveHandler(s, "/replay?job=backup")

	var out []entryJSON
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].JobName != "backup" {
		t.Errorf("expected backup, got %q", out[0].JobName)
	}
}
