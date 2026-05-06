// Package notifier wires together alert building and webhook delivery.
package notifier

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/cronwatch/internal/alert"
	"github.com/example/cronwatch/internal/webhook"
)

// Sender is the interface used to deliver a payload to a webhook endpoint.
type Sender interface {
	Send(ctx context.Context, payload webhook.Payload) error
}

// Notifier composes an alert builder with a webhook sender.
type Notifier struct {
	builder *alert.Builder
	sender  Sender
	logger  *log.Logger
}

// New creates a Notifier for the given webhook URL.
// If logger is nil a default logger is used.
func New(webhookURL string, logger *log.Logger) (*Notifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("notifier: webhook URL must not be empty")
	}
	sender, err := webhook.New(webhookURL)
	if err != nil {
		return nil, fmt.Errorf("notifier: %w", err)
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Notifier{
		builder: alert.NewBuilder(webhookURL),
		sender:  sender,
		logger:  logger,
	}, nil
}

// Missed sends a "missed" alert for the named job whose last expected run was at t.
func (n *Notifier) Missed(ctx context.Context, jobName string, t time.Time) error {
	payload := n.builder.Missed(jobName, t)
	if err := n.sender.Send(ctx, payload); err != nil {
		n.logger.Printf("notifier: missed alert for %q: %v", jobName, err)
		return err
	}
	return nil
}

// Failed sends a "failed" alert for the named job with the provided error detail.
func (n *Notifier) Failed(ctx context.Context, jobName string, t time.Time, detail string) error {
	payload := n.builder.Failed(jobName, t, detail)
	if err := n.sender.Send(ctx, payload); err != nil {
		n.logger.Printf("notifier: failed alert for %q: %v", jobName, err)
		return err
	}
	return nil
}

// Recovered sends a "recovered" alert indicating the job is healthy again.
func (n *Notifier) Recovered(ctx context.Context, jobName string, t time.Time) error {
	payload := n.builder.Recovered(jobName, t)
	if err := n.sender.Send(ctx, payload); err != nil {
		n.logger.Printf("notifier: recovered alert for %q: %v", jobName, err)
		return err
	}
	return nil
}
