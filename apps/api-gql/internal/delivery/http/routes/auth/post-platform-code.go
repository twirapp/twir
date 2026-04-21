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
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
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

	accessToken, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to encrypt access token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt user access token", err)
	}

	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to encrypt refresh token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt user refresh token", err)
	}

	existingUser, err := a.usersRepo.GetByPlatformID(ctx, platform.PlatformKick, platformUser.ID)
	userNotFound := errors.Is(err, usersmodel.ErrNotFound)
	if err != nil && !userNotFound {
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
	if userNotFound {
		a.logger.InfoContext(ctx, "kick auth: no existing platform account, creating new user")
		userID, err = a.createUser(ctx, platform.PlatformKick, platformUser.ID, platformUser.Login, platformUser.DisplayName, platformUser.Avatar)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to create user", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create user", err)
		}
		createdUser = true
		a.logger.InfoContext(ctx, "kick auth: created new user", slog.String("user_id", userID.String()))
	} else {
		userID = uuid.MustParse(existingUser.ID)
		_, updateErr := a.usersRepo.Update(ctx, userID.String(), usersrepo.UpdateInput{
			Login:       &platformUser.Login,
			DisplayName: &platformUser.DisplayName,
			Avatar:      &platformUser.Avatar,
		})
		if updateErr != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to update user profile", logger.Error(updateErr))
		}
		a.logger.InfoContext(ctx, "kick auth: found existing platform account", slog.String("user_id", userID.String()))
	}

	dbUser, err := a.usersRepo.GetByID(ctx, userID.String())
	if err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to get user", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user", err)
	}
	if dbUser.IsBanned {
		return nil, huma.Error403Forbidden("Forbidden", nil)
	}

	currentToken, err := a.tokensRepository.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, tokensrepository.ErrNotFound) {
		a.logger.ErrorContext(ctx, "kick auth: failed to get user token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user token", err)
	}

	tokenExpires := tokens.ExpiresIn
	tokenCreatedAt := time.Now().UTC()
	if currentToken != nil {
		_, err = a.tokensRepository.UpdateTokenByID(
			ctx,
			currentToken.ID,
			tokensrepository.UpdateTokenInput{
				AccessToken:         &accessToken,
				RefreshToken:        &refreshToken,
				ExpiresIn:           &tokenExpires,
				ObtainmentTimestamp: &tokenCreatedAt,
				Scopes:              tokens.Scopes,
			},
		)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to update user token", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot update user token", err)
		}
	} else {
		_, err = a.tokensRepository.CreateUserToken(
			ctx,
			tokensrepository.CreateInput{
				UserID:              userID,
				AccessToken:         accessToken,
				RefreshToken:        refreshToken,
				ExpiresIn:           tokens.ExpiresIn,
				ObtainmentTimestamp: tokenCreatedAt,
				Scopes:              tokens.Scopes,
			},
		)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to create user token", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create user token", err)
		}
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

	channel, err := a.channelsRepo.GetByKickUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, channelsrepo.ErrNotFound) {
			a.logger.ErrorContext(ctx, "kick auth: failed to get channel", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot get channel", err)
		}

		channel, err = a.createChannel(ctx, nil, &userID, defaultBot.ID, defaultKickBotID)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick auth: failed to create channel", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create channel", err)
		}
		a.logger.InfoContext(ctx, "kick auth: created channel", slog.String("channel_id", channel.ID.String()), slog.String("kick_bot_id", defaultKickBot.ID))
	} else {
		if defaultKickBotID != nil && channel.KickBotID == nil {
			updatedChannel, updateErr := a.channelsRepo.Update(ctx, channel.ID, channelsrepo.UpdateInput{KickBotID: defaultKickBotID})
			if updateErr != nil {
				a.logger.ErrorContext(ctx, "kick auth: failed to update channel kick bot", logger.Error(updateErr))
			} else {
				channel = updatedChannel
				a.logger.InfoContext(ctx, "kick auth: updated channel with kick bot", slog.String("channel_id", channel.ID.String()), slog.String("kick_bot_id", defaultKickBot.ID))
			}
		}
		a.logger.InfoContext(ctx, "kick auth: found existing channel", slog.String("channel_id", channel.ID.String()))
	}

	channelIDStr := channel.ID.String()

	if createdUser {
		if err = a.bus.Scheduler.CreateDefaultRoles.Publish(
			ctx,
			scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{channelIDStr}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default roles", logger.Error(err), slog.String("channel_id", channelIDStr))
		}

		if err = a.bus.Scheduler.CreateDefaultCommands.Publish(
			ctx,
			scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{channelIDStr}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default commands", logger.Error(err), slog.String("channel_id", channelIDStr))
		}
	}

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		ctx,
		buscoreeventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelIDStr},
	); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish eventsub subscribe", logger.Error(err), slog.String("channel_id", channelIDStr))
	}

	if err := a.sessions.SetSessionInternalUserID(ctx, userID); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set internal user id", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set internal user id", err)
	}

	if err := a.sessions.SetSessionCurrentPlatform(ctx, platform.PlatformKick.String()); err != nil {
		a.logger.ErrorContext(ctx, "kick auth: failed to set current platform", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot set current platform", err)
	}

	if err := a.sessions.SetSessionSelectedDashboard(ctx, channelIDStr); err != nil {
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

	a.logger.InfoContext(ctx, "kick auth: completed successfully", slog.String("redirect_to", string(redirectTo)), slog.String("user_id", userID.String()), slog.String("channel_id", channelIDStr))

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
