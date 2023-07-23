package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/auth"
	"github.com/satont/twir/libs/twitch"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth struct {
	*impl_deps.Deps
	TwitchScopes []string
}

func (c *Auth) AuthGetLink(ctx context.Context, request *auth.GetLinkRequest) (*auth.GetLinkResponse, error) {
	if request.State == "" {
		return nil, twirp.NewError(twirp.ErrorCode(400), "no state provided")
	}

	twitchClient, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:    c.Config.TwitchClientId,
			RedirectURI: c.Config.TwitchCallbackUrl,
		},
	)
	if err != nil {
		return nil, err
	}

	url := twitchClient.GetAuthorizationURL(
		&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       c.TwitchScopes,
			State:        request.State,
			ForceVerify:  false,
		},
	)

	return &auth.GetLinkResponse{Link: url}, nil
}

func (c *Auth) AuthPostCode(ctx context.Context, request *auth.PostCodeRequest) (*emptypb.Empty, error) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}
	tokens, err := twitchClient.RequestUserAccessToken(request.Code)
	if err != nil {
		return nil, err
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
	err = c.Db.WithContext(ctx).Where("id = ?", twitchUser.ID).Find(dbUser).Error
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

	if dbUser.ID == "" {
		accessToken, err := crypto.Encrypt(tokens.Data.AccessToken, c.Config.TokensCipherKey)
		if err != nil {
			return nil, err
		}

		refreshToken, err := crypto.Encrypt(tokens.Data.RefreshToken, c.Config.TokensCipherKey)
		if err != nil {
			return nil, err
		}

		tokenData := model.Tokens{
			ID:                  uuid.New().String(),
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			ExpiresIn:           int32(tokens.Data.ExpiresIn),
			ObtainmentTimestamp: time.Now().UTC(),
			Scopes:              tokens.Data.Scopes,
		}

		err = c.Db.WithContext(ctx).Create(&tokenData).Error
		if err != nil {
			return nil, err
		}

		newUser := &model.Users{
			ID:         twitchUser.ID,
			TokenID:    sql.NullString{String: tokenData.ID, Valid: true},
			IsTester:   false,
			IsBotAdmin: false,
			ApiKey:     uuid.New().String(),
			Channel: &model.Channels{
				ID:             twitchUser.ID,
				IsEnabled:      false,
				IsTwitchBanned: false,
				IsBanned:       false,
				IsBotMod:       false,
				BotID:          defaultBot.ID,
			},
		}
		err = c.Db.WithContext(ctx).Create(newUser).Error
		if err != nil {
			return nil, err
		}

		dbUser = newUser
	}
	if dbUser.Channel == nil || dbUser.Channel.ID == "" {
		dbUser.Channel = &model.Channels{
			ID:    twitchUser.ID,
			BotID: defaultBot.ID,
		}
	}

	_, err = c.Grpc.Scheduler.CreateDefaultRoles(
		ctx,
		&scheduler.CreateDefaultRolesRequest{UsersIds: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, err
	}

	_, err = c.Grpc.Scheduler.CreateDefaultCommands(
		ctx,
		&scheduler.CreateDefaultCommandsRequest{UsersIds: []string{twitchUser.ID}},
	)
	if err != nil {
		return nil, err
	}

	c.SessionManager.Put(ctx, "dbUser", &dbUser)
	c.SessionManager.Put(ctx, "twitchUser", &twitchUser)
	c.SessionManager.Put(ctx, "dashboardId", dbUser.ID)

	return &emptypb.Empty{}, nil
}
