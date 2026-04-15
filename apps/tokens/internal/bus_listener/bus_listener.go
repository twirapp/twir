package bus_listener

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"go.uber.org/fx"
	"gorm.io/gorm"

	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	twitchlib "github.com/twirapp/twir/libs/twitch"
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

	Config           cfg.Config
	Gorm             *gorm.DB
	Redsync          *redsync.Redsync
	Logger           *slog.Logger
	TwirBus          *buscore.Bus
	TokensRepository tokensrepository.Repository
	UsersRepository  usersrepository.Repository
}

type tokensImpl struct {
	globalClient   *helix.Client
	httpClient     *http.Client
	appAccessToken *appToken

	config           cfg.Config
	log              *slog.Logger
	redSync          *redsync.Redsync
	twirBus          *buscore.Bus
	tokensRepository tokensrepository.Repository
	usersRepository  usersrepository.Repository
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
	httpClient := &http.Client{}
	if opts.Config.TwitchMockEnabled {
		httpClient.Transport = twitchlib.NewMockRoundTripper(http.DefaultTransport, opts.Config)
	}

	helixOpts := &helix.Options{
		ClientID:      opts.Config.TwitchClientId,
		ClientSecret:  opts.Config.TwitchClientSecret,
		RedirectURI:   opts.Config.GetTwitchCallbackUrl(),
		RateLimitFunc: rateLimitFunc,
		HTTPClient:    httpClient,
	}
	if opts.Config.TwitchMockEnabled {
		helixOpts.APIBaseURL = opts.Config.TwitchMockApiUrl
	}

	helixClient, err := helix.NewClient(
		helixOpts,
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
		httpClient:   httpClient,
		appAccessToken: &appToken{
			AccessToken:    appAccessToken.Data.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      appAccessToken.Data.ExpiresIn,
		},

		config:           opts.Config,
		log:              opts.Logger,
		redSync:          opts.Redsync,
		twirBus:          opts.TwirBus,
		tokensRepository: opts.TokensRepository,
		usersRepository:  opts.UsersRepository,
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

	userID, err := uuid.Parse(data.UserId)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("cannot parse user id: %w", err)
	}

	token, err := c.tokensRepository.GetByUserID(ctx, userID)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf(
			"cannot get user token from repository: %w",
			err,
		)
	}

	decryptedRefreshToken, err := crypto.Decrypt(token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if decryptedRefreshToken == "" {
		return tokens.TokenResponse{}, errors.New("refresh token is empty")
	}

	if isTokenExpired(token.ExpiresIn, token.ObtainmentTimestamp) {
		user, err := c.usersRepository.GetByID(ctx, data.UserId)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("cannot get user: %w", err)
		}

		var refreshedAccessToken, refreshedRefreshToken string
		var refreshedExpiresIn int
		var refreshedScopes []string

		switch platformentity.Platform(user.Platform) {
		case platformentity.PlatformKick:
			kickTokens, err := c.refreshKickToken(ctx, decryptedRefreshToken)
			if err != nil {
				return tokens.TokenResponse{}, fmt.Errorf("kick refresh token: %w", err)
			}
			refreshedAccessToken = kickTokens.AccessToken
			refreshedRefreshToken = kickTokens.RefreshToken
			refreshedExpiresIn = kickTokens.ExpiresIn
			refreshedScopes = kickTokens.Scopes
		default:
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
			refreshedAccessToken = newToken.Data.AccessToken
			refreshedRefreshToken = newToken.Data.RefreshToken
			refreshedExpiresIn = newToken.Data.ExpiresIn
			refreshedScopes = newToken.Data.Scopes
		}

		newRefreshToken, err := crypto.Encrypt(refreshedRefreshToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}

		newAccessToken, err := crypto.Encrypt(refreshedAccessToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}

		timeStamp := time.Now().UTC()

		dbToken, err := c.tokensRepository.UpdateTokenByID(
			ctx, token.ID, tokensrepository.UpdateTokenInput{
				AccessToken:         &newAccessToken,
				RefreshToken:        &newRefreshToken,
				ExpiresIn:           &refreshedExpiresIn,
				ObtainmentTimestamp: &timeStamp,
				Scopes:              refreshedScopes,
			},
		)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf(
				"cannot update user token in repository: %w",
				err,
			)
		}

		token = dbToken

		c.log.Info(
			"user token refreshed",
			slog.String("user_id", data.UserId),
			slog.String("platform", user.Platform),
			slog.Int("expires_in", token.ExpiresIn),
		)
	}

	decryptedAccessToken, err := crypto.Decrypt(token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	return tokens.TokenResponse{
		AccessToken: decryptedAccessToken,
		Scopes:      token.Scopes,
	}, nil
}

type kickTokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
}

func (c *tokensImpl) refreshKickToken(ctx context.Context, refreshToken string) (*kickTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", c.config.KickClientId)
	form.Set("client_secret", c.config.KickClientSecret)
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://id.kick.com/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("create kick refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do kick refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read kick refresh response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("kick refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		ExpiresIn    int      `json:"expires_in"`
		Scope        string   `json:"scope"`
		Scopes       []string `json:"scopes"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode kick refresh response: %w", err)
	}

	responseScopes := parsed.Scopes
	if len(responseScopes) == 0 && parsed.Scope != "" {
		responseScopes = strings.Fields(parsed.Scope)
	}

	return &kickTokenResponse{
		AccessToken:  parsed.AccessToken,
		RefreshToken: parsed.RefreshToken,
		ExpiresIn:    parsed.ExpiresIn,
		Scopes:       responseScopes,
	}, nil
}

func (c *tokensImpl) RequestBotToken(
	ctx context.Context,
	data tokens.GetBotTokenRequest,
) (tokens.TokenResponse, error) {
	mu := c.redSync.NewMutex("tokens-bots-lock-" + data.BotId)
	mu.Lock()
	defer mu.Unlock()

	token, err := c.tokensRepository.GetByBotID(ctx, data.BotId)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf(
			"cannot get bot token from repository: %w",
			err,
		)
	}

	decryptedRefreshToken, err := crypto.Decrypt(token.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	if isTokenExpired(token.ExpiresIn, token.ObtainmentTimestamp) {
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
		token.RefreshToken = newRefreshToken

		newAccessToken, err := crypto.Encrypt(newToken.Data.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, err
		}

		timeStamp := time.Now().UTC()

		newDbToken, err := c.tokensRepository.UpdateTokenByID(
			ctx, token.ID, tokensrepository.UpdateTokenInput{
				AccessToken:         &newAccessToken,
				RefreshToken:        &newRefreshToken,
				ExpiresIn:           &newToken.Data.ExpiresIn,
				ObtainmentTimestamp: &timeStamp,
				Scopes:              newToken.Data.Scopes,
			},
		)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf(
				"cannot update bot token in repository: %w",
				err,
			)
		}

		token = newDbToken

		c.log.Info("bot token refreshed", slog.String("bot_id", data.BotId))
	}

	decryptedAccessToken, err := crypto.Decrypt(token.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, err
	}

	return tokens.TokenResponse{
		AccessToken: decryptedAccessToken,
		Scopes:      token.Scopes,
		ExpiresIn:   int32(token.ExpiresIn),
	}, nil
}
