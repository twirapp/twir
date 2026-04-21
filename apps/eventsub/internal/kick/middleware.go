package kick

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/logger"
)

const (
	redisPublicKeyKey = "kick:public_key"
	publicKeyTTL      = time.Hour
	maxBodySize       = 1 << 20
	maxRequestAge     = 5 * time.Minute
	maxFutureSkew     = time.Minute

	kickPublicKeyURL = "https://api.kick.com/public/v1/public-key"

	fallbackPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq/+l1WnlRrGSolDMA+A8
6rAhMbQGmQ2SapVcGM3zq8ANXjnhDWocMqfWcTd95btDydITa10kDvHzw9WQOqp2
MZI7ZyrfzJuz5nhTPCiJwTwnEtWft7nV14BYRDHvlfqPUaZ+1KR4OCaO/wWIk/rQ
L/TjY0M70gse8rlBkbo2a8rKhu69RQTRsoaf4DVhDPEeSeI5jVrRDGAMGL3cGuyY
6CLKGdjVEM78g3JfYOvDU/RvfqD7L89TZ3iN94jrmWdGz34JNlEI5hqK8dd7C5EF
BEbZ5jgB8s8ReQV8H+MkuffjdAj3ajDDX3DOJMIut1lBrUVD1AaSrGCKHooWoL2e
twIDAQAB
-----END PUBLIC KEY-----`
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

// Middleware verifies the Kick RSA-SHA256 webhook signature and populates context with event headers.
type Middleware struct {
	redis  *redis.Client
	logger *slog.Logger
}

func NewMiddleware(redisClient *redis.Client, log *slog.Logger) *Middleware {
	return &Middleware{
		redis:  redisClient,
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

		signedString := messageID + "." + timestamp + "." + string(body)

		pubKey, err := m.getPublicKey(r.Context())
		if err != nil {
			m.logger.ErrorContext(r.Context(), "failed to get kick public key", logger.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		sigBytes, err := decodeBase64Signature(sigHeader)
		if err != nil {
			m.logger.WarnContext(r.Context(), "failed to decode kick signature", logger.Error(err))
			http.Error(w, "invalid signature encoding", http.StatusForbidden)
			return
		}

		hash := sha256.Sum256([]byte(signedString))
		if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes); err != nil {
			m.logger.WarnContext(r.Context(), "kick signature verification failed", logger.Error(err))
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

func decodeBase64Signature(sig string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(sig)
	if err == nil {
		return b, nil
	}
	return base64.RawStdEncoding.DecodeString(sig)
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

type kickPublicKeyResponse struct {
	Data struct {
		PublicKey string `json:"public_key"`
	} `json:"data"`
}

func (m *Middleware) getPublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	cached, err := m.redis.Get(ctx, redisPublicKeyKey).Result()
	if err == nil && cached != "" {
		return parseRSAPublicKey(cached)
	}

	keyPEM, err := m.fetchPublicKeyFromAPI(ctx)
	if err != nil {
		m.logger.WarnContext(ctx, "falling back to hardcoded kick public key", logger.Error(err))
		return parseRSAPublicKey(fallbackPublicKey)
	}

	if setErr := m.redis.Set(ctx, redisPublicKeyKey, keyPEM, publicKeyTTL).Err(); setErr != nil {
		m.logger.WarnContext(ctx, "failed to cache kick public key in redis", logger.Error(setErr))
	}

	return parseRSAPublicKey(keyPEM)
}

func (m *Middleware) fetchPublicKeyFromAPI(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, kickPublicKeyURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch kick public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("kick public key API returned status %d", resp.StatusCode)
	}

	var result kickPublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode kick public key response: %w", err)
	}

	if result.Data.PublicKey == "" {
		return "", fmt.Errorf("empty public key in kick API response")
	}

	return result.Data.PublicKey, nil
}

func parseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}

	return rsaKey, nil
}
