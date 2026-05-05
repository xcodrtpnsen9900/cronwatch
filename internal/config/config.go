package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// JobConfig defines a single cron job to monitor.
type JobConfig struct {
	Name         string        `json:"name"`
	Schedule     string        `json:"schedule"`
	Timeout      time.Duration `json:"timeout"`
	WebhookURL   string        `json:"webhook_url"`
}

// Config is the top-level application configuration.
type Config struct {
	Jobs           []JobConfig   `json:"jobs"`
	DefaultWebhook string        `json:"default_webhook"`
	CheckInterval  time.Duration `json:"check_interval"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("at least one job must be configured")
	}
	for i, job := range c.Jobs {
		if job.Name == "" {
			return fmt.Errorf("job[%d]: name is required", i)
		}
		if job.Schedule == "" {
			return fmt.Errorf("job[%d] %q: schedule is required", i, job.Name)
		}
		weebhook := job.WebhookURL
		if weebhook == "" {
			weebhook = c.DefaultWebhook
		}
		if weebhook == "" {
			return fmt.Errorf("job[%d] %q: webhook_url or default_webhook is required", i, job.Name)
		}
	}
	if c.CheckInterval == 0 {
		c.CheckInterval = 30 * time.Second
	}
	return nil
}
