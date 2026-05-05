package alert_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/alert"
)

func fixedBuilder(t time.Time) *alert.Builder {
	b := alert.NewBuilder()
	// swap clock via exported option would be ideal; here we just verify fields
	_ = b
	return alert.NewBuilder()
}

func TestMissed(t *testing.T) {
	b := alert.NewBuilder()
	expected := time.Now().Add(-5 * time.Minute)
	p := b.Missed("backup", expected)

	if p.JobName != "backup" {
		t.Errorf("expected job name 'backup', got %q", p.JobName)
	}
	if p.AlertType != alert.TypeMissed {
		t.Errorf("expected type missed, got %q", p.AlertType)
	}
	if !strings.Contains(p.Message, "backup") {
		t.Errorf("message should contain job name, got %q", p.Message)
	}
	if !strings.Contains(p.Details, expected.Format(time.RFC3339)) {
		t.Errorf("details should contain expected time, got %q", p.Details)
	}
	if p.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestFailed(t *testing.T) {
	b := alert.NewBuilder()
	p := b.Failed("cleanup", "exit code 1")

	if p.AlertType != alert.TypeFailed {
		t.Errorf("expected type failed, got %q", p.AlertType)
	}
	if !strings.Contains(p.Details, "exit code 1") {
		t.Errorf("details should contain reason, got %q", p.Details)
	}
	if !strings.Contains(p.Message, "cleanup") {
		t.Errorf("message should contain job name, got %q", p.Message)
	}
}

func TestRecovered(t *testing.T) {
	b := alert.NewBuilder()
	p := b.Recovered("sync")

	if p.AlertType != alert.TypeRecovered {
		t.Errorf("expected type recovered, got %q", p.AlertType)
	}
	if p.Details != "" {
		t.Errorf("recovered payload should have no details, got %q", p.Details)
	}
	if !strings.Contains(p.Message, "sync") {
		t.Errorf("message should contain job name, got %q", p.Message)
	}
}

func TestPayload_TypeConstants(t *testing.T) {
	if alert.TypeMissed == alert.TypeFailed {
		t.Error("TypeMissed and TypeFailed must be distinct")
	}
	if alert.TypeFailed == alert.TypeRecovered {
		t.Error("TypeFailed and TypeRecovered must be distinct")
	}
}
