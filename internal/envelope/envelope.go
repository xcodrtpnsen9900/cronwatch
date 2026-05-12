// Package envelope wraps outgoing alert payloads with metadata such as
// source host, environment, and schema version before they are dispatched
// via webhook.
package envelope

import (
	"os"
	"time"
)

const SchemaVersion = "1"

// Envelope wraps an alert payload with delivery metadata.
type Envelope struct {
	SchemaVersion string      `json:"schema_version"`
	SentAt        time.Time   `json:"sent_at"`
	Host          string      `json:"host"`
	Environment   string      `json:"environment"`
	Payload       interface{} `json:"payload"`
}

// Builder constructs Envelope values.
type Builder struct {
	env  string
	host string
	now  func() time.Time
}

// New returns a Builder. environment defaults to "production" when empty.
// The host is resolved from os.Hostname; errors are silently ignored.
func New(environment string) *Builder {
	if environment == "" {
		environment = "production"
	}
	host, _ := os.Hostname()
	return &Builder{
		env:  environment,
		host: host,
		now:  time.Now,
	}
}

// withClock replaces the time source; used in tests.
func (b *Builder) withClock(fn func() time.Time) *Builder {
	b.now = fn
	return b
}

// Wrap returns an Envelope containing payload.
func (b *Builder) Wrap(payload interface{}) Envelope {
	return Envelope{
		SchemaVersion: SchemaVersion,
		SentAt:        b.now(),
		Host:          b.host,
		Environment:   b.env,
		Payload:       payload,
	}
}
