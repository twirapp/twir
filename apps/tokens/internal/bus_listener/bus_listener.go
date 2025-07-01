package bus_listener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/nicklaw5/helix/v2"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var appTokenScopes []string

type appToken struct {
	AccessToken    string
	ObtainmentTime time.Time
	ExpiresIn      int
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config  cfg.Config
	Gorm    *gorm.DB
	Redsync *redsync.Redsync
	Logger  logger.Logger
	TwirBus *buscore.Bus
}

type tokensImpl struct {
	globalClient   *helix.Client
	appAccessToken *appToken

	db      *gorm.DB
	config  cfg.Config
	log     logger.Logger
	redSync *redsync.Redsync
	twirBus *buscore.Bus
}

func rateLimitFunc(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	var reset64 int64
	reset64 = int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().UTC().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}

func NewTokens(opts Opts) error {
	helixClient, err := helix.NewClient(
		&helix.Options{
			ClientID:      opts.Config.TwitchClientId,
			ClientSecret:  opts.Config.TwitchClientSecret,
			RedirectURI:   opts.Config.GetTwitchCallbackUrl(),
			RateLimitFunc: rateLimitFunc,
		},
	)
	if err != nil {
		return err
	}
	appAccessToken, err := helixClient.RequestAppAccessToken(appTokenScopes)
	if err != nil {
		return err
	}

	impl := &tokensImpl{
		globalClient: helixClient,
		appAccessToken: &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		},

		db:      opts.Gorm,
		config:  opts.Config,
		log:     opts.Logger,
		redSync: opts.Redsync,
		twirBus: opts.TwirBus,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.twirBus.Tokens.RequestAppToken.SubscribeGroup(
					"tokens",
					impl.RequestAppToken,
				); err != nil {
					return err
				}
				if err := impl.twirBus.Tokens.RequestUserToken.SubscribeGroup(
					"tokens",
					impl.RequestUserToken,
				); err != nil {
					return err
				}
				if err := impl.twirBus.Tokens.RequestBotToken.SubscribeGroup(
					"tokens",
					impl.RequestBotToken,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.twirBus.Tokens.RequestAppToken.Unsubscribe()
				impl.twirBus.Tokens.RequestUserToken.Unsubscribe()
				impl.twirBus.Tokens.RequestBotToken.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *tokensImpl) RequestAppToken(
	ctx context.Context,
	_ struct{},
) (tokens.TokenResponse, error) {
	mu := c.redSync.NewMutex("tokens-app-lock")
	mu.Lock()
	defer mu.Unlock()

	if isTokenExpired(c.appAccessToken.ExpiresIn, c.appAccessToken.ObtainmentTime) {
		appAccessToken, err := c.globalClient.RequestAppAccessToken(appTokenScopes)
		if err != nil {
			return tokens.TokenResponse{}, err
		}

		c.appAccessToken = &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		}
		c.log.Info("app token refreshed")
	}

	return tokens.TokenResponse{
		AccessToken: c.appAccessToken.AccessToken,
		Scopes:      []string{},
	}, nil
}

func (c *tokensImpl) RequestUserToken(
	ctx context.Context,
	data tokens.GetUserTokenRequest,
) (tokens.TokenResponse, error) {
	mu := c.redSync.NewMutex("tokens-users-lock-" + data.UserId)
	mu.Lock()
	defer mu.Unlock()

	user := model.Users{}
	err := c.db.WithContext(ctx).Where("id = ?", data.UserId).Preload("Token").Find(&user).Error
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if user.ID == "" || user.Token == nil || user.Token.ID == "" {
		return tokens.TokenResponse{}, fmt.Errorf(
			"cannot find user token in db, userId: %s, token: %v",
			user.ID,
			user.Token,
		)
	}

	decryptedRefreshToken, err := crypto.Decrypt(user.Token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if decryptedRefreshToken == "" {
		return tokens.TokenResponse{}, errors.New("refresh token is empty")
	}

	if isTokenExpired(int(user.Token.ExpiresIn), user.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return tokens.TokenResponse{}, err
		}
		if newToken.ErrorMessage != "" {
			return tokens.TokenResponse{}, fmt.Errorf("refresh token error: %s", newToken.ErrorMessage)
		}

		if newToken.StatusCode != 200 || newToken.Data.AccessToken == "" {
			return tokens.TokenResponse{}, fmt.Errorf(
				"refresh token status code: %d",
				newToken.StatusCode,
			)
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}
		user.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}
		user.Token.AccessToken = newAccessToken

		user.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		user.Token.Scopes = newToken.Data.Scopes
		user.Token.ObtainmentTimestamp = time.Now().UTC()
		if err := c.db.WithContext(ctx).Save(&user.Token).Error; err != nil {
			return tokens.TokenResponse{}, err
		}

		c.log.Info(
			"user token refreshed",
			slog.String("user_id", user.ID),
			slog.Int("expires_in", int(user.Token.ExpiresIn)),
			slog.String("access_token", newAccessToken),
			slog.String("refresh_token", newRefreshToken),
		)
	}

	decryptedAccessToken, err := crypto.Decrypt(user.Token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	return tokens.TokenResponse{
		AccessToken: decryptedAccessToken,
		Scopes:      user.Token.Scopes,
	}, nil
}

func (c *tokensImpl) RequestBotToken(
	ctx context.Context,
	data tokens.GetBotTokenRequest,
) (tokens.TokenResponse, error) {
	mu := c.redSync.NewMutex("tokens-bots-lock-" + data.BotId)
	mu.Lock()
	defer mu.Unlock()

	bot := model.Bots{}
	err := c.db.WithContext(ctx).Where("id = ?", data.BotId).Preload("Token").Find(&bot).Error
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if bot.ID == "" || bot.Token == nil || bot.Token.ID == "" {
		return tokens.TokenResponse{}, errors.New("cannot find bot token in db")
	}

	decryptedRefreshToken, err := crypto.Decrypt(bot.Token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if isTokenExpired(int(bot.Token.ExpiresIn), bot.Token.ObtainmentTimestamp) {
		newToken, err := c.globalClient.RefreshUserAccessToken(decryptedRefreshToken)
		if err != nil {
			return tokens.TokenResponse{}, err
		}

		if newToken.ErrorMessage != "" {
			return tokens.TokenResponse{}, fmt.Errorf("refresh token error: %s", newToken.ErrorMessage)
		}

		newRefreshToken, err := crypto.Encrypt(newToken.Data.RefreshToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}
		bot.Token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}
		bot.Token.AccessToken = newAccessToken

		bot.Token.ExpiresIn = int32(newToken.Data.ExpiresIn)
		bot.Token.Scopes = newToken.Data.Scopes
		bot.Token.ObtainmentTimestamp = time.Now().UTC()
		if err := c.db.WithContext(ctx).Save(&bot.Token).Error; err != nil {
			return tokens.TokenResponse{}, err
		}
		c.log.Info("bot token refreshed", slog.String("bot_id", bot.ID))
	}

	decryptedAccessToken, err := crypto.Decrypt(bot.Token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	return tokens.TokenResponse{
		AccessToken: decryptedAccessToken,
		Scopes:      bot.Token.Scopes,
		ExpiresIn:   bot.Token.ExpiresIn,
	}, nil
}
