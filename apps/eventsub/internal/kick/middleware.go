package kick

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/logger"
)

const (
	maxBodySize   = 1 << 20
	maxRequestAge = 5 * time.Minute
	maxFutureSkew = time.Minute
)

type kickContextKey string

const (
	ctxKeyMessageID      kickContextKey = "kick-event-message-id"
	ctxKeyEventType      kickContextKey = "kick-event-type"
	ctxKeyEventVersion   kickContextKey = "kick-event-version"
	ctxKeySubscriptionID kickContextKey = "kick-event-subscription-id"
	ctxKeyTimestamp      kickContextKey = "kick-event-message-timestamp"
)

// KickMessageIDFromContext extracts the Kick-Event-Message-Id from context.
func KickMessageIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyMessageID).(string)
	return v
}

// KickEventTypeFromContext extracts the Kick-Event-Type from context.
func KickEventTypeFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyEventType).(string)
	return v
}

// KickEventVersionFromContext extracts the Kick-Event-Version from context.
func KickEventVersionFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyEventVersion).(string)
	return v
}

// KickSubscriptionIDFromContext extracts the Kick-Event-Subscription-Id from context.
func KickSubscriptionIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeySubscriptionID).(string)
	return v
}

// KickMessageTimestampFromContext extracts the Kick-Event-Message-Timestamp from context.
func KickMessageTimestampFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyTimestamp).(string)
	return v
}

// Middleware verifies the Kick webhook signature and populates context with event headers.
type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(log *slog.Logger) *Middleware {
	return &Middleware{
		logger: log,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageID := r.Header.Get("Kick-Event-Message-Id")
		timestamp := r.Header.Get("Kick-Event-Message-Timestamp")
		eventType := r.Header.Get("Kick-Event-Type")
		eventVersion := r.Header.Get("Kick-Event-Version")
		subscriptionID := r.Header.Get("Kick-Event-Subscription-Id")
		sigHeader := r.Header.Get("Kick-Event-Signature")

		m.logger.InfoContext(r.Context(), "incoming kick webhook request",
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("event_type", eventType),
			slog.String("message_id", messageID),
			slog.Bool("has_signature", sigHeader != ""),
		)

		if messageID == "" || timestamp == "" || eventType == "" || sigHeader == "" {
			m.logger.WarnContext(r.Context(), "missing required kick webhook headers")
			http.Error(w, "missing required headers", http.StatusBadRequest)
			return
		}

		if _, err := validateTimestamp(timestamp, time.Now()); err != nil {
			m.logger.WarnContext(r.Context(), "invalid kick webhook timestamp", logger.Error(err))
			http.Error(w, "invalid timestamp", http.StatusBadRequest)
			return
		}

		body, err := readLimitedRequestBody(w, r)
		if err != nil {
			m.logBodyReadError(r.Context(), err)
			return
		}

		if !gokick.ValidateEvent(r.Header, body) {
			m.logger.WarnContext(r.Context(), "kick signature verification failed")
			http.Error(w, "invalid signature", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyMessageID, messageID)
		ctx = context.WithValue(ctx, ctxKeyEventType, eventType)
		ctx = context.WithValue(ctx, ctxKeyEventVersion, eventVersion)
		ctx = context.WithValue(ctx, ctxKeySubscriptionID, subscriptionID)
		ctx = context.WithValue(ctx, ctxKeyTimestamp, timestamp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandlerWithoutVerification extracts headers and populates context without verifying the signature.
func (m *Middleware) HandlerWithoutVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageID := r.Header.Get("Kick-Event-Message-Id")
		timestamp := r.Header.Get("Kick-Event-Message-Timestamp")
		eventType := r.Header.Get("Kick-Event-Type")
		eventVersion := r.Header.Get("Kick-Event-Version")
		subscriptionID := r.Header.Get("Kick-Event-Subscription-Id")

		m.logger.InfoContext(r.Context(), "incoming kick webhook request (no sig verify)",
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("event_type", eventType),
			slog.String("message_id", messageID),
		)

		if messageID == "" || timestamp == "" || eventType == "" {
			m.logger.WarnContext(r.Context(), "missing required kick webhook headers")
			http.Error(w, "missing required headers", http.StatusBadRequest)
			return
		}

		if _, err := validateTimestamp(timestamp, time.Now()); err != nil {
			m.logger.WarnContext(r.Context(), "invalid kick webhook timestamp", logger.Error(err))
			http.Error(w, "invalid timestamp", http.StatusBadRequest)
			return
		}

		if _, err := readLimitedRequestBody(w, r); err != nil {
			m.logBodyReadError(r.Context(), err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyMessageID, messageID)
		ctx = context.WithValue(ctx, ctxKeyEventType, eventType)
		ctx = context.WithValue(ctx, ctxKeyEventVersion, eventVersion)
		ctx = context.WithValue(ctx, ctxKeySubscriptionID, subscriptionID)
		ctx = context.WithValue(ctx, ctxKeyTimestamp, timestamp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func readLimitedRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return nil, err
		}

		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	return body, nil
}

func (m *Middleware) logBodyReadError(ctx context.Context, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		m.logger.WarnContext(ctx, "kick webhook request body exceeds size limit", logger.Error(err))
		return
	}

	m.logger.ErrorContext(ctx, "failed to read kick webhook request body", logger.Error(err))
}

func validateTimestamp(timestamp string, now time.Time) (time.Time, error) {
	// Kick sometimes sends duplicated dates like: 2026-04-26T2026-04-26T11:27:44Z
	if strings.Count(timestamp, "T") == 2 {
		parts := strings.SplitN(timestamp, "T", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[1], parts[0]+"T") {
			timestamp = parts[0] + "T" + strings.TrimPrefix(parts[1], parts[0]+"T")
		}
	}

	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp: %w", err)
	}

	if parsed.Before(now.Add(-maxRequestAge)) {
		return time.Time{}, fmt.Errorf("timestamp older than %s", maxRequestAge)
	}

	if parsed.After(now.Add(maxFutureSkew)) {
		return time.Time{}, fmt.Errorf("timestamp more than %s in future", maxFutureSkew)
	}

	return parsed, nil
}
