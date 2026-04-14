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
	"github.com/nicklaw5/helix/v2"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	userplatformaccountsrepo "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	"github.com/twirapp/twir/libs/twitch"
)

type authBody struct {
	Code  string `json:"code" minLength:"20" required:"true"`
	State string `json:"state" required:"true"`
}

type authResponseDto struct {
	RedirectTo string `json:"redirect_to"`
}

func (a *Auth) handleAuthPostCode(
	ctx context.Context,
	input authBody,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	redirectTo, err := base64.StdEncoding.DecodeString(input.State)
	if err != nil {
		return nil, huma.Error400BadRequest("Cannot decode state", err)
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, a.config, a.bus)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot create twitch client", err)
	}

	tokens, err := twitchClient.RequestUserAccessToken(input.Code)
	if err != nil {
		return nil, err
	}

	if tokens.ErrorMessage != "" {
		return nil, huma.Error500InternalServerError(
			"Cannot get user access token",
			fmt.Errorf("error message: %s", tokens.ErrorMessage),
		)
	}

	twitchClient.SetUserAccessToken(tokens.Data.AccessToken)

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot get user data from twitch", err)
	}
	if len(users.Data.Users) == 0 {
		return nil, huma.Error500InternalServerError(
			"Cannot get user data from twitch",
			fmt.Errorf("twitch user not found"),
		)
	}

	twitchUser := users.Data.Users[0]

	accessToken, err := crypto.Encrypt(tokens.Data.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user access token", err)
	}

	refreshToken, err := crypto.Encrypt(tokens.Data.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user refresh token", err)
	}

	account, err := a.userPlatformAccountsRepo.GetByPlatformUserID(ctx, platform.PlatformTwitch, twitchUser.ID)
	accountNotFound := errors.Is(err, userplatformaccountsrepo.ErrNotFound)
	if err != nil && !accountNotFound {
		return nil, huma.Error500InternalServerError("Cannot get user platform account", err)
	}

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find default bot", err)
	}
	if defaultBot.ID == "" {
		return nil, huma.Error500InternalServerError("Cannot find default bot", fmt.Errorf("no default bot found"))
	}

	userID := uuid.Nil
	createdUser := false
	if accountNotFound {
		userID, err = a.createUser(ctx, &twitchUser.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create user", err)
		}
		createdUser = true
	} else {
		userID = account.UserID
	}

	_, err = a.userPlatformAccountsRepo.Upsert(ctx, userplatformaccountsrepo.UpsertInput{
		UserID:              userID,
		Platform:            platform.PlatformTwitch,
		PlatformUserID:      twitchUser.ID,
		PlatformLogin:       twitchUser.Login,
		PlatformDisplayName: twitchUser.DisplayName,
		PlatformAvatar:      twitchUser.ProfileImageURL,
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		Scopes:              tokens.Data.Scopes,
		ExpiresIn:           tokens.Data.ExpiresIn,
		ObtainmentTimestamp: time.Now().UTC(),
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot upsert user platform account", err)
	}

	dbUser, err := a.usersRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot get user", err)
	}
	isBanned := dbUser.IsBanned
	if isBanned {
		return nil, huma.Error403Forbidden("Forbidden", nil)
	}

	currentToken, err := a.tokensRepository.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, tokensrepository.ErrNotFound) {
		return nil, huma.Error500InternalServerError("Cannot get user token", err)
	}

	tokenExpires := tokens.Data.ExpiresIn
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
				Scopes:              tokens.Data.Scopes,
			},
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot update user token", err)
		}
	} else {
		_, err = a.tokensRepository.CreateUserToken(
			ctx,
			tokensrepository.CreateInput{
				UserID:              userID,
				AccessToken:         accessToken,
				RefreshToken:        refreshToken,
				ExpiresIn:           tokens.Data.ExpiresIn,
				ObtainmentTimestamp: tokenCreatedAt,
				Scopes:              tokens.Data.Scopes,
			},
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create user token", err)
		}
	}

	channel, err := a.channelsRepo.GetByUserIDAndPlatform(ctx, userID, platform.PlatformTwitch)
	if err != nil {
		if !errors.Is(err, channelsrepo.ErrNotFound) {
			return nil, huma.Error500InternalServerError("Cannot get channel", err)
		}

		channel, err = a.createChannel(ctx, userID, defaultBot.ID, platform.PlatformTwitch)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create channel", err)
		}
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
		return nil, huma.Error500InternalServerError("Cannot set internal user id", err)
	}

	if err := a.sessions.SetSessionCurrentPlatform(ctx, platform.PlatformTwitch.String()); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set current platform", err)
	}

	if err := a.sessions.SetSessionSelectedDashboard(ctx, channel.ID); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set selected dashboard", err)
	}

	if err := a.sessions.SetSessionTwitchUser(ctx, twitchUser); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set twitch user", err)
	}

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}
