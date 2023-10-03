package grpc_impl

import (
	"context"
	"errors"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/do"
	"github.com/satont/twir/apps/tokens/internal/di"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
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

	appLock   *redsync.Mutex
	usersLock *redsync.Mutex
	botsLock  *redsync.Mutex
}

func NewTokens() *TokensGrpcImpl {
	config := do.MustInvoke[cfg.Config](di.Provider)

	helixClient, err := helix.NewClient(
		&helix.Options{
			ClientID:     config.TwitchClientId,
			ClientSecret: config.TwitchClientSecret,
			RedirectURI:  config.TwitchCallbackUrl,
		},
	)
	if err != nil {
		panic(err)
	}

	appAccessToken, err := helixClient.RequestAppAccessToken(appTokenScopes)
	if err != nil {
		panic(err)
	}

	redisSync := do.MustInvoke[redsync.Redsync](di.Provider)

	return &TokensGrpcImpl{
		UnimplementedTokensServer: tokens.UnimplementedTokensServer{},
		globalClient:              helixClient,
		appAccessToken: &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		},

		botsLock:  redisSync.NewMutex("tokens-bots-lock"),
		usersLock: redisSync.NewMutex("tokens-users-lock"),
		appLock:   redisSync.NewMutex("tokens-app-lock"),
	}
}

func (c *TokensGrpcImpl) RequestAppToken(
	_ context.Context,
	_ *emptypb.Empty,
) (*tokens.Token, error) {
	c.appLock.Lock()
	defer c.appLock.Unlock()

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
	c.usersLock.Lock()
	defer c.usersLock.Unlock()

	db := do.MustInvoke[gorm.DB](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	user := model.Users{}
	err := db.Where("id = ?", data.UserId).Preload("Token").Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == "" || user.Token == nil || user.Token.ID == "" {
		return nil, errors.New("cannot find user token in db")
	}

	decryptedRefreshToken, err := crypto.Decrypt(user.Token.RefreshToken, config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	if isTokenExpired(int(user.Token.ExpiresIn), user.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return nil, err
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		user.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		user.Token.AccessToken = newAccessToken

		user.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		user.Token.Scopes = newToken.Data.Scopes
		user.Token.ObtainmentTimestamp = time.Now().UTC()
		db.Save(&user.Token)
	}

	decryptedAccessToken, err := crypto.Decrypt(user.Token.AccessToken, config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	return &tokens.Token{
		AccessToken: decryptedAccessToken,
		Scopes:      user.Token.Scopes,
	}, nil
}

func (c *TokensGrpcImpl) RequestBotToken(
	_ context.Context,
	data *tokens.GetBotTokenRequest,
) (*tokens.Token, error) {
	c.botsLock.Lock()
	defer c.botsLock.Unlock()

	db := do.MustInvoke[gorm.DB](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	bot := model.Bots{}
	err := db.Where("id = ?", data.BotId).Preload("Token").Find(&bot).Error
	if err != nil {
		return nil, err
	}

	if bot.ID == "" || bot.Token == nil || bot.Token.ID == "" {
		return nil, errors.New("cannot find bot token in db")
	}

	decryptedRefreshToken, err := crypto.Decrypt(bot.Token.RefreshToken, config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	if isTokenExpired(int(bot.Token.ExpiresIn), bot.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return nil, err
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		bot.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, config.TokensCipherKey)
		if err != nil {
			return nil, err
		}
		bot.Token.AccessToken = newAccessToken

		bot.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		bot.Token.Scopes = newToken.Data.Scopes
		bot.Token.ObtainmentTimestamp = time.Now().UTC()
		db.Save(&bot.Token)
	}

	decryptedAccessToken, err := crypto.Decrypt(bot.Token.AccessToken, config.TokensCipherKey)
	if err != nil {
		return nil, err
	}

	return &tokens.Token{
		AccessToken: decryptedAccessToken,
		Scopes:      bot.Token.Scopes,
		ExpiresIn:   bot.Token.ExpiresIn,
	}, nil
}
