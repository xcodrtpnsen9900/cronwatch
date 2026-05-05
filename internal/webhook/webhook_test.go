package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := webhook.New(ts.URL)
	p := webhook.Payload{
		JobName:   "backup",
		Status:    "missed",
		Message:   "job did not run within the expected window",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := client.Send(p); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.JobName != p.JobName {
		t.Errorf("job_name: got %q, want %q", received.JobName, p.JobName)
	}
	if received.Status != p.Status {
		t.Errorf("status: got %q, want %q", received.Status, p.Status)
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := webhook.New(ts.URL)
	err := client.Send(webhook.Payload{JobName: "test", Status: "failed"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_InvalidURL(t *testing.T) {
	client := webhook.New("http://127.0.0.1:0/no-server")
	err := client.Send(webhook.Payload{JobName: "test", Status: "missed"})
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}

func TestSend_SetsTimestampIfZero(t *testing.T) {
	var received webhook.Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	before := time.Now().UTC()
	client := webhook.New(ts.URL)
	// Send with zero timestamp — client should fill it in.
	_ = client.Send(webhook.Payload{JobName: "heartbeat", Status: "ok"})
	after := time.Now().UTC()

	if received.Timestamp.Before(before) || received.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", received.Timestamp, before, after)
	}
}
