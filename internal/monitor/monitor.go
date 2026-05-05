package monitor

import (
	"log"
	"sync"
	"time"

	"cronwatch/internal/config"
)

// JobState tracks the last seen execution time for a cron job.
type JobState struct {
	Name        string
	LastSeen    time.Time
	Missed      bool
}

// Monitor watches cron jobs and triggers alerts when runs are missed.
type Monitor struct {
	cfg      *config.Config
	states   map[string]*JobState
	mu       sync.Mutex
	alertFn  func(job string, lastSeen time.Time)
	stopCh   chan struct{}
}

// New creates a Monitor with the provided config and alert callback.
func New(cfg *config.Config, alertFn func(job string, lastSeen time.Time)) *Monitor {
	states := make(map[string]*JobState, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		states[j.Name] = &JobState{
			Name:     j.Name,
			LastSeen: time.Now(),
		}
	}
	return &Monitor{
		cfg:     cfg,
		states:  states,
		alertFn: alertFn,
		stopCh:  make(chan struct{}),
	}
}

// Heartbeat records a successful execution for the named job.
func (m *Monitor) Heartbeat(jobName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.states[jobName]; ok {
		s.LastSeen = time.Now()
		s.Missed = false
		log.Printf("[monitor] heartbeat received for job %q", jobName)
	}
}

// Start begins the periodic check loop.
func (m *Monitor) Start() {
	ticker := time.NewTicker(m.cfg.CheckInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.check()
			case <-m.stopCh:
				return
			}
		}
	}()
}

// Stop halts the monitor loop.
func (m *Monitor) Stop() {
	close(m.stopCh)
}

func (m *Monitor) check() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for _, j := range m.cfg.Jobs {
		s := m.states[j.Name]
		if now.Sub(s.LastSeen) > j.MaxDelay {
			if !s.Missed {
				s.Missed = true
				log.Printf("[monitor] job %q missed — last seen %s ago", j.Name, now.Sub(s.LastSeen).Round(time.Second))
				if m.alertFn != nil {
					go m.alertFn(j.Name, s.LastSeen)
				}
			}
		}
	}
}
