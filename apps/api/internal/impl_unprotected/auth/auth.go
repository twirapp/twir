package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/auth"
	"github.com/twitchtv/twirp"
)

type Auth struct {
	*impl_deps.Deps
	TwitchScopes []string
}

func (c *Auth) AuthPostCode(ctx context.Context, request *auth.PostCodeRequest) (
	*auth.PostCodeResponse,
	error,
) {
	if request.GetCode() == "" {
		return nil, twirp.NewError("400", "no code provided")
	}
	if request.GetState() == "" {
		return nil, twirp.NewError("400", "no state provided")
	}

	redirectTo, err := base64.StdEncoding.DecodeString(request.GetState())
	if err != nil {
		return nil, fmt.Errorf("cannot decode state: %w", err)
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, fmt.Errorf("cannot create twitch client: %w", err)
	}
	tokens, err := twitchClient.RequestUserAccessToken(request.GetCode())
	if err != nil {
		return nil, fmt.Errorf("cannot user data from twitch: %w", err)
	}
	if tokens.ErrorMessage != "" {
		return nil, errors.New(tokens.ErrorMessage)
	}

	twitchClient.SetUserAccessToken(tokens.Data.AccessToken)

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, err
	}
	if len(users.Data.Users) == 0 {
		return nil, errors.New("twitch user not found")
	}

	twitchUser := users.Data.Users[0]

	dbUser := &model.Users{}
	err = c.Db.WithContext(ctx).Where("id = ?", twitchUser.ID).Preload("Token").Find(dbUser).Error
	if err != nil {
		return nil, err
	}

	defaultBot := &model.Bots{}
	err = c.Db.WithContext(ctx).Where("type = ?", "DEFAULT").Find(defaultBot).Error
	if err != nil {
		return nil, err
	}

	if defaultBot.ID == "" {
		return nil, twirp.Internal.Error("no default bot found")
	}

	accessToken, err := crypto.Encrypt(tokens.Data.AccessToken, c.Config.TokensCipherKey)
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt user access token: %w", err)
	}

	refreshToken, err := crypto.Encrypt(tokens.Data.RefreshToken, c.Config.TokensCipherKey)
	if err != nil {
		return nil, fmt.Errorf("ecnrypt user refres token: %w", err)
	}

	if dbUser.ID == "" {
		newUser := &model.Users{
			ID:         twitchUser.ID,
			IsTester:   false,
			IsBotAdmin: false,
			ApiKey:     uuid.New().String(),
			Channel: &model.Channels{
				ID:    twitchUser.ID,
				BotID: defaultBot.ID,
			},
		}

		if err := c.Db.Create(newUser).Error; err != nil {
			return nil, fmt.Errorf("cannot create user: %w", err)
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

	if err := c.Db.WithContext(ctx).Save(tokenData).Error; err != nil {
		return nil, fmt.Errorf("cannot update user token: %w", err)
	}

	if err := c.Db.WithContext(ctx).Debug().Save(&tokenData).Error; err != nil {
		return nil, fmt.Errorf("cannot update db user: %w", err)
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

	if err := c.Db.WithContext(ctx).Debug().Save(dbUser).Error; err != nil {
		return nil, fmt.Errorf("cannot update db user: %w", err)
	}

	err = c.Bus.Scheduler.CreateDefaultRoles.Publish(
		scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, err
	}

	err = c.Bus.Scheduler.CreateDefaultCommands.Publish(
		scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, err
	}

	c.SessionManager.Put(ctx, "dbUser", &dbUser)
	c.SessionManager.Put(ctx, "twitchUser", &twitchUser)
	c.SessionManager.Put(ctx, "dashboardId", dbUser.ID)

	if err := c.Bus.EventSub.SubscribeToAllEvents.Publish(
		eventsub.EventsubSubscribeToAllEventsRequest{
			ChannelID: dbUser.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("cannot subscribe to eventsub: %w", err)
	}

	return &auth.PostCodeResponse{
		RedirectTo: string(redirectTo),
	}, nil
}
