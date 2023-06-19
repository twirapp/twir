package auth

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/crypto"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/auth"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Auth struct {
	*impl_deps.Deps
	TwitchScopes []string
}

func (c *Auth) AuthGetLink(ctx context.Context, request *auth.GetLinkRequest) (*auth.GetLinkResponse, error) {
	if request.State == "" {
		return nil, twirp.NewError(twirp.ErrorCode(400), "no state provided")
	}

	twitchClient, err := helix.NewClientWithContext(ctx, &helix.Options{
		ClientID:    c.Config.TwitchClientId,
		RedirectURI: c.Config.TwitchCallbackUrl,
	})
	if err != nil {
		return nil, err
	}

	url := twitchClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       c.TwitchScopes,
		State:        request.State,
		ForceVerify:  false,
	})

	return &auth.GetLinkResponse{Link: url}, nil
}

func (c *Auth) AuthPostCode(ctx context.Context, request *auth.PostCodeRequest) (*emptypb.Empty, error) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}
	tokens, err := twitchClient.RequestUserAccessToken(request.Code)
	if err != nil || tokens.ErrorMessage != "" {
		return nil, err
	}

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil || len(users.Data.Users) == 0 {
		return nil, err
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

	c.SessionManager.Put(ctx, "dbUser", dbUser)

	err = twirp.SetHTTPResponseHeader(ctx, "Dashboard-id", dbUser.ID)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &emptypb.Empty{}, nil
}
