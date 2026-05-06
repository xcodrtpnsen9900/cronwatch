package history_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/history"
)

func makeEntry(job string, success bool) history.Entry {
	return history.Entry{
		JobName:   job,
		Timestamp: time.Now(),
		Success:   success,
		Message:   "test",
	}
}

func TestRecord_And_Latest(t *testing.T) {
	s := history.New(10)

	e := makeEntry("backup", true)
	s.Record(e)

	got, ok := s.Latest("backup")
	if !ok {
		t.Fatal("expected an entry")
	}
	if got.JobName != "backup" || got.Success != true {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestLatest_Missing(t *testing.T) {
	s := history.New(10)
	_, ok := s.Latest("nonexistent")
	if ok {
		t.Fatal("expected no entry for unknown job")
	}
}

func TestAll_Order(t *testing.T) {
	s := history.New(10)
	for i := 0; i < 3; i++ {
		s.Record(makeEntry("job", i%2 == 0))
	}
	entries := s.All("job")
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestRecord_Prunes_Old_Entries(t *testing.T) {
	max := 5
	s := history.New(max)

	for i := 0; i < 12; i++ {
		s.Record(makeEntry("prune-job", true))
	}

	entries := s.All("prune-job")
	if len(entries) != max {
		t.Errorf("expected %d entries after pruning, got %d", max, len(entries))
	}
}

func TestRecord_IsolatesJobs(t *testing.T) {
	s := history.New(10)
	s.Record(makeEntry("jobA", true))
	s.Record(makeEntry("jobB", false))

	a := s.All("jobA")
	b := s.All("jobB")

	if len(a) != 1 || len(b) != 1 {
		t.Errorf("expected 1 entry each, got %d and %d", len(a), len(b))
	}
	if a[0].Success != true || b[0].Success != false {
		t.Error("job entries mixed up")
	}
}
