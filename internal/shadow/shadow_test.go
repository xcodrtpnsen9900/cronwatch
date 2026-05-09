package shadow_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cronwatch/internal/shadow"
)

func newRunner() *shadow.Runner {
	return shadow.New(log.New(os.Stderr, "", 0))
}

func TestRun_PrimaryErrorReturned(t *testing.T) {
	r := newRunner()
	primaryErr := errors.New("primary failed")
	err := r.Run(context.Background(), "job1",
		func(ctx context.Context) error { return primaryErr },
		func(ctx context.Context) error { return nil },
	)
	if !errors.Is(err, primaryErr) {
		t.Fatalf("expected primary error, got %v", err)
	}
}

func TestRun_ShadowResultRecorded(t *testing.T) {
	r := newRunner()
	shadowErr := errors.New("shadow failed")
	_ = r.Run(context.Background(), "job2",
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return shadowErr
		},
	)
	time.Sleep(50 * time.Millisecond)
	results := r.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Job != "job2" {
		t.Errorf("unexpected job name: %s", results[0].Job)
	}
	if !errors.Is(results[0].Err, shadowErr) {
		t.Errorf("expected shadow error, got %v", results[0].Err)
	}
}

func TestRun_ShadowSuccess(t *testing.T) {
	r := newRunner()
	_ = r.Run(context.Background(), "job3",
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
	)
	time.Sleep(30 * time.Millisecond)
	results := r.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Errorf("expected nil error, got %v", results[0].Err)
	}
	if results[0].Duration < 0 {
		t.Errorf("expected non-negative duration")
	}
}

func TestReset_ClearsResults(t *testing.T) {
	r := newRunner()
	_ = r.Run(context.Background(), "job4",
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
	)
	time.Sleep(30 * time.Millisecond)
	r.Reset()
	if got := r.Results(); len(got) != 0 {
		t.Fatalf("expected 0 results after reset, got %d", len(got))
	}
}
