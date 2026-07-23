package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

var errPlatformConflict = errors.New("platform account is already linked to another channel")
var errAuthForbidden = errors.New("forbidden")
var errPlatformUnavailable = errors.New("platform is not available")

type transactionRunner interface {
	Do(context.Context, func(context.Context) error) error
}

type eventSubPublisher interface {
	Publish(context.Context, buscoreeventsub.EventsubSubscribeToAllEventsRequest) error
}

type platformBindingConfig struct {
	BotUserID *uuid.UUID
	BotConfig json.RawMessage
}

type platformBindingConfigResolver func(context.Context) (platformBindingConfig, error)

type postPlatformAuthHook func(
	context.Context,
	completePlatformAuthResult,
	*appplatform.PlatformUser,
	*appplatform.PlatformTokens,
) error

type completePlatformAuthInput struct {
	Platform        platformentity.Platform
	PlatformUser    *appplatform.PlatformUser
	Tokens          *appplatform.PlatformTokens
	BindingConfig   platformBindingConfig
	TargetChannelID *uuid.UUID
	InitiatorUserID *uuid.UUID
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
	if a.transactionRunner == nil {
		return completePlatformAuthResult{}, fmt.Errorf("auth transaction runner is not configured")
	}

	sessionUser, hasLiveSession, err := a.getLiveSessionUser(ctx)
	if err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("get live session user: %w", err)
	}
	if hasLiveSession && sessionUser.IsBanned {
		return completePlatformAuthResult{}, errAuthForbidden
	}
	if input.TargetChannelID != nil {
		if input.InitiatorUserID == nil ||
			*input.InitiatorUserID == uuid.Nil ||
			!hasLiveSession ||
			sessionUser.ID != *input.InitiatorUserID {
			return completePlatformAuthResult{}, errAuthForbidden
		}
		if err := a.authorizeTargetDashboard(ctx, sessionUser, *input.TargetChannelID); err != nil {
			return completePlatformAuthResult{}, err
		}
	}

	platformUser, createdUser, err := a.getOrCreatePlatformUser(ctx, input.Platform, input.PlatformUser)
	if err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("get or create platform user: %w", err)
	}
	if platformUser.IsBanned {
		return completePlatformAuthResult{}, errAuthForbidden
	}

	result := completePlatformAuthResult{
		PlatformUserID:  platformUser.ID,
		CreatedUser:     createdUser,
		UsedLiveSession: hasLiveSession,
	}
	if hasLiveSession {
		result.SessionUserID = sessionUser.ID
	} else {
		result.SessionUserID = platformUser.ID
	}

	err = a.transactionRunner.Do(ctx, func(txCtx context.Context) error {
		if input.TargetChannelID != nil {
			channel, getChannelErr := a.channelsRepo.GetByID(txCtx, *input.TargetChannelID)
			if getChannelErr != nil {
				return fmt.Errorf("get target channel: %w", getChannelErr)
			}
			result.Channel = channel
		} else if hasLiveSession {
			channel, createdChannel, getChannelErr := a.getOrCreateChannelForUser(txCtx, sessionUser)
			if getChannelErr != nil {
				return fmt.Errorf("get or create session channel: %w", getChannelErr)
			}
			result.Channel = channel
			result.CreatedChannel = createdChannel
		} else {
			channel, createChannelErr := a.createChannel(txCtx)
			if createChannelErr != nil {
				return fmt.Errorf("create platform channel: %w", createChannelErr)
			}
			result.Channel = channel
			result.CreatedChannel = true
		}

		channel, linkErr := a.linkPlatformToChannel(
			txCtx,
			result.Channel,
			input.Platform,
			platformUser.ID,
			input.PlatformUser.ID,
			input.BindingConfig,
		)
		if linkErr != nil {
			return fmt.Errorf("link platform to channel: %w", linkErr)
		}
		result.Channel = channel

		return nil
	})
	if err != nil {
		return completePlatformAuthResult{}, err
	}

	if err := a.upsertPlatformUserToken(ctx, result.PlatformUserID, input.Tokens); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("upsert platform user token: %w", err)
	}

	if result.CreatedChannel {
		a.publishDefaultChannelResources(ctx, result.Channel.ID)
	}
	a.publishEventSubSubscription(ctx, result.Channel.ID, input.Platform)

	if err := a.sessions.SetSessionInternalUserID(ctx, result.SessionUserID); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set internal user id: %w", err)
	}

	sessionPlatform := input.Platform.String()
	if hasLiveSession {
		sessionPlatform = sessionUser.Platform.String()
	}
	if err := a.sessions.SetSessionCurrentPlatform(ctx, sessionPlatform); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set current platform: %w", err)
	}
	if err := a.sessions.SetSessionSelectedDashboard(ctx, result.Channel.ID.String()); err != nil {
		return completePlatformAuthResult{}, fmt.Errorf("set selected dashboard: %w", err)
	}

	return result, nil
}

