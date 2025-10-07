package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"

	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
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

	dbUser := &model.Users{}
	err = a.gorm.WithContext(ctx).Where("id = ?", twitchUser.ID).Find(dbUser).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find user", err)
	}

	if dbUser.IsBanned {
		return nil, huma.Error403Forbidden("Forbidden", nil)
	}

	defaultBot := &model.Bots{}
	err = a.gorm.WithContext(ctx).Where("type = ?", "DEFAULT").Find(defaultBot).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find default bot", err)
	}

	if defaultBot.ID == "" {
		return nil, huma.Error500InternalServerError(
			"Cannot find default bot",
			fmt.Errorf("no default bot found"),
		)
	}

	accessToken, err := crypto.Encrypt(tokens.Data.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user access token", err)
	}

	refreshToken, err := crypto.Encrypt(tokens.Data.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user refresh token", err)
	}

	if dbUser.ID == "" {
		newUser := &model.Users{
			ID:         twitchUser.ID,
			IsBotAdmin: false,
			ApiKey:     uuid.NewString(),
			Channel: &model.Channels{
				ID:    twitchUser.ID,
				BotID: defaultBot.ID,
			},
		}

		if err := a.gorm.Create(newUser).Error; err != nil {
			return nil, huma.Error500InternalServerError("Cannot create user", err)
		}

		dbUser = newUser
	}

	currentToken, err := a.tokensRepository.GetByUserID(ctx, dbUser.ID)
	if err != nil && !errors.Is(err, tokensrepository.ErrNotFound) {
		return nil, huma.Error500InternalServerError("Cannot get user token", err)
	}

	tokenExpires := tokens.Data.ExpiresIn
	tokenCreatedAt := time.Now().UTC()
	if currentToken != nil {
		_, err := a.tokensRepository.UpdateTokenByID(
			ctx, currentToken.ID, tokensrepository.UpdateTokenInput{
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
		newToken, err := a.tokensRepository.CreateUserToken(
			ctx, tokensrepository.CreateInput{
				UserID:              dbUser.ID,
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

		dbUser.TokenID = sql.NullString{
			String: newToken.ID.String(),
			Valid:  true,
		}
	}

	if dbUser.Channel == nil || dbUser.Channel.ID == "" {
		dbUser.Channel = &model.Channels{
			ID:    twitchUser.ID,
			BotID: defaultBot.ID,
		}
	}

	if err := a.gorm.WithContext(ctx).Save(dbUser).Error; err != nil {
		return nil, huma.Error500InternalServerError("Cannot update user", err)
	}

	err = a.bus.Scheduler.CreateDefaultRoles.Publish(
		ctx,
		scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot create default roles", err)
	}

	err = a.bus.Scheduler.CreateDefaultCommands.Publish(
		ctx,
		scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot create default commands", err)
	}

	if err := a.sessions.SetSessionAuthenticatedUser(ctx, *dbUser); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set session user", err)
	}

	if err := a.sessions.SetSessionSelectedDashboard(ctx, dbUser.ID); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set session selected dashboard", err)
	}

	if err := a.sessions.SetSessionTwitchUser(ctx, twitchUser); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set session twitch user", err)
	}

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		ctx,
		eventsub.EventsubSubscribeToAllEventsRequest{
			ChannelID: dbUser.ID,
		},
	); err != nil {
		return nil, huma.Error500InternalServerError("Cannot subscribe to eventsub", err)
	}

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}
