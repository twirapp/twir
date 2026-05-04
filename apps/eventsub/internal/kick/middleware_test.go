package kick

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestMiddleware(t *testing.T) *Middleware {
	t.Helper()
	return NewMiddleware(slog.Default())
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

func TestValidateTimestamp_DuplicatedDate(t *testing.T) {
	now := time.Date(2026, 4, 26, 11, 30, 0, 0, time.UTC)
	malformed := "2026-04-26T2026-04-26T11:27:44Z"

	parsed, err := validateTimestamp(malformed, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2026, 4, 26, 11, 27, 44, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, parsed)
	}
}
