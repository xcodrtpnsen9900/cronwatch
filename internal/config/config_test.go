package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "cronwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `{
		"default_webhook": "https://hooks.example.com/alert",
		"check_interval": 60000000000,
		"jobs": [
			{"name": "backup", "schedule": "0 2 * * *", "timeout": 3600000000000}
		]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", cfg.Jobs[0].Name)
	}
	if cfg.CheckInterval != 60*time.Second {
		t.Errorf("unexpected check interval: %v", cfg.CheckInterval)
	}
}

func TestLoad_DefaultCheckInterval(t *testing.T) {
	path := writeTempConfig(t, `{
		"default_webhook": "https://hooks.example.com/alert",
		"jobs": [{"name": "sync", "schedule": "*/5 * * * *"}]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CheckInterval != 30*time.Second {
		t.Errorf("expected default 30s interval, got %v", cfg.CheckInterval)
	}
}

func TestLoad_MissingWebhook(t *testing.T) {
	path := writeTempConfig(t, `{
		"jobs": [{"name": "sync", "schedule": "*/5 * * * *"}]
	}`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook, got nil")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	path := writeTempConfig(t, `{"default_webhook": "https://hooks.example.com/alert", "jobs": []}`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty jobs, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
