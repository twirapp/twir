package auth

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

func TestMCPOAuthAttemptGetReturnsMatchingAttemptWhenMultipleAreStored(t *testing.T) {
	// Given
	now := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)
	auth, ctx := newMCPOAuthTestSession(t, now)
	first := newMCPOAuthAttempt(now.Add(10*time.Minute), "first")
	second := newMCPOAuthAttempt(now.Add(10*time.Minute), "second")
	if err := auth.SetMCPOAuthAttempt(ctx, "opaque-attempt-first", first); err != nil {
		t.Fatalf("store first MCP OAuth attempt: %v", err)
	}
	if err := auth.SetMCPOAuthAttempt(ctx, "opaque-attempt-second", second); err != nil {
		t.Fatalf("store second MCP OAuth attempt: %v", err)
	}

	// When
	attempt, err := auth.GetMCPOAuthAttempt(ctx, "opaque-attempt-first")

	// Then
	if err != nil {
		t.Fatalf("get first MCP OAuth attempt: %v", err)
	}
	requireEqualMCPOAuthAttempt(t, attempt, first)
}

func TestMCPOAuthAttemptGetRejectsAndDeletesExpiredAttempt(t *testing.T) {
	// Given
	now := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)
	auth, ctx := newMCPOAuthTestSession(t, now)
	const attemptID = "opaque-expired-attempt"
	if err := auth.SetMCPOAuthAttempt(ctx, attemptID, newMCPOAuthAttempt(now.Add(-time.Second), "expired")); err != nil {
		t.Fatalf("store expired MCP OAuth attempt: %v", err)
	}

	// When
	_, err := auth.GetMCPOAuthAttempt(ctx, attemptID)

	// Then
	if !errors.Is(err, ErrMCPOAuthAttemptExpired) {
		t.Fatalf("expired MCP OAuth attempt error = %v, want ErrMCPOAuthAttemptExpired", err)
	}
	if _, err := auth.GetMCPOAuthAttempt(ctx, attemptID); !errors.Is(err, ErrMCPOAuthAttemptNotFound) {
		t.Fatalf("removed expired MCP OAuth attempt error = %v, want ErrMCPOAuthAttemptNotFound", err)
	}
}

func TestMCPOAuthAttemptDeleteMakesAttemptSingleUse(t *testing.T) {
	// Given
	now := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)
	auth, ctx := newMCPOAuthTestSession(t, now)
	const attemptID = "opaque-one-use-attempt"
	if err := auth.SetMCPOAuthAttempt(ctx, attemptID, newMCPOAuthAttempt(now.Add(10*time.Minute), "one-use")); err != nil {
		t.Fatalf("store MCP OAuth attempt: %v", err)
	}

	// When
	err := auth.DeleteMCPOAuthAttempt(ctx, attemptID)

	// Then
	if err != nil {
		t.Fatalf("delete MCP OAuth attempt: %v", err)
	}
	if _, err := auth.GetMCPOAuthAttempt(ctx, attemptID); !errors.Is(err, ErrMCPOAuthAttemptNotFound) {
		t.Fatalf("deleted MCP OAuth attempt error = %v, want ErrMCPOAuthAttemptNotFound", err)
	}
}

func TestSessionLogoutClearsMCPOAuthAttempts(t *testing.T) {
	// Given
	now := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)
	auth, ctx := newMCPOAuthTestSession(t, now)
	if err := auth.SetMCPOAuthAttempt(ctx, "opaque-attempt-first", newMCPOAuthAttempt(now.Add(10*time.Minute), "first")); err != nil {
		t.Fatalf("store first MCP OAuth attempt: %v", err)
	}
	if err := auth.SetMCPOAuthAttempt(ctx, "opaque-attempt-second", newMCPOAuthAttempt(now.Add(10*time.Minute), "second")); err != nil {
		t.Fatalf("store second MCP OAuth attempt: %v", err)
	}

	// When
	err := auth.SessionLogout(ctx)

	// Then
	if err != nil {
		t.Fatalf("logout session: %v", err)
	}
	if _, err := auth.GetMCPOAuthAttempt(ctx, "opaque-attempt-first"); !errors.Is(err, ErrMCPOAuthAttemptNotFound) {
		t.Fatalf("first MCP OAuth attempt after logout error = %v, want ErrMCPOAuthAttemptNotFound", err)
	}
	if _, err := auth.GetMCPOAuthAttempt(ctx, "opaque-attempt-second"); !errors.Is(err, ErrMCPOAuthAttemptNotFound) {
		t.Fatalf("second MCP OAuth attempt after logout error = %v, want ErrMCPOAuthAttemptNotFound", err)
	}
}

func newMCPOAuthTestSession(t *testing.T, now time.Time) (*Auth, context.Context) {
	t.Helper()
	registerSessionTypes()

	sessionManager := scs.New()
	ctx, err := sessionManager.Load(context.Background(), "")
	if err != nil {
		t.Fatalf("load session: %v", err)
	}

	return &Auth{
		sessionManager: sessionManager,
		now:            func() time.Time { return now },
	}, ctx
}

func newMCPOAuthAttempt(expiresAt time.Time, suffix string) MCPOAuthAttempt {
	return MCPOAuthAttempt{
		ClientID:        "client-" + suffix,
		RedirectURI:     "https://client.example/callback/" + suffix,
		ClientState:     "client-state-" + suffix,
		CodeChallenge:   "pkce-challenge-" + suffix,
		RequestedScopes: []string{"chat:read", "chat:write:" + suffix},
		Resource:        "https://resource.example/" + suffix,
		CSRFToken:       "csrf-token-" + suffix,
		ExpiresAt:       expiresAt,
	}
}

func requireEqualMCPOAuthAttempt(t *testing.T, got, want MCPOAuthAttempt) {
	t.Helper()
	if got.ClientID != want.ClientID || got.RedirectURI != want.RedirectURI || got.ClientState != want.ClientState || got.CodeChallenge != want.CodeChallenge || !slices.Equal(got.RequestedScopes, want.RequestedScopes) || got.Resource != want.Resource || got.CSRFToken != want.CSRFToken || !got.ExpiresAt.Equal(want.ExpiresAt) {
		t.Fatalf("MCP OAuth attempt = %+v, want %+v", got, want)
	}
}
