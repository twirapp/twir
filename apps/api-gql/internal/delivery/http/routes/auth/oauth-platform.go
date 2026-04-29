package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

var errPlatformConflict = errors.New("platform account is already linked to another channel")
var errAuthForbidden = errors.New("forbidden")

type completePlatformAuthInput struct {
	Platform         platformentity.Platform
	PlatformUser     *appplatform.PlatformUser
	Tokens           *appplatform.PlatformTokens
	DefaultBotID     string
	DefaultKickBotID *uuid.UUID
}

type completePlatformAuthResult struct {
	SessionUserID   uuid.UUID
	PlatformUserID  uuid.UUID
	Channel         channelsmodel.Channel
	CreatedUser     bool
	CreatedChannel  bool
	UsedLiveSession bool
}

func (a *Auth) completePlatformAuth(
	ctx context.Context,
	input completePlatformAuthInput,
) (completePlatformAuthResult, error) {
	if input.PlatformUser == nil {
		return completePlatformAuthResult{}, fmt.Errorf("platform user is required")
	}

	if input.Tokens == nil {
		return completePlatformAuthResult{}, fmt.Errorf("platform tokens are required")
	}

	sessionUser, hasLiveSession, err := a.getLiveSessionUser(ctx)
	if err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("get live session user: %w", err)
	}

	platformUser, createdUser, err := a.getOrCreatePlatformUser(ctx, input.Platform, input.PlatformUser)
	if err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("get or create platform user: %w", err)
	}

	if platformUser.IsBanned {
		return completePlatformAuthResult{}, errAuthForbidden
	}

	platformUserID, err := uuid.Parse(platformUser.ID)
	if err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("parse platform user id: %w", err)
	}

	var (
		channel        channelsmodel.Channel
		createdChannel bool
		sessionUserID  uuid.UUID
	)

	if hasLiveSession {
		if sessionUser.IsBanned {
			return completePlatformAuthResult{}, errAuthForbidden
		}

		sessionUserID, err = uuid.Parse(sessionUser.ID)
		if err != nil {
			return completePlatformAuthResult{}, fmt.Errorf("parse session user id: %w", err)
		}

		channel, createdChannel, err = a.getOrCreateChannelForUser(
			ctx,
			sessionUserID,
			sessionUser.Platform,
			input.DefaultBotID,
			input.DefaultKickBotID,
		)
		if err != nil {
			return completePlatformAuthResult{}, fmt.Errorf("get or create session channel: %w", err)
		}

		if sessionUser.Platform == input.Platform {
			if sessionUser.PlatformID != input.PlatformUser.ID {
				return completePlatformAuthResult{}, errPlatformConflict
			}
		} else {
			channel, err = a.linkPlatformToChannel(
				ctx,
				channel,
				input.Platform,
				platformUserID,
				input.DefaultKickBotID,
			)
			if err != nil {
				return completePlatformAuthResult{}, fmt.Errorf("link platform to channel: %w", err)
			}
		}
	} else {
		sessionUserID = platformUserID
		channel, createdChannel, err = a.getOrCreateChannelForUser(
			ctx,
			platformUserID,
			input.Platform,
			input.DefaultBotID,
			input.DefaultKickBotID,
		)
		if err != nil {
			return completePlatformAuthResult{}, fmt.Errorf("get or create platform channel: %w", err)
		}
	}

	if err := a.upsertPlatformUserToken(ctx, platformUserID, input.Tokens); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("upsert platform user token: %w", err)
	}

	if createdChannel {
		if err = a.bus.Scheduler.CreateDefaultRoles.Publish(
			ctx,
			scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{channel.ID.String()}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default roles", logger.Error(err), slog.String("channel_id", channel.ID.String()))
		}

		if err = a.bus.Scheduler.CreateDefaultCommands.Publish(
			ctx,
			scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{channel.ID.String()}},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish create default commands", logger.Error(err), slog.String("channel_id", channel.ID.String()))
		}
	}

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		ctx,
		buscoreeventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channel.ID.String()},
	); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish eventsub subscribe", logger.Error(err), slog.String("channel_id", channel.ID.String()))
	}

	if err := a.sessions.SetSessionInternalUserID(ctx, sessionUserID); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set internal user id: %w", err)
	}

	sessionPlatform := input.Platform.String()
	if hasLiveSession {
		sessionPlatform = sessionUser.Platform.String()
	}

	if err := a.sessions.SetSessionCurrentPlatform(ctx, sessionPlatform); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set current platform: %w", err)
	}

	if err := a.sessions.SetSessionSelectedDashboard(ctx, channel.ID.String()); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set selected dashboard: %w", err)
	}

	return completePlatformAuthResult{
		SessionUserID:   sessionUserID,
		PlatformUserID:  platformUserID,
		Channel:         channel,
		CreatedUser:     createdUser,
		CreatedChannel:  createdChannel,
		UsedLiveSession: hasLiveSession,
	}, nil
}

func (a *Auth) getLiveSessionUser(ctx context.Context) (usersmodel.User, bool, error) {
	userID, err := a.sessions.GetInternalUserID(ctx)
	if err != nil {
		return usersmodel.Nil, false, nil
	}

	user, err := a.usersRepo.GetByID(ctx, userID.String())
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return usersmodel.Nil, false, nil
		}

		return usersmodel.Nil, false, err
	}

	return user, true, nil
}

