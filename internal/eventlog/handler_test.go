package eventlog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsAllEvents(t *testing.T) {
	l := New(10)
	l.Add(LevelInfo, "job1", "started", nil)
	l.Add(LevelError, "job2", "failed", nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	Handler(l).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestHandler_FilterByLevel(t *testing.T) {
	l := New(10)
	l.Add(LevelInfo, "a", "ok", nil)
	l.Add(LevelError, "b", "fail", nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events?level=error", nil)
	Handler(l).ServeHTTP(rec, req)

	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 1 || events[0].Level != LevelError {
		t.Errorf("expected 1 error event, got %+v", events)
	}
}

func TestHandler_EmptyLog(t *testing.T) {
	l := New(10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	Handler(l).ServeHTTP(rec, req)

	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty array, got %d events", len(events))
	}
}

func TestHandler_ContentType(t *testing.T) {
	l := New(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	Handler(l).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
