package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
)

var errOAuthAttemptPlatformMismatch = errors.New("oauth attempt belongs to another platform")

type platformCodeBody struct {
	Code     string `json:"code" minLength:"1" required:"true"`
	State    string `json:"state" required:"true"`
	DeviceID string `json:"device_id"`
}

type kickCodeBody = platformCodeBody

type platformCodeInput struct {
	Platform platformentity.Platform
	Code     string
	State    string
	DeviceID string
}

type platformCodeResult struct {
	AuthResult   completePlatformAuthResult
	RedirectTo   string
	PlatformUser *appplatform.PlatformUser
	Tokens       *appplatform.PlatformTokens
}

func (a *Auth) StartTwitchAuth(ctx context.Context, redirectTo string) (string, error) {
	return a.startPlatformAuth(ctx, platformentity.PlatformTwitch, redirectTo)
}

func (a *Auth) StartPlatformAuth(
	ctx context.Context,
	platform platformentity.Platform,
	redirectTo string,
) (string, error) {
	return a.startPlatformAuth(ctx, platform, redirectTo)
}

func (a *Auth) StartPlatformAuthForChannel(
	ctx context.Context,
	channelID uuid.UUID,
	platform platformentity.Platform,
	redirectTo string,
) (string, error) {
	sessionUser, hasLiveSession, err := a.getLiveSessionUser(ctx)
	if err != nil {
		return "", fmt.Errorf("get live session user: %w", err)
	}
	if !hasLiveSession || sessionUser.IsBanned {
		return "", errAuthForbidden
	}
	if err := a.authorizeTargetDashboard(ctx, sessionUser, channelID); err != nil {
		return "", err
	}

	return a.startPlatformAuthForChannel(ctx, platform, redirectTo, &channelID)
}

func (a *Auth) startPlatformAuth(
	ctx context.Context,
	platform platformentity.Platform,
	redirectTo string,
) (string, error) {
	return a.startPlatformAuthForChannel(ctx, platform, redirectTo, nil)
}

func (a *Auth) startPlatformAuthForChannel(
	ctx context.Context,
	platform platformentity.Platform,
	redirectTo string,
	targetChannelID *uuid.UUID,
) (string, error) {
	provider, err := a.platformProvider(platform)
	if err != nil {
		return "", err
	}

	codeVerifier, codeChallenge, err := generatePKCE()
	if err != nil {
		return "", err
	}
	if redirectTo == "" {
		redirectTo = "/dashboard"
	}

	state := uuid.NewString()
	if err := a.sessions.SetOAuthAttempt(ctx, state, authsessions.OAuthAttempt{
		Platform:        platform,
		RedirectTo:      redirectTo,
		CodeVerifier:    codeVerifier,
		TargetChannelID: targetChannelID,
	}); err != nil {
		return "", fmt.Errorf("store OAuth attempt: %w", err)
	}

	authorizeURL := provider.GetAuthURL(state, codeChallenge)
	if authorizeURL == "" {
		_ = a.sessions.DeleteOAuthAttempt(ctx, state)
		return "", fmt.Errorf("build %s authorization URL", platform)
	}

	return authorizeURL, nil
}

func (a *Auth) completePlatformCode(ctx context.Context, input platformCodeInput) (platformCodeResult, error) {
	provider, err := a.platformProvider(input.Platform)
	if err != nil {
		return platformCodeResult{}, err
	}

	attempt, err := a.sessions.GetOAuthAttempt(ctx, input.State)
	if err != nil {
		return platformCodeResult{}, fmt.Errorf("get OAuth attempt: %w", err)
	}
	if attempt.Platform != input.Platform {
		return platformCodeResult{}, errOAuthAttemptPlatformMismatch
	}
	if input.DeviceID != "" {
		attempt.DeviceID = input.DeviceID
		if err := a.sessions.SetOAuthAttempt(ctx, input.State, attempt); err != nil {
			return platformCodeResult{}, fmt.Errorf("store callback device ID: %w", err)
		}
	}

	result, err := a.completePlatformExchange(
		ctx,
		input.Platform,
		provider,
		input.Code,
		attempt.CodeVerifier,
		attempt.DeviceID,
		attempt.RedirectTo,
		attempt.TargetChannelID,
	)
	if err != nil {
		return platformCodeResult{}, err
	}
	if err := a.sessions.DeleteOAuthAttempt(ctx, input.State); err != nil {
		return platformCodeResult{}, fmt.Errorf("delete OAuth attempt: %w", err)
	}

	return result, nil
}

func (a *Auth) completePlatformExchange(
	ctx context.Context,
	platform platformentity.Platform,
	provider appplatform.PlatformProvider,
	code string,
	codeVerifier string,
	deviceID string,
	redirectTo string,
	targetChannelID *uuid.UUID,
) (platformCodeResult, error) {
	tokens, err := provider.ExchangeCode(ctx, appplatform.ExchangeCodeInput{
		Code:         code,
		CodeVerifier: codeVerifier,
		DeviceID:     deviceID,
	})
	if err != nil {
		return platformCodeResult{}, fmt.Errorf("exchange platform code: %w", err)
	}

	platformUser, err := provider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		return platformCodeResult{}, fmt.Errorf("get platform user: %w", err)
	}

	bindingConfig, err := a.platformBindingConfig(ctx, platform)
	if err != nil {
		return platformCodeResult{}, fmt.Errorf("get platform binding configuration: %w", err)
	}

	authResult, err := a.completePlatformAuth(ctx, completePlatformAuthInput{
		Platform:        platform,
		PlatformUser:    platformUser,
		Tokens:          tokens,
		BindingConfig:   bindingConfig,
		TargetChannelID: targetChannelID,
	})
	if err != nil {
		return platformCodeResult{}, err
	}

	if hook, ok := a.postPlatformAuthHooks[platform]; ok {
		if err := hook(ctx, authResult, platformUser, tokens); err != nil {
			return platformCodeResult{}, err
		}
	}

	return platformCodeResult{
		AuthResult:   authResult,
		RedirectTo:   redirectTo,
		PlatformUser: platformUser,
		Tokens:       tokens,
	}, nil
}

