package alert

import (
	"fmt"
	"time"
)

// Type represents the kind of alert being sent.
type Type string

const (
	TypeMissed    Type = "missed"
	TypeFailed    Type = "failed"
	TypeRecovered Type = "recovered"
)

// Payload is the structured alert sent to a webhook.
type Payload struct {
	JobName   string    `json:"job_name"`
	AlertType Type      `json:"alert_type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details,omitempty"`
}

// Builder constructs alert payloads for cron job events.
type Builder struct {
	now func() time.Time
}

// NewBuilder returns a Builder using the real clock.
func NewBuilder() *Builder {
	return &Builder{now: time.Now}
}

// Missed returns a Payload indicating a job did not run within its expected window.
func (b *Builder) Missed(jobName string, expectedAt time.Time) Payload {
	return Payload{
		JobName:   jobName,
		AlertType: TypeMissed,
		Message:   fmt.Sprintf("Job %q missed its scheduled run", jobName),
		Timestamp: b.now(),
		Details:   fmt.Sprintf("expected at %s", expectedAt.Format(time.RFC3339)),
	}
}

// Failed returns a Payload indicating a job reported a failure.
func (b *Builder) Failed(jobName string, reason string) Payload {
	return Payload{
		JobName:   jobName,
		AlertType: TypeFailed,
		Message:   fmt.Sprintf("Job %q failed", jobName),
		Timestamp: b.now(),
		Details:   reason,
	}
}

// Recovered returns a Payload indicating a previously alerting job is healthy again.
func (b *Builder) Recovered(jobName string) Payload {
	return Payload{
		JobName:   jobName,
		AlertType: TypeRecovered,
		Message:   fmt.Sprintf("Job %q has recovered", jobName),
		Timestamp: b.now(),
	}
}

// MissedWithDuration returns a Payload for a missed job that includes how long
// ago the job was expected to run, providing more context in the alert details.
func (b *Builder) MissedWithDuration(jobName string, expectedAt time.Time) Payload {
	now := b.now()
	overdue := now.Sub(expectedAt).Truncate(time.Second)
	return Payload{
		JobName:   jobName,
		AlertType: TypeMissed,
		Message:   fmt.Sprintf("Job %q missed its scheduled run", jobName),
		Timestamp: now,
		Details:   fmt.Sprintf("expected at %s (overdue by %s)", expectedAt.Format(time.RFC3339), overdue),
	}
}
