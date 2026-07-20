package gsi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"golang.org/x/time/rate"
)

const maxBodySize = 1 << 20

type MatchProcessor interface {
	Process(ctx context.Context, channelID uuid.UUID, payload Payload) error
}

type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	repo       dotarepository.Repository
	processor  MatchProcessor

	rateLimit rate.Limit
	rateBurst int

	limitersMu sync.Mutex
	limiters   map[string]*rateLimiterEntry

	sweepCtx    context.Context
	sweepCancel context.CancelFunc
}

type Option func(*Server)

func WithRateLimit(limit rate.Limit, burst int) Option {
	return func(s *Server) {
		s.rateLimit = limit
		s.rateBurst = burst
	}
}

func New(
	config cfg.Config,
	logger *slog.Logger,
	repo dotarepository.Repository,
	processor MatchProcessor,
	opts ...Option,
) *Server {
	s := &Server{
		logger:    logger,
		repo:      repo,
		processor: processor,
		rateLimit: rate.Limit(10),
		rateBurst: 10,
		limiters:  make(map[string]*rateLimiterEntry),
	}

	for _, opt := range opts {
		opt(s)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /gsi/{token}", s.handleGsi)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.DotaHttpPort),
		Handler: mux,
	}

	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("dota gsi http server listen: %w", err)
	}

	s.sweepCtx, s.sweepCancel = context.WithCancel(context.Background())
	go s.sweepLimiters(s.sweepCtx)

	s.logger.Info("dota GSI HTTP server started", slog.String("addr", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("dota GSI HTTP server error", logger.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("dota GSI HTTP server stopping")
	if s.sweepCancel != nil {
		s.sweepCancel()
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *Server) handleGsi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	settings, err := s.repo.GetByGsiToken(ctx, token)
	if err != nil {
		if errors.Is(err, dotarepository.ErrNotFound) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		s.logger.ErrorContext(ctx, "dota gsi: failed to get settings by token", logger.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "malformed body", http.StatusBadRequest)
		return
	}

	// The token in the URL is authoritative. payload.auth.token is only
	// verified when present; payloads without it are still accepted.
	if payload.Auth.Token != "" && payload.Auth.Token != token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !s.allow(token) {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	if err := s.processor.Process(ctx, settings.ChannelID, payload); err != nil {
		s.logger.ErrorContext(ctx, "dota gsi: failed to process payload", logger.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) allow(token string) bool {
	s.limitersMu.Lock()
	defer s.limitersMu.Unlock()

	entry, ok := s.limiters[token]
	if !ok {
		entry = &rateLimiterEntry{
			limiter: rate.NewLimiter(s.rateLimit, s.rateBurst),
		}
		s.limiters[token] = entry
	}
	entry.lastSeen = time.Now()

	return entry.limiter.Allow()
}

func (s *Server) sweepLimiters(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.limitersMu.Lock()
			for token, entry := range s.limiters {
				if time.Since(entry.lastSeen) > 5*time.Minute {
					delete(s.limiters, token)
				}
			}
			s.limitersMu.Unlock()
		}
	}
}
