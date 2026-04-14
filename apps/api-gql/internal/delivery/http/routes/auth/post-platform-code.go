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
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	userplatformaccountsrepo "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
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

	if a.userPlatformAccountsRepo == nil {
		a.logger.ErrorContext(ctx, "kick auth: user platform accounts repo is nil")
		return nil, huma.Error500InternalServerError("User platform accounts repository is not configured", fmt.Errorf("user platform accounts repo is nil"))
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
	a.logger.InfoContext(ctx, "kick auth: code verifier retrieved")

	tokens, err := a.kickProvider.ExchangeCode(ctx, input.Code, codeVerifier)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to exchange code", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot exchange code", err)
	}
	a.logger.InfoContext(ctx, "kick auth: token exchange successful")

	platformUser, err := a.kickProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get user from kick", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user data from kick", err)
	}
	a.logger.InfoContext(ctx, "kick auth: got user from kick", slog.String("kick_user_id", platformUser.ID), slog.String("kick_login", platformUser.Login))

	encryptedAccess, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to encrypt access token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt user access token", err)
	}

	encryptedRefresh, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to encrypt refresh token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt user refresh token", err)
	}

	account, err := a.userPlatformAccountsRepo.GetByPlatformUserID(ctx, platform.PlatformKick, platformUser.ID)
	accountNotFound := errors.Is(err, userplatformaccountsrepo.ErrNotFound)
	if err != nil && !accountNotFound {
		a.logger.ErrorContext(ctx, "kick auth: failed to get platform account", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user platform account", err)
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

	userID := uuid.Nil
	createdUser := false
	if accountNotFound {
		a.logger.InfoContext(ctx, "kick auth: no existing platform account, creating new user")
		emptyTwitchID := ""
		userID, err = a.createUser(ctx, &emptyTwitchID)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to create user", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create user", err)
		}
		createdUser = true
		a.logger.InfoContext(ctx, "kick auth: created new user", slog.String("user_id", userID.String()))
	} else {
		userID = account.UserID
		a.logger.InfoContext(ctx, "kick auth: found existing platform account", slog.String("user_id", userID.String()))
	}

	_, err = a.userPlatformAccountsRepo.Upsert(ctx, userplatformaccountsrepo.UpsertInput{
		UserID:              userID,
		Platform:            platform.PlatformKick,
		PlatformUserID:      platformUser.ID,
		PlatformLogin:       platformUser.Login,
		PlatformDisplayName: platformUser.DisplayName,
		PlatformAvatar:      platformUser.Avatar,
		AccessToken:         encryptedAccess,
		RefreshToken:        encryptedRefresh,
		Scopes:              tokens.Scopes,
		ExpiresIn:           tokens.ExpiresIn,
		ObtainmentTimestamp: time.Now().UTC(),
	})
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to upsert platform account", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot upsert user platform account", err)
	}
	a.logger.InfoContext(ctx, "kick auth: upserted platform account")

	channel, err := a.channelsRepo.GetByUserIDAndPlatform(ctx, userID, platform.PlatformKick)
	if err != nil {
		if !errors.Is(err, channelsrepo.ErrNotFound) {
			a.logger.ErrorContext(ctx, "kick auth: failed to get channel", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot get channel", err)
		}

		channel, err = a.createChannel(ctx, userID, defaultBot.ID, platform.PlatformKick)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to create channel", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create channel", err)
		}
		a.logger.InfoContext(ctx, "kick auth: created channel", slog.String("channel_id", channel.ID))
	} else {
		a.logger.InfoContext(ctx, "kick auth: found existing channel", slog.String("channel_id", channel.ID))
	}

	if createdUser {
		if err = a.bus.Scheduler.CreateDefaultRoles.Publish(
			ctx,
			scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{channel.ID}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default roles", logger.Error(err), slog.String("channel_id", channel.ID))
		}

		if err = a.bus.Scheduler.CreateDefaultCommands.Publish(
			ctx,
			scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{channel.ID}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default commands", logger.Error(err), slog.String("channel_id", channel.ID))
		}
	}

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		ctx,
		buscoreeventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channel.ID},
	); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish eventsub subscribe", logger.Error(err), slog.String("channel_id", channel.ID))
	}

	if err := a.sessions.SetSessionInternalUserID(ctx, userID); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set internal user id", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set internal user id", err)
	}

	if err := a.sessions.SetSessionCurrentPlatform(ctx, platform.PlatformKick.String()); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set current platform", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set current platform", err)
	}

	if err := a.sessions.SetSessionSelectedDashboard(ctx, channel.ID); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set selected dashboard", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set selected dashboard", err)
	}

	if err := a.sessions.SetSessionKickUser(ctx, authsessions.KickSessionUser{
		ID:     platformUser.ID,
		Login:  platformUser.Login,
		Avatar: platformUser.Avatar,
	}); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set kick user", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set kick user", err)
	}

	a.logger.InfoContext(ctx, "kick auth: completed successfully", slog.String("redirect_to", string(redirectTo)), slog.String("user_id", userID.String()), slog.String("channel_id", channel.ID))

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

func (a *Auth) createUser(ctx context.Context, twitchID *string) (uuid.UUID, error) {
	newID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate user id: %w", err)
	}

	apiKey := uuid.New().String()

	finalTwitchID := twitchID
	if finalTwitchID == nil || *finalTwitchID == "" {
		idStr := newID.String()
		finalTwitchID = &idStr
	}

	user, err := a.usersRepo.Create(ctx, usersrepo.CreateInput{
		ID:         newID.String(),
		TwitchID:   finalTwitchID,
		ApiKey:     &apiKey,
		IsBotAdmin: false,
		IsBanned:   false,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return uuid.MustParse(user.ID), nil
}

func (a *Auth) createChannel(
	ctx context.Context,
	userID uuid.UUID,
	botID string,
	platformVal platform.Platform,
) (channelsmodel.Channel, error) {
	channel, err := a.channelsRepo.Create(ctx, channelsrepo.CreateInput{
		UserID:   userID,
		BotID:    botID,
		Platform: platformVal,
	})
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("create channel: %w", err)
	}

	return channel, nil
}