func (a *Auth) authorizeTargetDashboard(
	ctx context.Context,
	user usersmodel.User,
	channelID uuid.UUID,
) error {
	if a.dashboardAccess == nil {
		return fmt.Errorf("dashboard access service is not configured")
	}

	hasAccess, err := a.dashboardAccess.CanAccess(ctx, dashboardaccess.Subject{
		ID:         user.ID.String(),
		IsBotAdmin: user.IsBotAdmin,
	}, channelID, "")
	if err != nil {
		return fmt.Errorf("check target dashboard access: %w", err)
	}
	if !hasAccess {
		return errAuthForbidden
	}

	return nil
}

func (a *Auth) getLiveSessionUser(ctx context.Context) (usersmodel.User, bool, error) {
	userID, err := a.sessions.GetInternalUserID(ctx)
	if err != nil {
		return usersmodel.Nil, false, nil
	}

	user, err := a.usersRepo.GetByID(ctx, userID)
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
		a.logger.ErrorContext(ctx, "cannot update platform user profile", logger.Error(updateErr), slog.String("user_id", user.ID.String()), slog.String("platform", platform.String()))
	}

	return user, false, nil
}

func (a *Auth) getOrCreateChannelForUser(
	ctx context.Context,
	user usersmodel.User,
) (channelsmodel.Channel, bool, error) {
	channel, err := a.channelsRepo.GetByBindingUserID(ctx, user.Platform, user.ID)
	if err == nil {
		return channel, false, nil
	}
	if !errors.Is(err, channelsrepo.ErrNotFound) {
		return channelsmodel.Nil, false, err
	}
	bindingConfig, err := a.platformBindingConfig(ctx, user.Platform)
	if err != nil {
		return channelsmodel.Nil, false, fmt.Errorf("get %s binding configuration: %w", user.Platform, err)
	}

	channel, err = a.createChannel(ctx)
	if err != nil {
		return channelsmodel.Nil, false, err
	}

	channel, err = a.linkPlatformToChannel(ctx, channel, user.Platform, user.ID, user.PlatformID, bindingConfig)
	if err != nil {
		return channelsmodel.Nil, false, err
	}

	return channel, true, nil
}

func (a *Auth) createChannel(ctx context.Context) (channelsmodel.Channel, error) {
	if a.botsRepo == nil {
		return channelsmodel.Nil, fmt.Errorf("bots repository is not configured")
	}

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("get default bot: %w", err)
	}
	if defaultBot.ID == "" {
		return channelsmodel.Nil, fmt.Errorf("default bot not found")
	}

	channel, err := a.channelsRepo.Create(ctx, channelsrepo.CreateInput{BotID: defaultBot.ID})
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("create channel: %w", err)
	}

	return channel, nil
}

