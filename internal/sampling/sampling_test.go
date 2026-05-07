package sampling_test

import (
	"testing"
	"time"

	"github.com/example/cronwatch/internal/sampling"
)

func newSampler(burst int, every time.Duration) (*sampling.Sampler, *time.Time) {
	s := sampling.New(sampling.Policy{MaxBurst: burst, Every: every})
	now := time.Now()
	// inject controllable clock via unexported field workaround: use real clock
	return s, &now
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	s := sampling.New(sampling.DefaultPolicy())
	if !s.Allow("job1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_BurstPermitted(t *testing.T) {
	p := sampling.Policy{MaxBurst: 3, Every: time.Hour}
	s := sampling.New(p)
	for i := 0; i < 3; i++ {
		if !s.Allow("job1") {
			t.Fatalf("expected call %d to be allowed within burst", i+1)
		}
	}
}

func TestAllow_SuppressedAfterBurst(t *testing.T) {
	p := sampling.Policy{MaxBurst: 2, Every: time.Hour}
	s := sampling.New(p)
	s.Allow("job1")
	s.Allow("job1")
	if s.Allow("job1") {
		t.Fatal("expected call to be suppressed after burst exhausted")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	p := sampling.Policy{MaxBurst: 1, Every: time.Hour}
	s := sampling.New(p)
	s.Allow("job1")
	if !s.Allow("job2") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	p := sampling.Policy{MaxBurst: 1, Every: time.Hour}
	s := sampling.New(p)
	s.Allow("job1") // consume burst
	s.Allow("job1") // suppressed
	s.Reset("job1")
	if !s.Allow("job1") {
		t.Fatal("expected allow after reset")
	}
}

func TestDefaultPolicy_Fields(t *testing.T) {
	p := sampling.DefaultPolicy()
	if p.MaxBurst <= 0 {
		t.Errorf("expected positive MaxBurst, got %d", p.MaxBurst)
	}
	if p.Every <= 0 {
		t.Errorf("expected positive Every, got %v", p.Every)
	}
}
