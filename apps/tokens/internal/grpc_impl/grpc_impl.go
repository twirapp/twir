package grpc_impl

import (
	"context"
	"errors"
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/tokens/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

var appTokenScopes = []string{}

type appToken struct {
	AccessToken    string
	ObtainmentTime time.Time
	ExpiresIn      int
}

type TokensGrpcImpl struct {
	tokens.UnimplementedTokensServer

	globalClient   *helix.Client
	appAccessToken *appToken
}

func NewTokens() *TokensGrpcImpl {
	config := do.MustInvoke[cfg.Config](di.Provider)

	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:     config.TwitchClientId,
		ClientSecret: config.TwitchClientSecret,
		RedirectURI:  config.TwitchCallbackUrl,
	})

	if err != nil {
		panic(err)
	}

	appAccessToken, err := helixClient.RequestAppAccessToken(appTokenScopes)

	if err != nil {
		panic(err)
	}

	return &TokensGrpcImpl{
		UnimplementedTokensServer: tokens.UnimplementedTokensServer{},

		globalClient: helixClient,
		appAccessToken: &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		},
	}
}

func (c *TokensGrpcImpl) RequestAppToken(
	_ context.Context,
	_ *emptypb.Empty,
) (*tokens.Token, error) {
	if isTokenExpired(c.appAccessToken.ExpiresIn, c.appAccessToken.ObtainmentTime) {
		appAccessToken, err := c.globalClient.RequestAppAccessToken(appTokenScopes)

		if err != nil {
			return nil, err
		}

		c.appAccessToken = &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		}
	}

	return &tokens.Token{
		AccessToken: c.appAccessToken.AccessToken,
		Scopes:      []string{},
	}, nil
}

func (c *TokensGrpcImpl) RequestUserToken(
	_ context.Context,
	data *tokens.GetUserTokenRequest,
) (*tokens.Token, error) {
	db := do.MustInvoke[gorm.DB](di.Provider)

	user := model.Users{}
	err := db.Where("id = ?", data.UserId).Preload("Token").Find(&user).Error

	if err != nil {
		return nil, err
	}

	if user.ID == "" || user.Token == nil || user.Token.ID == "" {
		return nil, errors.New("cannot find user token in db")
	}

	if isTokenExpired(int(user.Token.ExpiresIn), user.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(user.Token.RefreshToken)

		if err != nil {
			return nil, err
		}

		user.Token.RefreshToken = newToken.Data.RefreshToken
		user.Token.AccessToken = newToken.Data.AccessToken
		user.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		user.Token.Scopes = newToken.Data.Scopes
		user.Token.ObtainmentTimestamp = time.Now().UTC()
		db.Save(&user.Token)
	}

	return &tokens.Token{
		AccessToken: user.Token.AccessToken,
		Scopes:      user.Token.Scopes,
	}, nil
}

func (c *TokensGrpcImpl) RequestBotToken(
	_ context.Context,
	data *tokens.GetBotTokenRequest,
) (*tokens.Token, error) {
	db := do.MustInvoke[gorm.DB](di.Provider)

	bot := model.Bots{}
	err := db.Where("id = ?", data.BotId).Preload("Token").Find(&bot).Error

	if err != nil {
		return nil, err
	}

	if bot.ID == "" || bot.Token == nil || bot.Token.ID == "" {
		return nil, errors.New("cannot find bot token in db")
	}

	if isTokenExpired(int(bot.Token.ExpiresIn), bot.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(bot.Token.RefreshToken)

		if err != nil {
			return nil, err
		}

		bot.Token.RefreshToken = newToken.Data.RefreshToken
		bot.Token.AccessToken = newToken.Data.AccessToken
		bot.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		bot.Token.Scopes = newToken.Data.Scopes
		bot.Token.ObtainmentTimestamp = time.Now().UTC()
		db.Save(&bot.Token)
	}

	return &tokens.Token{
		AccessToken: bot.Token.AccessToken,
		Scopes:      bot.Token.Scopes,
	}, nil
}