func (a *Auth) handlePlatformCode(
	ctx context.Context,
	input platformCodeInput,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	result, err := a.completePlatformCode(ctx, input)
	if err != nil {
		return nil, a.platformAuthHTTPError(err)
	}

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: result.RedirectTo}), nil
}

func (a *Auth) handleKickCode(
	ctx context.Context,
	input kickCodeBody,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	return a.handlePlatformCode(ctx, platformCodeInput{
		Platform: platformentity.PlatformKick,
		Code:     input.Code,
		State:    input.State,
		DeviceID: input.DeviceID,
	})
}

func (a *Auth) platformProvider(platform platformentity.Platform) (appplatform.PlatformProvider, error) {
	if a.platformRegistry == nil {
		return nil, errPlatformUnavailable
	}

	provider, ok := a.platformRegistry.Get(platform)
	if !ok || provider == nil {
		return nil, errPlatformUnavailable
	}

	return provider, nil
}

func (a *Auth) platformBindingConfig(
	ctx context.Context,
	platform platformentity.Platform,
) (platformBindingConfig, error) {
	resolver, ok := a.bindingConfigResolvers[platform]
	if !ok {
		return platformBindingConfig{BotConfig: json.RawMessage(`{}`)}, nil
	}

	return resolver(ctx)
}

func (a *Auth) twitchBindingConfig(ctx context.Context) (platformBindingConfig, error) {
	if a.botsRepo == nil {
		return platformBindingConfig{}, fmt.Errorf("bots repository is not configured")
	}

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		return platformBindingConfig{}, fmt.Errorf("get default bot: %w", err)
	}
	if defaultBot.ID == "" {
		return platformBindingConfig{}, fmt.Errorf("default bot not found")
	}

	config, err := json.Marshal(struct {
		BotID          string `json:"bot_id"`
		IsBotMod       bool   `json:"is_bot_mod"`
		IsTwitchBanned bool   `json:"is_twitch_banned"`
	}{BotID: defaultBot.ID})
	if err != nil {
		return platformBindingConfig{}, fmt.Errorf("encode Twitch binding configuration: %w", err)
	}

	return platformBindingConfig{BotConfig: config}, nil
}

func (a *Auth) kickBindingConfig(ctx context.Context) (platformBindingConfig, error) {
	if a.kickBotsRepo == nil {
		return platformBindingConfig{}, fmt.Errorf("kick bots repository is not configured")
	}

	defaultBot, err := a.kickBotsRepo.GetDefault(ctx)
	if errors.Is(err, kickbotsrepo.ErrNotFound) {
		return platformBindingConfig{BotConfig: json.RawMessage(`{}`)}, nil
	}
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get default kick bot", logger.Error(err))
		return platformBindingConfig{BotConfig: json.RawMessage(`{}`)}, nil
	}

	config, err := json.Marshal(struct {
		KickBotID string `json:"kick_bot_id"`
	}{KickBotID: defaultBot.ID.String()})
	if err != nil {
		return platformBindingConfig{}, fmt.Errorf("encode Kick binding configuration: %w", err)
	}

	botUserID := defaultBot.KickUserID
	return platformBindingConfig{BotUserID: &botUserID, BotConfig: config}, nil
}

func (a *Auth) updateKickBotTokenAfterAuth(
	ctx context.Context,
	result completePlatformAuthResult,
	platformUser *appplatform.PlatformUser,
	tokens *appplatform.PlatformTokens,
) error {
	if !result.CreatedUser && a.kickBotsRepo != nil {
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
					_, updateErr := a.kickBotsRepo.UpdateToken(ctx, existingKickBot.ID, kickbotsrepo.UpdateTokenInput{
						AccessToken:         accessToken,
						RefreshToken:        refreshToken,
						Scopes:              tokens.Scopes,
						ExpiresIn:           tokens.ExpiresIn,
						ObtainmentTimestamp: time.Now().UTC(),
					})
					if updateErr != nil {
						a.logger.ErrorContext(ctx, "kick auth: failed to update kick bot token on re-login", logger.Error(updateErr))
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
		return fmt.Errorf("set kick session user: %w", err)
	}

	return nil
}

func (a *Auth) platformAuthHTTPError(err error) error {
	switch {
	case errors.Is(err, errAuthForbidden):
		return huma.Error403Forbidden("Forbidden", nil)
	case errors.Is(err, errPlatformConflict):
		return huma.Error409Conflict("Platform account already linked to another dashboard", err)
	case errors.Is(err, errPlatformUnavailable):
		return huma.Error404NotFound("Platform is not available", nil)
	case errors.Is(err, authsessions.ErrOAuthAttemptNotFound), errors.Is(err, errOAuthAttemptPlatformMismatch):
		return huma.Error400BadRequest("Invalid or expired OAuth state", nil)
	default:
		return huma.Error500InternalServerError("Cannot complete platform auth", err)
	}
}

func generatePKCE() (string, string, error) {
	randomBytes := make([]byte, 96)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("generate PKCE verifier: %w", err)
	}

	verifier := base64.RawURLEncoding.EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(verifier))
	return verifier, base64.RawURLEncoding.EncodeToString(hash[:]), nil
}
