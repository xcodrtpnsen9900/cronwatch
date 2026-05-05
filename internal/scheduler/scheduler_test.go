package scheduler_test

import (
	"testing"
	"time"

	"github.com/user/cronwatch/internal/scheduler"
)

func mustParse(t *testing.T, expr string) *scheduler.Schedule {
	t.Helper()
	s, err := scheduler.Parse(expr)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", expr, err)
	}
	return s
}

func TestParse_Valid(t *testing.T) {
	_, err := scheduler.Parse("*/5 * * * *")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := scheduler.Parse("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestNext(t *testing.T) {
	s := mustParse(t, "0 * * * *") // top of every hour
	base := time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)
	next := s.Next(base)
	want := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("Next() = %v, want %v", next, want)
	}
}

func TestWasExpected_True(t *testing.T) {
	s := mustParse(t, "*/10 * * * *") // every 10 minutes
	since := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 12, 15, 0, 0, time.UTC)
	if !s.WasExpected(since, until) {
		t.Error("WasExpected() = false, want true")
	}
}

func TestWasExpected_False(t *testing.T) {
	s := mustParse(t, "0 0 * * *") // midnight only
	since := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)
	if s.WasExpected(since, until) {
		t.Error("WasExpected() = true, want false")
	}
}

func TestString(t *testing.T) {
	expr := "*/5 * * * *"
	s := mustParse(t, expr)
	if s.String() != expr {
		t.Errorf("String() = %q, want %q", s.String(), expr)
	}
}
