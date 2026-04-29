package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
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

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find default bot", err)
	}
	if defaultBot.ID == "" {
		return nil, huma.Error500InternalServerError("Cannot find default bot", fmt.Errorf("no default bot found"))
	}

	result, err := a.completePlatformAuth(ctx, completePlatformAuthInput{
		Platform: platform.PlatformTwitch,
		PlatformUser: &appplatform.PlatformUser{
			ID:          twitchUser.ID,
			Login:       twitchUser.Login,
			DisplayName: twitchUser.DisplayName,
			Avatar:      twitchUser.ProfileImageURL,
		},
		Tokens: &appplatform.PlatformTokens{
			AccessToken:  tokens.Data.AccessToken,
			RefreshToken: tokens.Data.RefreshToken,
			ExpiresIn:    tokens.Data.ExpiresIn,
			Scopes:       tokens.Data.Scopes,
		},
		DefaultBotID: defaultBot.ID,
	})
	if err != nil {
		if errors.Is(err, errAuthForbidden) {
			return nil, huma.Error403Forbidden("Forbidden", nil)
		}

		if errors.Is(err, errPlatformConflict) {
			return nil, huma.Error409Conflict("Platform account already linked to another dashboard", err)
		}

		return nil, huma.Error500InternalServerError("Cannot complete auth", err)
	}

	if err := a.sessions.SetSessionTwitchUser(ctx, twitchUser); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set twitch user", err)
	}

	a.logger.InfoContext(ctx, "twitch auth: completed successfully", slog.String("channel_id", result.Channel.ID.String()), slog.String("user_id", result.SessionUserID.String()))

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}

func (a *Auth) createChannel(
	ctx context.Context,
	twitchUserID *uuid.UUID,
	kickUserID *uuid.UUID,
	botID string,
	kickBotID *uuid.UUID,
) (channelsmodel.Channel, error) {
	channel, err := a.channelsRepo.Create(ctx, channelsrepo.CreateInput{
		TwitchUserID: twitchUserID,
		KickUserID:   kickUserID,
		BotID:        botID,
		KickBotID:    kickBotID,
	})
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("create channel: %w", err)
	}

	return channel, nil
}
