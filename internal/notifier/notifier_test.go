package notifier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/notifier"
	"github.com/example/cronwatch/internal/webhook"
)

// mockSender captures the last payload sent and optionally returns an error.
type mockSender struct {
	last webhook.Payload
	err  error
}

func (m *mockSender) Send(_ context.Context, p webhook.Payload) error {
	m.last = p
	return m.err
}

func newTestNotifier(t *testing.T, ms *mockSender) *notifier.Notifier {
	t.Helper()
	n, err := notifier.New("http://example.com/hook", nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Inject mock sender via the exported field — we use the package-level
	// constructor and rely on the Sender interface for injection in tests.
	_ = ms // sender injection tested via New; direct field access omitted
	return n
}

func TestNew_EmptyURL(t *testing.T) {
	_, err := notifier.New("", nil)
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNew_ValidURL(t *testing.T) {
	n, err := notifier.New("http://localhost/hook", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

// senderNotifier builds a Notifier whose internal sender is replaced with ms.
func senderNotifier(t *testing.T, ms *mockSender) *notifier.Notifier {
	t.Helper()
	n, err := notifier.NewWithSender("http://example.com/hook", ms, nil)
	if err != nil {
		t.Fatalf("NewWithSender: %v", err)
	}
	return n
}

func TestMissed_Delegates(t *testing.T) {
	ms := &mockSender{}
	n := senderNotifier(t, ms)
	if err := n.Missed(context.Background(), "backup", time.Now()); err != nil {
		t.Fatalf("Missed: %v", err)
	}
	if ms.last.Job != "backup" {
		t.Errorf("expected job=backup, got %q", ms.last.Job)
	}
}

func TestFailed_Delegates(t *testing.T) {
	ms := &mockSender{}
	n := senderNotifier(t, ms)
	if err := n.Failed(context.Background(), "sync", time.Now(), "exit 1"); err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if ms.last.Job != "sync" {
		t.Errorf("expected job=sync, got %q", ms.last.Job)
	}
}

func TestSendError_Propagates(t *testing.T) {
	ms := &mockSender{err: errors.New("network down")}
	n := senderNotifier(t, ms)
	err := n.Recovered(context.Background(), "cleanup", time.Now())
	if err == nil {
		t.Fatal("expected error to propagate, got nil")
	}
}
