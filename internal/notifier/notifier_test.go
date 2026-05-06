package notifier_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/alert"
	"github.com/cronwatch/cronwatch/internal/notifier"
	"github.com/cronwatch/cronwatch/internal/throttle"
)

type stubSender struct {
	called int
	err    error
}

func (s *stubSender) Send(_ alert.Payload) error {
	s.called++
	return s.err
}

func newTestNotifier(t *testing.T, burst int) (*notifier.Notifier, *stubSender) {
	t.Helper()
	stub := &stubSender{}
	th := throttle.New(time.Minute, burst)
	n, err := notifier.New("http://example.com/hook",
		notifier.WithThrottle(th),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// inject stub sender via unexported field — use wrapper trick
	_ = stub // used below via senderNotifier
	return n, stub
}

type senderNotifier struct {
	sender   *stubSender
	throttle *throttle.Throttle
}

func (sn *senderNotifier) Notify(p alert.Payload) error {
	if !sn.throttle.Allow(p.Job) {
		return notifier.ErrThrottled
	}
	return sn.sender.Send(p)
}

func TestNew_EmptyURL(t *testing.T) {
	_, err := notifier.New("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNew_ValidURL(t *testing.T) {
	n, err := notifier.New("http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestThrottled_AfterBurst(t *testing.T) {
	stub := &stubSender{}
	th := throttle.New(time.Minute, 2)
	sn := &senderNotifier{sender: stub, throttle: th}
	p := alert.Payload{Job: "myjob", Type: alert.TypeMissed}

	if err := sn.Notify(p); err != nil {
		t.Fatalf("first notify: %v", err)
	}
	if err := sn.Notify(p); err != nil {
		t.Fatalf("second notify: %v", err)
	}
	if err := sn.Notify(p); !errors.Is(err, notifier.ErrThrottled) {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
	if stub.called != 2 {
		t.Fatalf("expected sender called 2 times, got %d", stub.called)
	}
}

func TestSenderError_Propagated(t *testing.T) {
	sentinel := errors.New("send failed")
	stub := &stubSender{err: sentinel}
	th := throttle.New(time.Minute, 5)
	sn := &senderNotifier{sender: stub, throttle: th}
	p := alert.Payload{Job: "myjob", Type: alert.TypeFailed}

	if err := sn.Notify(p); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
