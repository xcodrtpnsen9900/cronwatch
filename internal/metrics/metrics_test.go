package metrics_test

import (
	"testing"

	"github.com/example/cronwatch/internal/metrics"
)

func TestGlobal_Counters(t *testing.T) {
	r := metrics.New()
	g := r.Global()
	g.AlertsTotal.Add(3)
	g.ChecksTotal.Add(10)

	if v := g.AlertsTotal.Load(); v != 3 {
		t.Fatalf("AlertsTotal: want 3, got %d", v)
	}
	if v := g.ChecksTotal.Load(); v != 10 {
		t.Fatalf("ChecksTotal: want 10, got %d", v)
	}
}

func TestJob_CreatedOnFirstAccess(t *testing.T) {
	r := metrics.New()
	c1 := r.Job("backup")
	c2 := r.Job("backup")
	if c1 != c2 {
		t.Fatal("expected same pointer for same job name")
	}
}

func TestJob_IndependentCounters(t *testing.T) {
	r := metrics.New()
	r.Job("jobA").MissedTotal.Add(2)
	r.Job("jobB").FailedTotal.Add(5)

	if v := r.Job("jobA").MissedTotal.Load(); v != 2 {
		t.Fatalf("jobA MissedTotal: want 2, got %d", v)
	}
	if v := r.Job("jobB").MissedTotal.Load(); v != 0 {
		t.Fatalf("jobB MissedTotal: want 0, got %d", v)
	}
}

func TestSnapshot_ContainsGlobalAndJobs(t *testing.T) {
	r := metrics.New()
	r.Global().AlertsTotal.Add(1)
	r.Job("nightly").RecoveredTotal.Add(4)

	snap := r.Snapshot()

	g, ok := snap["__global__"]
	if !ok {
		t.Fatal("snapshot missing __global__ key")
	}
	if g["alerts_total"] != 1 {
		t.Fatalf("global alerts_total: want 1, got %d", g["alerts_total"])
	}

	j, ok := snap["nightly"]
	if !ok {
		t.Fatal("snapshot missing nightly key")
	}
	if j["recovered_total"] != 4 {
		t.Fatalf("nightly recovered_total: want 4, got %d", j["recovered_total"])
	}
}

func TestSnapshot_AllKeysPresent(t *testing.T) {
	r := metrics.New()
	snap := r.Snapshot()
	wantKeys := []string{"alerts_total", "checks_total", "missed_total", "failed_total", "recovered_total"}
	for _, k := range wantKeys {
		if _, ok := snap["__global__"][k]; !ok {
			t.Errorf("global snapshot missing key %q", k)
		}
	}
}
