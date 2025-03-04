package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
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

	twitchClient, err := twitch.NewAppClientWithContext(ctx, a.config, a.tokensGrpc)
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
	err = a.gorm.WithContext(ctx).Where("id = ?", twitchUser.ID).Preload("Token").Find(dbUser).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find user", err)
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
			IsTester:   false,
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

	tokenData := model.Tokens{
		ID:                  uuid.New().String(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		ExpiresIn:           int32(tokens.Data.ExpiresIn),
		ObtainmentTimestamp: time.Now().UTC(),
		Scopes:              tokens.Data.Scopes,
	}
	if dbUser.TokenID.Valid {
		tokenData.ID = dbUser.TokenID.String
	}

	if err := a.gorm.WithContext(ctx).Save(tokenData).Error; err != nil {
		return nil, huma.Error500InternalServerError("Cannot update user token", err)
	}

	if err := a.gorm.WithContext(ctx).Debug().Save(&tokenData).Error; err != nil {
		return nil, huma.Error500InternalServerError("Cannot update user token", err)
	}

	dbUser.TokenID = sql.NullString{
		String: tokenData.ID,
		Valid:  true,
	}

	if dbUser.Channel == nil || dbUser.Channel.ID == "" {
		dbUser.Channel = &model.Channels{
			ID:    twitchUser.ID,
			BotID: defaultBot.ID,
		}
	}

	if err := a.gorm.WithContext(ctx).Debug().Save(dbUser).Error; err != nil {
		return nil, huma.Error500InternalServerError("Cannot update user", err)
	}

	err = a.bus.Scheduler.CreateDefaultRoles.Publish(
		scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot create default roles", err)
	}

	err = a.bus.Scheduler.CreateDefaultCommands.Publish(
		scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot create default commands", err)
	}

	a.sessions.Put(ctx, "dbUser", &dbUser)
	a.sessions.Put(ctx, "twitchUser", &twitchUser)
	a.sessions.Put(ctx, "dashboardId", dbUser.ID)

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		eventsub.EventsubSubscribeToAllEventsRequest{
			ChannelID: dbUser.ID,
		},
	); err != nil {
		return nil, huma.Error500InternalServerError("Cannot subscribe to eventsub", err)
	}

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}
