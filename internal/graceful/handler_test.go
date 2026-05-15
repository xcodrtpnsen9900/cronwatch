package graceful_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/graceful"
)

func serveHandler(c *graceful.Coordinator) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/graceful", nil)
	graceful.Handler(c).ServeHTTP(w, r)
	return w
}

func TestHandler_ContentType(t *testing.T) {
	c := graceful.New(time.Second)
	w := serveHandler(c)
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHandler_NoActiveJobs(t *testing.T) {
	c := graceful.New(time.Second)
	w := serveHandler(c)

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["active_count"].(float64) != 0 {
		t.Error("expected zero active jobs")
	}
	if resp["shutting_down"].(bool) {
		t.Error("expected shutting_down to be false")
	}
}

func TestHandler_ReportsActiveJobs(t *testing.T) {
	c := graceful.New(time.Second)
	c.Acquire("job-x")
	c.Acquire("job-y")
	defer c.Release("job-x")
	defer c.Release("job-y")

	w := serveHandler(c)
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["active_count"].(float64) != 2 {
		t.Errorf("expected 2 active jobs, got %v", resp["active_count"])
	}
}

func TestHandler_ReportsShuttingDown(t *testing.T) {
	c := graceful.New(50 * time.Millisecond)
	c.Acquire("job-z") // keep alive so Shutdown blocks briefly
	go func() {
		time.Sleep(200 * time.Millisecond)
		c.Release("job-z")
	}()
	go c.Shutdown(nil) //nolint:staticcheck — intentional nil ctx for test
	time.Sleep(10 * time.Millisecond) // let shutdown begin

	w := serveHandler(c)
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !resp["shutting_down"].(bool) {
		t.Error("expected shutting_down to be true")
	}
}
