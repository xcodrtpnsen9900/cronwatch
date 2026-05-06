package dedup

import (
	"testing"
	"time"
)

func newDedup(ttl time.Duration) *Deduplicator {
	d := New(ttl)
	return d
}

func TestIsDuplicate_FirstCallAllowed(t *testing.T) {
	d := newDedup(time.Minute)
	k := Key{Job: "backup", Kind: "missed"}
	if d.IsDuplicate(k) {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallSuppressed(t *testing.T) {
	d := newDedup(time.Minute)
	k := Key{Job: "backup", Kind: "missed"}
	d.IsDuplicate(k)
	if !d.IsDuplicate(k) {
		t.Fatal("expected second call within TTL to be a duplicate")
	}
}

func TestIsDuplicate_AfterTTLAllowed(t *testing.T) {
	now := time.Now()
	d := newDedup(time.Minute)
	d.now = func() time.Time { return now }

	k := Key{Job: "backup", Kind: "failed"}
	d.IsDuplicate(k)

	d.now = func() time.Time { return now.Add(2 * time.Minute) }
	if d.IsDuplicate(k) {
		t.Fatal("expected call after TTL to not be a duplicate")
	}
}

func TestIsDuplicate_IndependentKeys(t *testing.T) {
	d := newDedup(time.Minute)
	k1 := Key{Job: "jobA", Kind: "missed"}
	k2 := Key{Job: "jobB", Kind: "missed"}
	d.IsDuplicate(k1)
	if d.IsDuplicate(k2) {
		t.Fatal("expected different keys to be independent")
	}
}

func TestReset_AllowsImmediateAlert(t *testing.T) {
	d := newDedup(time.Minute)
	k := Key{Job: "deploy", Kind: "failed"}
	d.IsDuplicate(k)
	d.Reset(k)
	if d.IsDuplicate(k) {
		t.Fatal("expected reset key to allow immediate alert")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	d := newDedup(time.Minute)
	d.now = func() time.Time { return now }

	k := Key{Job: "nightly", Kind: "missed"}
	d.IsDuplicate(k)

	d.now = func() time.Time { return now.Add(2 * time.Minute) }
	d.Purge()

	if len(d.records) != 0 {
		t.Fatalf("expected records to be empty after purge, got %d", len(d.records))
	}
}

func TestPurge_KeepsActiveEntries(t *testing.T) {
	now := time.Now()
	d := newDedup(time.Hour)
	d.now = func() time.Time { return now }

	k := Key{Job: "weekly", Kind: "failed"}
	d.IsDuplicate(k)

	d.now = func() time.Time { return now.Add(30 * time.Minute) }
	d.Purge()

	if len(d.records) != 1 {
		t.Fatalf("expected 1 active record after purge, got %d", len(d.records))
	}
}
