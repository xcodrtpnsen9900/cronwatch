package monitor

import (
	"sync/atomic"
	"testing"
	"time"

	"cronwatch/internal/config"
)

func makeConfig(maxDelay time.Duration) *config.Config {
	return &config.Config{
		WebhookURL:    "http://example.com/hook",
		CheckInterval: 50 * time.Millisecond,
		Jobs: []config.Job{
			{Name: "backup", Schedule: "@daily", MaxDelay: maxDelay},
		},
	}
}

func TestHeartbeat_ResetsState(t *testing.T) {
	cfg := makeConfig(200 * time.Millisecond)
	m := New(cfg, nil)

	// Manually age the last-seen time.
	m.mu.Lock()
	m.states["backup"].LastSeen = time.Now().Add(-1 * time.Hour)
	m.states["backup"].Missed = true
	m.mu.Unlock()

	m.Heartbeat("backup")

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.states["backup"].Missed {
		t.Error("expected Missed to be false after heartbeat")
	}
	if time.Since(m.states["backup"].LastSeen) > time.Second {
		t.Error("expected LastSeen to be updated to roughly now")
	}
}

func TestCheck_TriggersAlert(t *testing.T) {
	cfg := makeConfig(10 * time.Millisecond)
	var alertCount int32
	alertFn := func(job string, _ time.Time) {
		if job != "backup" {
			t.Errorf("unexpected job name: %s", job)
		}
		atomic.AddInt32(&alertCount, 1)
	}

	m := New(cfg, alertFn)
	// Age the heartbeat so the job appears missed.
	m.mu.Lock()
	m.states["backup"].LastSeen = time.Now().Add(-1 * time.Hour)
	m.mu.Unlock()

	m.check()
	time.Sleep(20 * time.Millisecond) // let goroutine fire

	if atomic.LoadInt32(&alertCount) != 1 {
		t.Errorf("expected 1 alert, got %d", alertCount)
	}
}

func TestCheck_NoDoubleAlert(t *testing.T) {
	cfg := makeConfig(10 * time.Millisecond)
	var alertCount int32
	m := New(cfg, func(_ string, _ time.Time) {
		atomic.AddInt32(&alertCount, 1)
	})
	m.mu.Lock()
	m.states["backup"].LastSeen = time.Now().Add(-1 * time.Hour)
	m.mu.Unlock()

	m.check()
	m.check() // second call should not fire again
	time.Sleep(20 * time.Millisecond)

	if atomic.LoadInt32(&alertCount) != 1 {
		t.Errorf("expected exactly 1 alert, got %d", alertCount)
	}
}

func TestStartStop(t *testing.T) {
	cfg := makeConfig(10 * time.Millisecond)
	m := New(cfg, nil)
	m.Start()
	time.Sleep(30 * time.Millisecond)
	m.Stop() // should not block or panic
}
