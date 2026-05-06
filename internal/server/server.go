// Package server wires up the HTTP server used by cronwatch.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds HTTP server settings.
type Config struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Port:            8080,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
}

// Server wraps net/http.Server with graceful-shutdown support.
type Server struct {
	httpServer *http.Server
	cfg        Config
}

// New creates a Server with the provided mux and config.
func New(mux http.Handler, cfg Config) *Server {
	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// Start begins listening in a goroutine and returns immediately.
func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// non-fatal: caller should watch the process exit
			_ = err
		}
	}()
}

// Shutdown gracefully stops the server within the configured timeout.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