func (a *Auth) getOrCreatePlatformUser(
	ctx context.Context,
	platform platformentity.Platform,
	platformUser *appplatform.PlatformUser,
) (usersmodel.User, bool, error) {
	user, err := a.usersRepo.GetByPlatformID(ctx, platform, platformUser.ID)
	userNotFound := errors.Is(err, usersmodel.ErrNotFound)
	if err != nil && !userNotFound {
		return usersmodel.Nil, false, err
	}

	if userNotFound {
		createdUser, createErr := a.usersRepo.Create(ctx, usersrepo.CreateInput{
			Platform:    platform,
			PlatformID:  platformUser.ID,
			IsBotAdmin:  false,
			IsBanned:    false,
			Login:       platformUser.Login,
			DisplayName: platformUser.DisplayName,
			Avatar:      platformUser.Avatar,
		})
		if createErr != nil {
			return usersmodel.Nil, false, createErr
		}

		return createdUser, true, nil
	}

	_, updateErr := a.usersRepo.Update(ctx, user.ID, usersrepo.UpdateInput{
		Login:       &platformUser.Login,
		DisplayName: &platformUser.DisplayName,
		Avatar:      &platformUser.Avatar,
	})
	if updateErr != nil {
		a.logger.ErrorContext(ctx, "cannot update platform user profile", logger.Error(updateErr), slog.String("user_id", user.ID), slog.String("platform", platform.String()))
	}

	return user, false, nil
}

func (a *Auth) getOrCreateChannelForUser(
	ctx context.Context,
	userID uuid.UUID,
	platform platformentity.Platform,
	defaultBotID string,
	defaultKickBotID *uuid.UUID,
) (channelsmodel.Channel, bool, error) {
	channel, err := a.getChannelByPlatformUserID(ctx, platform, userID)
	if err == nil {
		return channel, false, nil
	}

	if !errors.Is(err, channelsrepo.ErrNotFound) {
		return channelsmodel.Nil, false, err
	}

	switch platform {
	case platformentity.PlatformTwitch:
		channel, err = a.createChannel(ctx, &userID, nil, defaultBotID, nil)
	case platformentity.PlatformKick:
		channel, err = a.createChannel(ctx, nil, &userID, defaultBotID, defaultKickBotID)
	default:
		return channelsmodel.Nil, false, fmt.Errorf("unsupported platform: %s", platform)
	}
	if err != nil {
		return channelsmodel.Nil, false, err
	}

	return channel, true, nil
}

func (a *Auth) getChannelByPlatformUserID(
	ctx context.Context,
	platform platformentity.Platform,
	userID uuid.UUID,
) (channelsmodel.Channel, error) {
	switch platform {
	case platformentity.PlatformTwitch:
		return a.channelsRepo.GetByTwitchUserID(ctx, userID)
	case platformentity.PlatformKick:
		return a.channelsRepo.GetByKickUserID(ctx, userID)
	default:
		return channelsmodel.Nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

func (a *Auth) linkPlatformToChannel(
	ctx context.Context,
	channel channelsmodel.Channel,
	platform platformentity.Platform,
	platformUserID uuid.UUID,
	defaultKickBotID *uuid.UUID,
) (channelsmodel.Channel, error) {
	linkedUserID := linkedPlatformUserID(channel, platform)
	if linkedUserID != nil {
		if *linkedUserID == platformUserID {
			return channel, nil
		}

		return channelsmodel.Nil, errPlatformConflict
	}

	existingPlatformChannel, err := a.getChannelByPlatformUserID(ctx, platform, platformUserID)
	if err != nil && !errors.Is(err, channelsrepo.ErrNotFound) {
		return channelsmodel.Nil, err
	}

	if err == nil && existingPlatformChannel.ID != channel.ID {
		return channelsmodel.Nil, errPlatformConflict
	}

	updateInput := channelsrepo.UpdateInput{}
	switch platform {
	case platformentity.PlatformTwitch:
		updateInput.TwitchUserID = &platformUserID
	case platformentity.PlatformKick:
		updateInput.KickUserID = &platformUserID
		if channel.KickBotID == nil && defaultKickBotID != nil {
			updateInput.KickBotID = defaultKickBotID
		}
	default:
		return channelsmodel.Nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	return a.channelsRepo.Update(ctx, channel.ID, updateInput)
}

func linkedPlatformUserID(channel channelsmodel.Channel, platform platformentity.Platform) *uuid.UUID {
	switch platform {
	case platformentity.PlatformTwitch:
		return channel.TwitchUserID
	case platformentity.PlatformKick:
		return channel.KickUserID
	default:
		return nil
	}
}

func (a *Auth) upsertPlatformUserToken(
	ctx context.Context,
	userID uuid.UUID,
	tokens *appplatform.PlatformTokens,
) error {
	accessToken, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt access token: %w", err)
	}

	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt refresh token: %w", err)
	}

	currentToken, err := a.tokensRepository.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, tokensrepository.ErrNotFound) {
		return fmt.Errorf("get user token: %w", err)
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
			return fmt.Errorf("update user token: %w", err)
		}

		return nil
	}

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
		return fmt.Errorf("create user token: %w", err)
	}

	return nil
}
