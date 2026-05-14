package labeler_test

import (
	"testing"

	"github.com/yourorg/cronwatch/internal/labeler"
)

func newStore() *labeler.Store { return labeler.New() }

func TestSet_ReplacesLabels(t *testing.T) {
	s := newStore()
	s.Set("job1", map[string]string{"env": "prod", "team": "ops"})
	s.Set("job1", map[string]string{"env": "staging"})

	all := s.All("job1")
	if len(all) != 1 {
		t.Fatalf("expected 1 label, got %d", len(all))
	}
	if all["env"] != "staging" {
		t.Errorf("expected staging, got %s", all["env"])
	}
}

func TestPut_AddsSingleLabel(t *testing.T) {
	s := newStore()
	s.Put("job1", "region", "us-east-1")
	v, ok := s.Get("job1", "region")
	if !ok || v != "us-east-1" {
		t.Errorf("expected us-east-1, got %q (ok=%v)", v, ok)
	}
}

func TestPut_UpdatesExistingLabel(t *testing.T) {
	s := newStore()
	s.Put("job1", "env", "prod")
	s.Put("job1", "env", "dev")
	v, _ := s.Get("job1", "env")
	if v != "dev" {
		t.Errorf("expected dev, got %s", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := newStore()
	_, ok := s.Get("unknown", "key")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestAll_EmptyForUnknownJob(t *testing.T) {
	s := newStore()
	if got := s.All("nope"); len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := newStore()
	s.Put("job1", "env", "prod")
	if err := s.Delete("job1", "env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("job1", "env")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestDelete_UnknownJobReturnsError(t *testing.T) {
	s := newStore()
	if err := s.Delete("ghost", "key"); err == nil {
		t.Error("expected error for unknown job")
	}
}

func TestJobs_ReturnsLabeledJobs(t *testing.T) {
	s := newStore()
	s.Put("alpha", "k", "v")
	s.Put("beta", "k", "v")

	jobs := s.Jobs()
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := newStore()
	s.Put("job1", "env", "prod")
	copy := s.All("job1")
	copy["env"] = "mutated"
	v, _ := s.Get("job1", "env")
	if v != "prod" {
		t.Error("store was mutated through returned map")
	}
}

func TestDelete_LastKeyRemovesJob(t *testing.T) {
	s := newStore()
	s.Put("job1", "env", "prod")
	if err := s.Delete("job1", "env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// After deleting the only label, the job should no longer appear in Jobs().
	jobs := s.Jobs()
	for _, j := range jobs {
		if j == "job1" {
			t.Error("expected job1 to be removed from Jobs() after last label deleted")
		}
	}
}
