package envelope

import (
	"os"
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestWrap_SetsSchemaVersion(t *testing.T) {
	b := New("staging").withClock(fixedNow)
	e := b.Wrap(map[string]string{"job": "backup"})
	if e.SchemaVersion != SchemaVersion {
		t.Fatalf("want schema_version %q, got %q", SchemaVersion, e.SchemaVersion)
	}
}

func TestWrap_SetsEnvironment(t *testing.T) {
	b := New("staging").withClock(fixedNow)
	e := b.Wrap(nil)
	if e.Environment != "staging" {
		t.Fatalf("want environment %q, got %q", "staging", e.Environment)
	}
}

func TestWrap_DefaultsToProduction(t *testing.T) {
	b := New("").withClock(fixedNow)
	e := b.Wrap(nil)
	if e.Environment != "production" {
		t.Fatalf("want environment %q, got %q", "production", e.Environment)
	}
}

func TestWrap_SetsSentAt(t *testing.T) {
	b := New("test").withClock(fixedNow)
	e := b.Wrap(nil)
	if !e.SentAt.Equal(fixedNow()) {
		t.Fatalf("want sent_at %v, got %v", fixedNow(), e.SentAt)
	}
}

func TestWrap_SetsHost(t *testing.T) {
	expected, _ := os.Hostname()
	b := New("test").withClock(fixedNow)
	e := b.Wrap(nil)
	if e.Host != expected {
		t.Fatalf("want host %q, got %q", expected, e.Host)
	}
}

func TestWrap_PayloadPassedThrough(t *testing.T) {
	type inner struct{ Name string }
	payload := inner{Name: "nightly-backup"}
	b := New("prod").withClock(fixedNow)
	e := b.Wrap(payload)
	got, ok := e.Payload.(inner)
	if !ok {
		t.Fatal("payload type mismatch")
	}
	if got.Name != payload.Name {
		t.Fatalf("want payload.Name %q, got %q", payload.Name, got.Name)
	}
}
