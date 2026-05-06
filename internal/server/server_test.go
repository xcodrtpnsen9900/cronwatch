package server_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/server"
)

func freePort() int {
	// Use a high ephemeral port unlikely to conflict in CI.
	return 19876
}

func TestDefaultConfig(t *testing.T) {
	cfg := server.DefaultConfig()
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.ShutdownTimeout <= 0 {
		t.Error("shutdown timeout must be positive")
	}
}

func TestStartStop(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "pong")
	})

	port := freePort()
	cfg := server.Config{
		Port:            port,
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		ShutdownTimeout: 3 * time.Second,
	}
	srv := server.New(mux, cfg)
	srv.Start()

	// Give the server a moment to bind.
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", port))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "pong" {
		t.Errorf("expected pong, got %s", body)
	}

	if err := srv.Shutdown(); err != nil {
		t.Errorf("shutdown error: %v", err)
	}
}
