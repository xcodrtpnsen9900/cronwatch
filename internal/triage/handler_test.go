package triage_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/triage"
)

func serveHandler(c *triage.Classifier) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/triage", nil)
	triage.Handler(c).ServeHTTP(rec, req)
	return rec
}

func TestHandler_ContentType(t *testing.T) {
	c := triage.New(triage.DefaultPolicy())
	rec := serveHandler(c)
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_EmptyClassifier(t *testing.T) {
	c := triage.New(triage.DefaultPolicy())
	rec := serveHandler(c)
	var snap struct {
		Jobs []interface{} `json:"jobs"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snap.Jobs) != 0 {
		t.Fatalf("expected empty jobs, got %d", len(snap.Jobs))
	}
}

func TestHandler_ReportsJobLevel(t *testing.T) {
	c := triage.New(triage.Policy{ErrorAfter: 2, CritAfter: 4, Window: time.Hour})
	c.Record("backup", time.Now())

	rec := serveHandler(c)
	var snap struct {
		Jobs []struct {
			Job   string `json:"job"`
			Level string `json:"level"`
		} `json:"jobs"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(snap.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(snap.Jobs))
	}
	if snap.Jobs[0].Job != "backup" {
		t.Errorf("unexpected job name %q", snap.Jobs[0].Job)
	}
	if snap.Jobs[0].Level != "warn" {
		t.Errorf("expected warn, got %q", snap.Jobs[0].Level)
	}
}
