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

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(config cfg.Config, logger *slog.Logger, redisClient *redis.Client, kickHandlers *kick.Handlers) *Server {
	mux := http.NewServeMux()

	s := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.EventsubHttpPort),
			Handler: mux,
		},
		logger: logger,
	}

	mux.HandleFunc("GET /health", s.handleHealth)

	if config.EventSubDisableSignatureVerification {
		logger.Warn("kick webhook signature verification is disabled")
		kickMiddleware := kick.NewMiddleware(redisClient, logger)
		mux.Handle("POST /webhook/kick", kickMiddleware.HandlerWithoutVerification(http.HandlerFunc(kickHandlers.HandleWebhook)))
	} else {
		kickMiddleware := kick.NewMiddleware(redisClient, logger)
		mux.Handle("POST /webhook/kick", kickMiddleware.Handler(http.HandlerFunc(kickHandlers.HandleWebhook)))
	}

	return s
}

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

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("eventsub HTTP server stopping")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
