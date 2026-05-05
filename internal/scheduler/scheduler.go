// Package scheduler parses cron expressions and determines
// whether a job was expected to run within a given time window.
package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Schedule wraps a parsed cron schedule.
type Schedule struct {
	raw      string
	parsed   cron.Schedule
}

// Parse parses a standard 5-field cron expression and returns a Schedule.
func Parse(expr string) (*Schedule, error) {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s, err := p.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("scheduler: invalid cron expression %q: %w", expr, err)
	}
	return &Schedule{raw: expr, parsed: s}, nil
}

// Next returns the next activation time after t.
func (s *Schedule) Next(t time.Time) time.Time {
	return s.parsed.Next(t)
}

// WasExpected reports whether the schedule had at least one activation
// in the half-open interval (since, until].
func (s *Schedule) WasExpected(since, until time.Time) bool {
	next := s.parsed.Next(since)
	return !next.IsZero() && !next.After(until)
}

// String returns the original cron expression.
func (s *Schedule) String() string {
	return s.raw
}
