package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Config holds HTTP server configuration.
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:         "",
		Port:         8080,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Server wraps an http.Server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

// New creates a Server bound to the given mux using cfg.
func New(cfg Config, mux http.Handler) (*Server, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("server listen %s: %w", addr, err)
	}
	return &Server{
		httpServer: &http.Server{
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		listener: ln,
	}, nil
}

// Addr returns the address the server is listening on.
func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

// Start begins serving requests in a background goroutine.
func (s *Server) Start() {
	go func() { _ = s.httpServer.Serve(s.listener) }()
}

// Stop gracefully shuts down the server, waiting up to timeout.
func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
