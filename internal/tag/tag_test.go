package tag_test

import (
	"sort"
	"testing"

	"github.com/cronwatch/cronwatch/internal/tag"
)

func newStore() *tag.Store { return tag.New() }

func TestSet_ReplacesTags(t *testing.T) {
	s := newStore()
	s.Set("job1", []string{"prod", "critical"})
	s.Set("job1", []string{"staging"})
	if s.Has("job1", "prod") {
		t.Fatal("expected old tag to be replaced")
	}
	if !s.Has("job1", "staging") {
		t.Fatal("expected new tag to be present")
	}
}

func TestAdd_AppendsTags(t *testing.T) {
	s := newStore()
	s.Add("job1", "prod")
	s.Add("job1", "critical")
	if !s.Has("job1", "prod", "critical") {
		t.Fatal("expected both tags to be present")
	}
}

func TestHas_AllRequired(t *testing.T) {
	s := newStore()
	s.Set("job1", []string{"prod", "team-a"})
	if !s.Has("job1", "prod") {
		t.Fatal("single tag should match")
	}
	if s.Has("job1", "prod", "team-b") {
		t.Fatal("missing tag should return false")
	}
}

func TestHas_UnknownJob(t *testing.T) {
	s := newStore()
	if s.Has("ghost", "prod") {
		t.Fatal("unknown job should return false")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := newStore()
	s.Set("job1", []string{"a", "b", "c"})
	tags := s.Get("job1")
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
}

func TestJobsWithTag(t *testing.T) {
	s := newStore()
	s.Add("job1", "prod")
	s.Add("job2", "staging")
	s.Add("job3", "prod")

	jobs := s.JobsWithTag("prod")
	sort.Strings(jobs)
	if len(jobs) != 2 || jobs[0] != "job1" || jobs[1] != "job3" {
		t.Fatalf("unexpected jobs: %v", jobs)
	}
}

func TestRemove_DeletesTag(t *testing.T) {
	s := newStore()
	s.Set("job1", []string{"prod", "critical"})
	s.Remove("job1", "critical")
	if s.Has("job1", "critical") {
		t.Fatal("removed tag should not be present")
	}
	if !s.Has("job1", "prod") {
		t.Fatal("remaining tag should still be present")
	}
}
