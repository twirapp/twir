package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	cfg "github.com/twirapp/twir/libs/config"
)

// Server is the HTTP server for the eventsub app.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// New creates a new Server and registers all routes.
func New(config cfg.Config, logger *slog.Logger, redisClient *redis.Client) *Server {
	mux := http.NewServeMux()

	s := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.EventsubHttpPort),
			Handler: mux,
		},
		logger: logger,
	}

	kickMiddleware := kick.NewMiddleware(redisClient, logger)

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.Handle("POST /webhook/kick", kickMiddleware.Handler(http.HandlerFunc(s.handleWebhookKick)))

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

func (s *Server) handleWebhookKick(w http.ResponseWriter, r *http.Request) {
	messageID := kick.KickMessageIDFromContext(r.Context())
	eventType := kick.KickEventTypeFromContext(r.Context())
	s.logger.InfoContext(r.Context(), "kick webhook received",
		slog.String("message_id", messageID),
		slog.String("event_type", eventType),
	)
	w.WriteHeader(http.StatusOK)
}
