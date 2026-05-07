package digest

import (
	"sync"
	"testing"
	"time"
)

type mockSender struct {
	mu      sync.Mutex
	reports []Report
	err     error
}

func (m *mockSender) SendDigest(r Report) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.reports = append(m.reports, r)
	return nil
}

func (m *mockSender) last() (Report, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.reports) == 0 {
		return Report{}, false
	}
	return m.reports[len(m.reports)-1], true
}

func TestFlush_NoEntries_NoSend(t *testing.T) {
	s := &mockSender{}
	d := New(time.Minute, s)
	if err := d.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.last(); ok {
		t.Fatal("expected no report to be sent when no entries")
	}
}

func TestFlush_SendsReport(t *testing.T) {
	s := &mockSender{}
	d := New(time.Minute, s)
	d.Record("job-a", "missed")
	d.Record("job-b", "failed")
	d.Record("job-b", "recovered")

	if err := d.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r, ok := s.last()
	if !ok {
		t.Fatal("expected a report")
	}
	if r.Missed != 1 {
		t.Errorf("missed: want 1, got %d", r.Missed)
	}
	if r.Failed != 1 {
		t.Errorf("failed: want 1, got %d", r.Failed)
	}
	if r.Recovered != 1 {
		t.Errorf("recovered: want 1, got %d", r.Recovered)
	}
	if len(r.Entries) != 3 {
		t.Errorf("entries: want 3, got %d", len(r.Entries))
	}
}

func TestFlush_DrainsBetweenCalls(t *testing.T) {
	s := &mockSender{}
	d := New(time.Minute, s)
	d.Record("job-a", "failed")
	_ = d.Flush()
	// second flush should not send again
	_ = d.Flush()
	s.mu.Lock()
	count := len(s.reports)
	s.mu.Unlock()
	if count != 1 {
		t.Errorf("expected exactly 1 report, got %d", count)
	}
}

func TestStartStop_FlushesOnStop(t *testing.T) {
	s := &mockSender{}
	d := New(10*time.Second, s) // long window so ticker doesn't fire
	d.Start()
	d.Record("job-x", "missed")
	if err := d.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}
	if _, ok := s.last(); !ok {
		t.Fatal("expected final flush on Stop")
	}
}
