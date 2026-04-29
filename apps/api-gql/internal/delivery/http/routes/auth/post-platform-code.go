package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
)

type kickCodeBody struct {
	Code  string `json:"code" minLength:"1" required:"true"`
	State string `json:"state" required:"true"`
}

func (a *Auth) handleKickCode(
	ctx context.Context,
	input kickCodeBody,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	a.logger.InfoContext(ctx, "kick auth: started", slog.String("state", input.State))

	redirectTo, err := decodeRedirectState(input.State)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to decode state", logger.Error(err))
		return nil, huma.Error400BadRequest("Cannot decode state", err)
	}

	if a.kickProvider == nil {
		a.logger.ErrorContext(ctx, "kick auth: kick provider is nil")
		return nil, huma.Error500InternalServerError("Kick provider is not configured", fmt.Errorf("kick provider is nil"))
	}

	if a.channelsRepo == nil {
		a.logger.ErrorContext(ctx, "kick auth: channels repo is nil")
		return nil, huma.Error500InternalServerError("Channels repository is not configured", fmt.Errorf("channels repo is nil"))
	}

	if a.botsRepo == nil {
		a.logger.ErrorContext(ctx, "kick auth: bots repo is nil")
		return nil, huma.Error500InternalServerError("Bots repository is not configured", fmt.Errorf("bots repo is nil"))
	}

	codeVerifier, err := a.getKickCodeVerifier(ctx)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get code verifier", logger.Error(err))
		return nil, huma.Error400BadRequest("Cannot get code verifier", err)
	}

	tokens, err := a.kickProvider.ExchangeCode(ctx, input.Code, codeVerifier)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to exchange code", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot exchange code", err)
	}

	platformUser, err := a.kickProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get user from kick", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user data from kick", err)
	}

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get default bot", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot find default bot", err)
	}
	if defaultBot.ID == "" {
		a.logger.ErrorContext(ctx, "kick auth: no default bot found")
		return nil, huma.Error500InternalServerError("Cannot find default bot", fmt.Errorf("no default bot found"))
	}

	defaultKickBot, kickBotErr := a.kickBotsRepo.GetDefault(ctx)
	if kickBotErr != nil && !errors.Is(kickBotErr, kickbotsrepo.ErrNotFound) {
		a.logger.ErrorContext(ctx, "kick auth: failed to get default kick bot", logger.Error(kickBotErr))
	}

	var defaultKickBotID *uuid.UUID
	if kickBotErr == nil {
		botID, parseErr := uuid.Parse(defaultKickBot.ID)
		if parseErr == nil {
			defaultKickBotID = &botID
		}
	}

	result, err := a.completePlatformAuth(ctx, completePlatformAuthInput{
		Platform: platform.PlatformKick,
		PlatformUser: &appplatform.PlatformUser{
			ID:          platformUser.ID,
			Login:       platformUser.Login,
			DisplayName: platformUser.DisplayName,
			Avatar:      platformUser.Avatar,
		},
		Tokens: &appplatform.PlatformTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    tokens.ExpiresIn,
			Scopes:       tokens.Scopes,
		},
		DefaultBotID:     defaultBot.ID,
		DefaultKickBotID: defaultKickBotID,
	})
	if err != nil {
		if errors.Is(err, errAuthForbidden) {
			return nil, huma.Error403Forbidden("Forbidden", nil)
		}

		if errors.Is(err, errPlatformConflict) {
			return nil, huma.Error409Conflict("Platform account already linked to another dashboard", err)
		}

		a.logger.ErrorContext(ctx, "kick auth: failed to complete auth", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot complete auth", err)
	}

	if !result.CreatedUser {
		accessToken, encryptAccessErr := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
		if encryptAccessErr != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to encrypt access token for kick bot update", logger.Error(encryptAccessErr))
		} else {
			refreshToken, encryptRefreshErr := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
			if encryptRefreshErr != nil {
				a.logger.ErrorContext(ctx, "kick auth: failed to encrypt refresh token for kick bot update", logger.Error(encryptRefreshErr))
			} else {
				existingKickBot, kickBotByUserErr := a.kickBotsRepo.GetByKickUserID(ctx, result.PlatformUserID)
				if kickBotByUserErr != nil && !errors.Is(kickBotByUserErr, kickbotsrepo.ErrNotFound) {
					a.logger.ErrorContext(ctx, "kick auth: failed to get kick bot by user id", logger.Error(kickBotByUserErr))
				}

				if kickBotByUserErr == nil {
					existingBotID, parseErr := uuid.Parse(existingKickBot.ID)
					if parseErr == nil {
						_, updateErr := a.kickBotsRepo.UpdateToken(
							ctx,
							existingBotID,
							kickbotsrepo.UpdateTokenInput{
								AccessToken:         accessToken,
								RefreshToken:        refreshToken,
								Scopes:              tokens.Scopes,
								ExpiresIn:           tokens.ExpiresIn,
								ObtainmentTimestamp: time.Now().UTC(),
							},
						)
						if updateErr != nil {
							a.logger.ErrorContext(ctx, "kick auth: failed to update kick bot token on re-login", logger.Error(updateErr))
						}
					}
				}
			}
		}
	}

	if err := a.sessions.SetSessionKickUser(ctx, authsessions.KickSessionUser{
		ID:     platformUser.ID,
		Login:  platformUser.Login,
		Avatar: platformUser.Avatar,
	}); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set kick user", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set kick user", err)
	}

	a.logger.InfoContext(ctx, "kick auth: completed successfully", slog.String("redirect_to", string(redirectTo)), slog.String("user_id", result.SessionUserID.String()), slog.String("channel_id", result.Channel.ID.String()))

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}

func decodeRedirectState(state string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(state)
	if err == nil {
		return decoded, nil
	}

	decoded, rawErr := base64.RawURLEncoding.DecodeString(state)
	if rawErr == nil {
		return decoded, nil
	}

	decoded, stdErr := base64.StdEncoding.DecodeString(state)
	if stdErr == nil {
		return decoded, nil
	}

	decoded, rawStdErr := base64.RawStdEncoding.DecodeString(state)
	if rawStdErr == nil {
		return decoded, nil
	}

	return nil, fmt.Errorf("decode state: %w", err)
}

func (a *Auth) getKickCodeVerifier(ctx context.Context) (string, error) {
	codeVerifier, ok := a.sessions.Get(ctx, kickCodeVerifierSessionKey).(string)
	if !ok || codeVerifier == "" {
		return "", fmt.Errorf("kick code verifier not found in session")
	}

	return codeVerifier, nil
}
