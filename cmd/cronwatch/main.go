package main

import (
	"flag"
	"log"
	"os"

	"github.com/cronwatch/internal/config"
)

func main() {
	configPath := flag.String("config", "configs/example.json", "path to JSON config file")
	flag.Parse()

	logger := log.New(os.Stdout, "[cronwatch] ", log.LstdFlags|log.Lmsgprefix)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	logger.Printf("loaded %d job(s), check interval: %s", len(cfg.Jobs), cfg.CheckInterval)
	for _, job := range cfg.Jobs {
		hook := job.WebhookURL
		if hook == "" {
			hook = cfg.DefaultWebhook
		}
		logger.Printf("  job %q | schedule: %s | timeout: %s | webhook: %s",
			job.Name, job.Schedule, job.Timeout, hook)
	}

	logger.Println("cronwatch started — monitoring cron jobs...")
	// Further monitor/alert logic will be wired here in subsequent phases.
	select {}
}
