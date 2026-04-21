package kick

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func newTestMiddleware(t *testing.T) *Middleware {
	t.Helper()

	redisClient, _ := redismock.NewClientMock()

	return NewMiddleware(redisClient, slog.Default())
}

func newTestMiddlewareWithRedisMock(t *testing.T) (*Middleware, redismock.ClientMock) {
	t.Helper()

	redisClient, redisMock := redismock.NewClientMock()

	return NewMiddleware(redisClient, slog.Default()), redisMock
}

func generateTestKeyPair(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyDER})

	return privateKey, string(publicKeyPEM)
}

func signKickRequest(t *testing.T, privateKey *rsa.PrivateKey, messageID, timestamp string, body []byte) string {
	t.Helper()

	hash := sha256.Sum256([]byte(messageID + "." + timestamp + "." + string(body)))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("sign request: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature)
}

func TestHandler_ValidSignature(t *testing.T) {
	middleware, redisMock := newTestMiddlewareWithRedisMock(t)
	privateKey, publicKeyPEM := generateTestKeyPair(t)
	body := []byte(`{"message":"hello"}`)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	signature := signKickRequest(t, privateKey, "message-valid", timestamp, body)

	redisMock.ExpectGet(redisPublicKeyKey).SetVal(publicKeyPEM)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Kick-Event-Message-Id", "message-valid")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "channel.followed")
	req.Header.Set("Kick-Event-Version", "1")
	req.Header.Set("Kick-Event-Subscription-Id", "sub-123")
	req.Header.Set("Kick-Event-Signature", signature)

	w := httptest.NewRecorder()

	middleware.Handler(next).ServeHTTP(w, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations not met: %v", err)
	}
}

func TestHandler_InvalidSignature(t *testing.T) {
	middleware, redisMock := newTestMiddlewareWithRedisMock(t)
	_, publicKeyPEM := generateTestKeyPair(t)
	body := []byte(`{"message":"hello"}`)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	redisMock.ExpectGet(redisPublicKeyKey).SetVal(publicKeyPEM)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Kick-Event-Message-Id", "message-invalid")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "channel.followed")
	req.Header.Set("Kick-Event-Signature", base64.StdEncoding.EncodeToString([]byte("wrong-signature")))

	w := httptest.NewRecorder()

	middleware.Handler(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status code %d, got %d", http.StatusForbidden, w.Code)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations not met: %v", err)
	}
}

func TestHandler_MissingSignature(t *testing.T) {
	middleware := newTestMiddleware(t)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(`{"message":"hello"}`)))
	req.Header.Set("Kick-Event-Message-Id", "message-missing-signature")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "channel.followed")

	w := httptest.NewRecorder()

	middleware.Handler(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerWithoutVerification_Success(t *testing.T) {
	middleware := newTestMiddleware(t)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		if got := r.Context().Value(ctxKeyMessageID); got != "message-123" {
			t.Fatalf("expected message ID in context %q, got %v", "message-123", got)
		}
		if got := r.Context().Value(ctxKeyEventType); got != "channel.followed" {
			t.Fatalf("expected event type in context %q, got %v", "channel.followed", got)
		}
		if got := r.Context().Value(ctxKeyEventVersion); got != "1" {
			t.Fatalf("expected event version in context %q, got %v", "1", got)
		}
		if got := r.Context().Value(ctxKeySubscriptionID); got != "sub-123" {
			t.Fatalf("expected subscription ID in context %q, got %v", "sub-123", got)
		}
		if got := r.Context().Value(ctxKeyTimestamp); got != timestamp {
			t.Fatalf("expected timestamp in context %q, got %v", timestamp, got)
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "channel.followed")
	req.Header.Set("Kick-Event-Version", "1")
	req.Header.Set("Kick-Event-Subscription-Id", "sub-123")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlerWithoutVerification_MissingHeaders(t *testing.T) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			name: "missing message id",
			headers: map[string]string{
				"Kick-Event-Message-Timestamp": timestamp,
			},
		},
		{
			name: "missing timestamp",
			headers: map[string]string{
				"Kick-Event-Message-Id": "message-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := newTestMiddleware(t)

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()

			middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

			if nextCalled {
				t.Fatal("expected next handler not to be called")
			}

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestHandlerWithoutVerification_MissingEventType(t *testing.T) {
	middleware := newTestMiddleware(t)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	// Missing Kick-Event-Type

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerWithoutVerification_StaleTimestamp(t *testing.T) {
	middleware := newTestMiddleware(t)
	staleTimestamp := time.Now().UTC().Add(-10 * time.Minute).Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", staleTimestamp)
	req.Header.Set("Kick-Event-Type", "chat.message.sent")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerWithoutVerification_MalformedTimestamp(t *testing.T) {
	middleware := newTestMiddleware(t)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", "not-a-timestamp")
	req.Header.Set("Kick-Event-Type", "chat.message.sent")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerWithoutVerification_FutureTimestamp(t *testing.T) {
	middleware := newTestMiddleware(t)
	futureTimestamp := time.Now().UTC().Add(5 * time.Minute).Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", futureTimestamp)
	req.Header.Set("Kick-Event-Type", "chat.message.sent")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerWithoutVerification_BodyTooLarge(t *testing.T) {
	middleware := newTestMiddleware(t)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Create a body larger than 1MB
	largeBody := make([]byte, maxBodySize+1)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(largeBody))
	req.Header.Set("Kick-Event-Message-Id", "message-123")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "chat.message.sent")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status code %d, got %d", http.StatusRequestEntityTooLarge, w.Code)
	}
}

func TestHandlerWithoutVerification_ContextValues(t *testing.T) {
	middleware := newTestMiddleware(t)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := KickMessageIDFromContext(r.Context()); got != "message-ctx" {
			t.Fatalf("expected KickMessageIDFromContext %q, got %q", "message-ctx", got)
		}
		if got := KickEventTypeFromContext(r.Context()); got != "livestream.status.updated" {
			t.Fatalf("expected KickEventTypeFromContext %q, got %q", "livestream.status.updated", got)
		}
		if got := KickEventVersionFromContext(r.Context()); got != "2" {
			t.Fatalf("expected KickEventVersionFromContext %q, got %q", "2", got)
		}
		if got := KickSubscriptionIDFromContext(r.Context()); got != "sub-ctx" {
			t.Fatalf("expected KickSubscriptionIDFromContext %q, got %q", "sub-ctx", got)
		}
		if got := KickMessageTimestampFromContext(r.Context()); got != timestamp {
			t.Fatalf("expected KickMessageTimestampFromContext %q, got %q", timestamp, got)
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("Kick-Event-Message-Id", "message-ctx")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "livestream.status.updated")
	req.Header.Set("Kick-Event-Version", "2")
	req.Header.Set("Kick-Event-Subscription-Id", "sub-ctx")

	w := httptest.NewRecorder()

	middleware.HandlerWithoutVerification(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
