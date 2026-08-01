package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
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
	targetChannelID := uuid.New()
	initiatorUserID := uuid.New()
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := auth.SetOAuthAttempt(ctx, state, OAuthAttempt{
		Platform:        platform.PlatformVKVideoLive,
		RedirectTo:      "/dashboard/settings",
		CodeVerifier:    "pkce-verifier",
		TargetChannelID: &targetChannelID,
		InitiatorUserID: &initiatorUserID,
		ExpiresAt:       expiresAt,
	}); err != nil {
		t.Fatalf("store OAuth attempt: %v", err)
	}

	attempt, err := auth.GetOAuthAttempt(ctx, state)
	if err != nil {
		t.Fatalf("load OAuth attempt: %v", err)
	}
	if attempt.Platform != platform.PlatformVKVideoLive || attempt.RedirectTo != "/dashboard/settings" || attempt.CodeVerifier != "pkce-verifier" || attempt.DeviceID != "" || attempt.TargetChannelID == nil || *attempt.TargetChannelID != targetChannelID || attempt.InitiatorUserID == nil || *attempt.InitiatorUserID != initiatorUserID || !attempt.ExpiresAt.Equal(expiresAt) {
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

func TestIntegrationOAuthAttemptConsumesStateOnlyOnce(t *testing.T) {
	registerSessionTypes()
	sessionManager := scs.New()
	ctx, err := sessionManager.Load(context.Background(), "")
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	auth := &Auth{sessionManager: sessionManager}
	channelID := uuid.New()
	userID := uuid.New()

	state, err := auth.CreateIntegrationOAuthAttempt(
		ctx,
		integrationsmodel.ServiceNightbot,
		channelID,
		userID,
	)
	if err != nil {
		t.Fatalf("create integration OAuth attempt: %v", err)
	}
	if state == "" {
		t.Fatal("created integration OAuth state is empty")
	}

	if err := auth.ConsumeIntegrationOAuthAttempt(
		ctx,
		state,
		integrationsmodel.ServiceNightbot,
		channelID,
		userID,
		time.Now(),
	); err != nil {
		t.Fatalf("consume integration OAuth attempt: %v", err)
	}

	err = auth.ConsumeIntegrationOAuthAttempt(
		ctx,
		state,
		integrationsmodel.ServiceNightbot,
		channelID,
		userID,
		time.Now(),
	)
	if !errors.Is(err, ErrOAuthAttemptNotFound) {
		t.Fatalf("replayed integration OAuth attempt error = %v, want ErrOAuthAttemptNotFound", err)
	}
}

func TestIntegrationOAuthAttemptRejectsMismatchedBindingsAndExpiry(t *testing.T) {
	channelID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name            string
		service         integrationsmodel.Service
		channelID       uuid.UUID
		initiatorUserID uuid.UUID
		expireAttempt   bool
		wantErr         error
	}{
		{
			name:            "wrong service",
			service:         integrationsmodel.ServiceStreamElements,
			channelID:       channelID,
			initiatorUserID: userID,
			wantErr:         ErrOAuthAttemptMismatch,
		},
		{
			name:            "wrong channel",
			service:         integrationsmodel.ServiceNightbot,
			channelID:       uuid.New(),
			initiatorUserID: userID,
			wantErr:         ErrOAuthAttemptMismatch,
		},
		{
			name:            "wrong initiator",
			service:         integrationsmodel.ServiceNightbot,
			channelID:       channelID,
			initiatorUserID: uuid.New(),
			wantErr:         ErrOAuthAttemptMismatch,
		},
		{
			name:            "expired",
			service:         integrationsmodel.ServiceNightbot,
			channelID:       channelID,
			initiatorUserID: userID,
			expireAttempt:   true,
			wantErr:         ErrOAuthAttemptExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerSessionTypes()
			sessionManager := scs.New()
			ctx, err := sessionManager.Load(context.Background(), "")
			if err != nil {
				t.Fatalf("load session: %v", err)
			}
			auth := &Auth{sessionManager: sessionManager}

			state, err := auth.CreateIntegrationOAuthAttempt(
				ctx,
				integrationsmodel.ServiceNightbot,
				channelID,
				userID,
			)
			if err != nil {
				t.Fatalf("create integration OAuth attempt: %v", err)
			}

			now := time.Now()
			if tt.expireAttempt {
				attempt, err := auth.GetOAuthAttempt(ctx, state)
				if err != nil {
					t.Fatalf("load integration OAuth attempt: %v", err)
				}
				now = attempt.ExpiresAt.Add(time.Nanosecond)
			}

			err = auth.ConsumeIntegrationOAuthAttempt(
				ctx,
				state,
				tt.service,
				tt.channelID,
				tt.initiatorUserID,
				now,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("consume integration OAuth attempt error = %v, want %v", err, tt.wantErr)
			}

			if _, err := auth.GetOAuthAttempt(ctx, state); err != nil {
				t.Fatalf("rejected integration OAuth attempt was removed: %v", err)
			}
		})
	}
}
