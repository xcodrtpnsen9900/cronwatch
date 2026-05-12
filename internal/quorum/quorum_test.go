package quorum

import (
	"testing"
	"time"
)

func newStore(window time.Duration) *Store {
	s := New(window)
	s.now = func() time.Time { return time.Unix(1_000_000, 0) }
	return s
}

func TestCheck_NotMet_NoReports(t *testing.T) {
	s := newStore(time.Minute)
	s.Require("backup", 2)
	st := s.Check("backup")
	if st.Met {
		t.Fatal("expected quorum not met")
	}
	if st.Reported != 0 {
		t.Fatalf("expected 0 reported, got %d", st.Reported)
	}
}

func TestCheck_Met_EnoughReports(t *testing.T) {
	s := newStore(time.Minute)
	s.Require("backup", 2)
	s.Report("backup", "worker-1")
	s.Report("backup", "worker-2")
	st := s.Check("backup")
	if !st.Met {
		t.Fatal("expected quorum met")
	}
	if st.Reported != 2 {
		t.Fatalf("expected 2 reported, got %d", st.Reported)
	}
}

func TestCheck_NotMet_PartialReports(t *testing.T) {
	s := newStore(time.Minute)
	s.Require("backup", 3)
	s.Report("backup", "worker-1")
	st := s.Check("backup")
	if st.Met {
		t.Fatal("expected quorum not met with only 1 of 3")
	}
}

func TestCheck_EvictsExpiredReports(t *testing.T) {
	base := time.Unix(1_000_000, 0)
	s := New(time.Minute)
	s.now = func() time.Time { return base }
	s.Require("sync", 1)
	s.Report("sync", "worker-1")

	// advance time beyond window
	s.now = func() time.Time { return base.Add(2 * time.Minute) }
	st := s.Check("sync")
	if st.Met {
		t.Fatal("expected quorum not met after eviction")
	}
	if st.Reported != 0 {
		t.Fatalf("expected 0 after eviction, got %d", st.Reported)
	}
}

func TestCheck_ZeroRequired_NeverMet(t *testing.T) {
	s := newStore(time.Minute)
	// no Require call → required == 0
	s.Report("job", "w1")
	st := s.Check("job")
	if st.Met {
		t.Fatal("quorum with required=0 should never be met")
	}
}

func TestCheck_InstancesListed(t *testing.T) {
	s := newStore(time.Minute)
	s.Require("etl", 2)
	s.Report("etl", "alpha")
	s.Report("etl", "beta")
	st := s.Check("etl")
	if len(st.Instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(st.Instances))
	}
}
