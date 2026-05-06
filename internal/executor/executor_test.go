package executor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwatch/internal/alert"
	"github.com/example/cronwatch/internal/executor"
	"github.com/example/cronwatch/internal/history"
	"github.com/example/cronwatch/internal/notifier"
)

func setup(t *testing.T, webhookURL string) (*executor.Executor, *history.History) {
	t.Helper()
	h := history.New()
	n, err := notifier.New(webhookURL)
	if err != nil {
		t.Fatalf("notifier.New: %v", err)
	}
	ab := alert.NewBuilder("cronwatch")
	ex := executor.New("test-job", "echo", []string{"hello"}, h, n, ab)
	return ex, h
}

func TestRun_SuccessRecordsEntry(t *testing.T) {
	var called bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ex, h := setup(t, srv.URL)
	if err := ex.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no webhook call on success")
	}
	if _, ok := h.Latest("test-job"); !ok {
		t.Error("expected history entry to be recorded")
	}
}

func TestRun_FailureSendsAlert(t *testing.T) {
	var called bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := history.New()
	n, _ := notifier.New(srv.URL)
	ab := alert.NewBuilder("cronwatch")
	ex := executor.New("fail-job", "false", nil, h, n, ab)

	_ = ex.Run(context.Background())
	if !called {
		t.Error("expected webhook to be called on failure")
	}
}

func TestRun_RecoverySendsAlert(t *testing.T) {
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := history.New()
	n, _ := notifier.New(srv.URL)
	ab := alert.NewBuilder("cronwatch")

	// First run: failure
	exFail := executor.New("recover-job", "false", nil, h, n, ab)
	_ = exFail.Run(context.Background())

	// Second run: success after failure — should trigger recovery alert
	exOK := executor.New("recover-job", "echo", []string{"ok"}, h, n, ab)
	_ = exOK.Run(context.Background())

	if callCount < 2 {
		t.Errorf("expected at least 2 webhook calls (failure + recovery), got %d", callCount)
	}
}
