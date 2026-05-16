// Package baseline tracks the expected duration of cron jobs and flags
// runs that deviate significantly from the historical average.
package baseline

import (
	"errors"
	"math"
	"sync"
	"time"
)

// Entry holds a single duration sample for a job.
type Entry struct {
	Job      string
	Duration time.Duration
	RecordedAt time.Time
}

// Snapshot summarises the baseline for one job.
type Snapshot struct {
	Job      string
	Samples  int
	MeanMs   float64
	StdDevMs float64
}

// Tracker accumulates duration samples and detects anomalies.
type Tracker struct {
	mu      sync.Mutex
	samples map[string][]float64
	maxPer  int
	sigma   float64 // number of standard deviations before flagging
}

// New returns a Tracker. maxPer caps stored samples per job; sigma sets the
// anomaly threshold (e.g. 2.0 means flag if > mean + 2σ).
func New(maxPer int, sigma float64) (*Tracker, error) {
	if maxPer <= 0 {
		return nil, errors.New("baseline: maxPer must be positive")
	}
	if sigma <= 0 {
		return nil, errors.New("baseline: sigma must be positive")
	}
	return &Tracker{
		samples: make(map[string][]float64),
		maxPer:  maxPer,
		sigma:   sigma,
	}, nil
}

// Record adds a duration sample for job.
func (t *Tracker) Record(job string, d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	ms := float64(d.Milliseconds())
	buf := append(t.samples[job], ms)
	if len(buf) > t.maxPer {
		buf = buf[len(buf)-t.maxPer:]
	}
	t.samples[job] = buf
}

// IsAnomaly returns true when d exceeds mean + sigma*stddev for job.
// Returns false if fewer than 2 samples exist (insufficient data).
func (t *Tracker) IsAnomaly(job string, d time.Duration) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	buf := t.samples[job]
	if len(buf) < 2 {
		return false
	}
	mean, std := stats(buf)
	return float64(d.Milliseconds()) > mean+t.sigma*std
}

// Snapshot returns a summary for job, or zero-value if unknown.
func (t *Tracker) Snapshot(job string) Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	buf := t.samples[job]
	if len(buf) == 0 {
		return Snapshot{Job: job}
	}
	mean, std := stats(buf)
	return Snapshot{Job: job, Samples: len(buf), MeanMs: mean, StdDevMs: std}
}

func stats(buf []float64) (mean, std float64) {
	for _, v := range buf {
		mean += v
	}
	mean /= float64(len(buf))
	for _, v := range buf {
		diff := v - mean
		std += diff * diff
	}
	std = math.Sqrt(std / float64(len(buf)))
	return
}
