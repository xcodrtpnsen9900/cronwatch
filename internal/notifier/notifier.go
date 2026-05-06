// Package notifier dispatches alert payloads via webhook with optional
// rate-limiting, deduplication, and throttle controls.
package notifier

import (
	"errors"
	"time"

	"github.com/cronwatch/cronwatch/internal/alert"
	"github.com/cronwatch/cronwatch/internal/throttle"
	"github.com/cronwatch/cronwatch/internal/webhook"
)

// Sender is the interface satisfied by webhook.Client.
type Sender interface {
	Send(payload alert.Payload) error
}

// Notifier wraps a Sender with throttle protection.
type Notifier struct {
	sender   Sender
	throttle *throttle.Throttle
}

// Option configures a Notifier.
type Option func(*Notifier)

// WithThrottle sets a custom throttle on the notifier.
func WithThrottle(th *throttle.Throttle) Option {
	return func(n *Notifier) { n.throttle = th }
}

// New creates a Notifier for the given webhook URL.
// Returns an error if the URL is empty.
func New(webhookURL string, opts ...Option) (*Notifier, error) {
	if webhookURL == "" {
		return nil, errors.New("notifier: webhook URL must not be empty")
	}
	client := webhook.New(webhookURL)
	n := &Notifier{
		sender:   client,
		throttle: throttle.New(time.Minute, 10),
	}
	for _, o := range opts {
		o(n)
	}
	return n, nil
}

// Notify dispatches p if the throttle permits it for p.Job.
// Returns ErrThrottled if the quota is exhausted.
func (n *Notifier) Notify(p alert.Payload) error {
	if !n.throttle.Allow(p.Job) {
		return ErrThrottled
	}
	return n.sender.Send(p)
}

// ErrThrottled is returned when the per-job alert quota is exhausted.
var ErrThrottled = errors.New("notifier: alert throttled")
