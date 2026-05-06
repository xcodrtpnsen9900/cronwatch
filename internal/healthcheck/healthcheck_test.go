package healthcheck_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/cronwatch/internal/healthcheck"
)

func TestHandler_AllChecksPass(t *testing.T) {
	c := healthcheck.New()
	c.Register("db", func() error { return nil })
	c.Register("webhook", func() error { return nil })

	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp healthcheck.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != healthcheck.StatusOK {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
	if resp.Checks["db"] != "ok" || resp.Checks["webhook"] != "ok" {
		t.Errorf("unexpected checks: %v", resp.Checks)
	}
}

func TestHandler_OneCheckFails(t *testing.T) {
	c := healthcheck.New()
	c.Register("db", func() error { return nil })
	c.Register("webhook", func() error { return errors.New("connection refused") })

	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var resp healthcheck.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != healthcheck.StatusDegraded {
		t.Errorf("expected degraded, got %s", resp.Status)
	}
	if resp.Checks["webhook"] != "connection refused" {
		t.Errorf("expected error message in checks, got %q", resp.Checks["webhook"])
	}
}

func TestHandler_NoChecks(t *testing.T) {
	c := healthcheck.New()

	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with no checks, got %d", rec.Code)
	}
}

func TestHandler_ContentType(t *testing.T) {
	c := healthcheck.New()

	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestHandler_TimestampPresent(t *testing.T) {
	c := healthcheck.New()

	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	var resp healthcheck.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
