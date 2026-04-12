package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	cfg "github.com/twirapp/twir/libs/config"
)

// Server is the HTTP server for the eventsub app.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// New creates a new Server and registers all routes.
func New(config cfg.Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.EventsubHttpPort),
			Handler: mux,
		},
		logger: logger,
	}

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /webhook/kick", s.handleWebhookKick)

	return s
}

// Start begins listening and serving HTTP requests in a goroutine.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("eventsub http server listen: %w", err)
	}

	s.logger.Info("eventsub HTTP server started", slog.String("addr", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			s.logger.Error("eventsub HTTP server error", slog.Any("error", err))
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("eventsub HTTP server stopping")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleWebhookKick(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