func (a *Auth) linkPlatformToChannel(
	ctx context.Context,
	channel channelsmodel.Channel,
	platform platformentity.Platform,
	platformUserID uuid.UUID,
	platformChannelID string,
	bindingConfig platformBindingConfig,
) (channelsmodel.Channel, error) {
	binding, err := a.channelPlatformsRepo.GetByChannelAndPlatform(ctx, channel.ID, platform)
	if err == nil {
		if binding.UserID != platformUserID || binding.PlatformChannelID != platformChannelID {
			return channelsmodel.Nil, errPlatformConflict
		}

		return channel, nil
	}
	if !errors.Is(err, channelplatformsrepo.ErrNotFound) {
		return channelsmodel.Nil, err
	}

	existingPlatformChannel, err := a.channelsRepo.GetByBindingUserID(ctx, platform, platformUserID)
	if err != nil && !errors.Is(err, channelsrepo.ErrNotFound) {
		return channelsmodel.Nil, err
	}
	if err == nil && existingPlatformChannel.ID != channel.ID {
		return channelsmodel.Nil, errPlatformConflict
	}

	if len(bindingConfig.BotConfig) == 0 {
		bindingConfig.BotConfig = json.RawMessage(`{}`)
	}
	createdBinding, err := a.channelPlatformsRepo.Create(ctx, channelplatformsrepo.CreateInput{
		ChannelID:         channel.ID,
		Platform:          platform,
		UserID:            platformUserID,
		PlatformChannelID: platformChannelID,
		Enabled:           true,
		BotUserID:         bindingConfig.BotUserID,
		BotConfig:         bindingConfig.BotConfig,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return channelsmodel.Nil, errPlatformConflict
		}
		return channelsmodel.Nil, err
	}

	updatedChannel := channel
	updatedChannel.Bindings = append(
		append([]channelplatformsmodel.ChannelPlatform(nil), channel.Bindings...),
		createdBinding,
	)

	return updatedChannel, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (a *Auth) publishDefaultChannelResources(ctx context.Context, channelID uuid.UUID) {
	if a.bus == nil || a.bus.Scheduler == nil {
		return
	}

	if err := a.bus.Scheduler.CreateDefaultRoles.Publish(
		ctx,
		scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{channelID.String()}},
	); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish create default roles", logger.Error(err), slog.String("channel_id", channelID.String()))
	}
	if err := a.bus.Scheduler.CreateDefaultCommands.Publish(
		ctx,
		scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{channelID.String()}},
	); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish create default commands", logger.Error(err), slog.String("channel_id", channelID.String()))
	}
}

func (a *Auth) publishEventSubSubscription(ctx context.Context, channelID uuid.UUID, platform platformentity.Platform) {
	if a.eventSubPublisher == nil {
		return
	}

	if err := a.eventSubPublisher.Publish(ctx, buscoreeventsub.EventsubSubscribeToAllEventsRequest{
		ChannelID: channelID.String(),
		Platform:  platform,
	}); err != nil {
		a.logger.ErrorContext(ctx, "cannot publish eventsub subscribe", logger.Error(err), slog.String("channel_id", channelID.String()))
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

	var encryptedDeviceID *string
	if tokens.DeviceID != "" {
		deviceID, encryptErr := crypto.Encrypt(tokens.DeviceID, a.config.TokensCipherKey)
		if encryptErr != nil {
			return fmt.Errorf("encrypt device ID: %w", encryptErr)
		}
		encryptedDeviceID = &deviceID
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
				DeviceID:            encryptedDeviceID,
			},
		)
		if err != nil {
			return fmt.Errorf("update user token: %w", err)
		}

		return nil
	}

	createdToken, err := a.tokensRepository.CreateUserToken(
		ctx,
		tokensrepository.CreateInput{
			UserID:              userID,
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			ExpiresIn:           tokens.ExpiresIn,
			ObtainmentTimestamp: tokenCreatedAt,
			Scopes:              tokens.Scopes,
			DeviceID:            encryptedDeviceID,
		},
	)
	if err != nil {
		return fmt.Errorf("create user token: %w", err)
	}

	tokenID := createdToken.ID.String()
	_, err = a.usersRepo.Update(ctx, userID, usersrepo.UpdateInput{TokenID: &tokenID})
	if err != nil {
		return fmt.Errorf("bind user token: %w", err)
	}

	return nil
}
