package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestOAuthAttemptKeepsPKCEAndCallbackDeviceIDInSession(t *testing.T) {
	registerSessionTypes()
	sessionManager := scs.New()
	ctx, err := sessionManager.Load(context.Background(), "")
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	auth := &Auth{sessionManager: sessionManager}

	const state = "opaque-state"
	if err := auth.SetOAuthAttempt(ctx, state, OAuthAttempt{
		Platform:     platform.PlatformVKVideoLive,
		RedirectTo:   "/dashboard/settings",
		CodeVerifier: "pkce-verifier",
	}); err != nil {
		t.Fatalf("store OAuth attempt: %v", err)
	}

	attempt, err := auth.GetOAuthAttempt(ctx, state)
	if err != nil {
		t.Fatalf("load OAuth attempt: %v", err)
	}
	if attempt.Platform != platform.PlatformVKVideoLive || attempt.RedirectTo != "/dashboard/settings" || attempt.CodeVerifier != "pkce-verifier" || attempt.DeviceID != "" {
		t.Fatalf("stored OAuth attempt = %+v", attempt)
	}

	attempt.DeviceID = "vk-device-id"
	if err := auth.SetOAuthAttempt(ctx, state, attempt); err != nil {
		t.Fatalf("store callback device ID: %v", err)
	}
	storedAttempt, err := auth.GetOAuthAttempt(ctx, state)
	if err != nil {
		t.Fatalf("reload OAuth attempt: %v", err)
	}
	if storedAttempt.DeviceID != "vk-device-id" {
		t.Fatalf("stored device ID = %q, want vk-device-id", storedAttempt.DeviceID)
	}

	if err := auth.DeleteOAuthAttempt(ctx, state); err != nil {
		t.Fatalf("delete OAuth attempt: %v", err)
	}
	if _, err := auth.GetOAuthAttempt(ctx, state); !errors.Is(err, ErrOAuthAttemptNotFound) {
		t.Fatalf("deleted OAuth attempt error = %v, want ErrOAuthAttemptNotFound", err)
	}
}
