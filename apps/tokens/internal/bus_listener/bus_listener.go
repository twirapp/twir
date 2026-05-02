package bus_listener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/scorfly/gokick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsintegrationsrepository "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsspotifyrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"go.uber.org/fx"
	"gorm.io/gorm"

	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
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
	KickBotsRepo     kickbotsrepository.Repository
	IntegrationsRepo integrationsrepository.Repository
	ChannelIntegrationsRepo channelsintegrationsrepository.Repository
	SpotifyIntegrationsRepo channelsintegrationsspotifyrepository.Repository
	TokensRepository tokensrepository.Repository
	UsersRepository  usersrepository.Repository
}

type lockableMutex interface {
	Lock() error
	Unlock() (bool, error)
}

type kickTokenRefresher interface {
	RefreshToken(ctx context.Context, refreshToken string) (gokick.TokenResponse, error)
}

type tokensImpl struct {
	globalClient   *helix.Client
	httpClient     *http.Client
	appAccessToken *appToken
	kickAppToken   *appToken

	config           cfg.Config
	log              *slog.Logger
	redSync          *redsync.Redsync
	twirBus          *buscore.Bus
	kickBotsRepo     kickbotsrepository.Repository
	integrationsRepo integrationsrepository.Repository
	channelIntegrationsRepo channelsintegrationsrepository.Repository
	spotifyIntegrationsRepo channelsintegrationsspotifyrepository.Repository
	tokensRepository tokensrepository.Repository
	usersRepository  usersrepository.Repository
	newMutex         func(name string) lockableMutex
	newKickTokenRefresher func() (kickTokenRefresher, error)
	spotifyTokenURL  string
	nightbotTokenURL string
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
		kickBotsRepo:     opts.KickBotsRepo,
		integrationsRepo: opts.IntegrationsRepo,
		channelIntegrationsRepo: opts.ChannelIntegrationsRepo,
		spotifyIntegrationsRepo: opts.SpotifyIntegrationsRepo,
		tokensRepository: opts.TokensRepository,
		usersRepository:  opts.UsersRepository,
		newMutex: func(name string) lockableMutex {
			return opts.Redsync.NewMutex(name)
		},
		newKickTokenRefresher: func() (kickTokenRefresher, error) {
			return gokick.NewClient(&gokick.ClientOptions{
				HTTPClient:   httpClient,
				ClientID:     opts.Config.KickClientId,
				ClientSecret: opts.Config.KickClientSecret,
			})
		},
		spotifyTokenURL:  "https://accounts.spotify.com/api/token",
		nightbotTokenURL: "https://api.nightbot.tv/oauth2/token",
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
				if err := impl.twirBus.Tokens.RequestChannelIntegrationToken.SubscribeGroup(
					"tokens",
					impl.RequestChannelIntegrationToken,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.twirBus.Tokens.RequestAppToken.Unsubscribe()
				impl.twirBus.Tokens.RequestUserToken.Unsubscribe()
				impl.twirBus.Tokens.RequestBotToken.Unsubscribe()
				impl.twirBus.Tokens.RequestChannelIntegrationToken.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *tokensImpl) RequestAppToken(
	ctx context.Context,
	data tokens.GetAppTokenRequest,
) (tokens.TokenResponse, error) {
	platform := data.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	if platform == platformentity.PlatformKick {
		return c.requestKickAppToken(ctx)
	}

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
		ExpiresIn:   int32(c.appAccessToken.ExpiresIn),
	}, nil
}

func (c *tokensImpl) requestKickAppToken(ctx context.Context) (tokens.TokenResponse, error) {
	mu := c.redSync.NewMutex("tokens-kick-app-lock")
	mu.Lock()
	defer mu.Unlock()

	if c.kickAppToken == nil || isTokenExpired(c.kickAppToken.ExpiresIn, c.kickAppToken.ObtainmentTime) {
		client, err := gokick.NewClient(
			&gokick.ClientOptions{
				ClientID:     c.config.KickClientId,
				ClientSecret: c.config.KickClientSecret,
			},
		)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("create kick client: %w", err)
		}

		resp, err := client.GetAppAccessToken(ctx)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("get kick app token: %w", err)
		}

		c.kickAppToken = &appToken{
			AccessToken:    resp.AccessToken,
			ObtainmentTime: time.Now().UTC(),
			ExpiresIn:      resp.ExpiresIn,
		}
		c.log.Info("kick app token refreshed")
	}

	return tokens.TokenResponse{
		AccessToken: c.kickAppToken.AccessToken,
		Scopes:      []string{},
		ExpiresIn:   int32(c.kickAppToken.ExpiresIn),
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
		user, err := c.usersRepository.GetByID(ctx, userID)
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
			if kickTokens.RefreshToken == "" {
				kickTokens.RefreshToken = decryptedRefreshToken
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
			slog.String("platform", string(user.Platform)),
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
	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			ClientID:     c.config.KickClientId,
			ClientSecret: c.config.KickClientSecret,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create kick client: %w", err)
	}

	res, err := client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh kick token: %w", err)
	}

	return &kickTokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
		Scopes:       strings.Fields(res.Scope),
	}, nil
}

func (c *tokensImpl) RequestBotToken(
	ctx context.Context,
	data tokens.GetBotTokenRequest,
) (tokens.TokenResponse, error) {
	platform := data.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	if platform == platformentity.PlatformKick {
		return c.requestKickBotToken(ctx)
	}

	mu := c.newMutex("tokens-bots-lock-" + data.BotId)
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

func (c *tokensImpl) requestKickBotToken(ctx context.Context) (tokens.TokenResponse, error) {
	mu := c.newMutex("tokens-kick-bot-lock")
	mu.Lock()
	defer mu.Unlock()

	bot, err := c.kickBotsRepo.GetDefault(ctx)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("cannot get default kick bot from repository: %w", err)
	}

	decryptedRefreshToken, err := crypto.Decrypt(bot.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("decrypt kick bot refresh token: %w", err)
	}

	if isTokenExpired(bot.ExpiresIn, bot.ObtainmentTimestamp) {
		client, err := c.newKickTokenRefresher()
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("create kick client: %w", err)
		}

		resp, err := client.RefreshToken(ctx, decryptedRefreshToken)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("refresh kick bot token: %w", err)
		}

		refreshToken := resp.RefreshToken
		if refreshToken == "" {
			refreshToken = decryptedRefreshToken
		}

		encryptedAccessToken, err := crypto.Encrypt(resp.AccessToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("encrypt kick bot access token: %w", err)
		}

		encryptedRefreshToken, err := crypto.Encrypt(refreshToken, c.config.TokensCipherKey)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("encrypt kick bot refresh token: %w", err)
		}

		updatedBot, err := c.kickBotsRepo.UpdateToken(
			ctx,
			bot.ID,
			kickbotsrepository.UpdateTokenInput{
				AccessToken:         encryptedAccessToken,
				RefreshToken:        encryptedRefreshToken,
				Scopes:              strings.Fields(resp.Scope),
				ExpiresIn:           resp.ExpiresIn,
				ObtainmentTimestamp: time.Now().UTC(),
			},
		)
		if err != nil {
			return tokens.TokenResponse{}, fmt.Errorf("persist kick bot token: %w", err)
		}

		bot = updatedBot
		c.log.Info("kick bot token refreshed", slog.String("kick_bot_id", bot.ID.String()))
	}

	decryptedAccessToken, err := crypto.Decrypt(bot.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("decrypt kick bot access token: %w", err)
	}

	return tokens.TokenResponse{
		AccessToken: decryptedAccessToken,
		Scopes:      bot.Scopes,
		ExpiresIn:   int32(bot.ExpiresIn),
	}, nil
}

func (c *tokensImpl) RequestChannelIntegrationToken(
	ctx context.Context,
	data tokens.GetChannelIntegrationTokenRequest,
) (tokens.TokenResponse, error) {
	switch data.Service {
	case integrationsmodel.ServiceSpotify:
		return c.requestSpotifyChannelIntegrationToken(ctx, data.ChannelID)
	case integrationsmodel.ServiceNightbot:
		return c.requestNightbotChannelIntegrationToken(ctx, data.ChannelID)
	default:
		return tokens.TokenResponse{}, fmt.Errorf("unsupported integration service: %s", data.Service)
	}
}

func (c *tokensImpl) requestSpotifyChannelIntegrationToken(ctx context.Context, channelID string) (tokens.TokenResponse, error) {
	mu := c.newMutex("tokens-channel-integration-spotify-" + channelID)
	mu.Lock()
	defer mu.Unlock()

	integration, err := c.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceSpotify)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("get spotify integration settings: %w", err)
	}
	if integration.ClientID == nil || integration.ClientSecret == nil {
		return tokens.TokenResponse{}, fmt.Errorf("spotify integration missing client credentials")
	}

	channelIntegration, err := c.spotifyIntegrationsRepo.GetByChannelID(ctx, channelID)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("get spotify channel integration: %w", err)
	}
	if channelIntegration.RefreshToken == "" {
		return tokens.TokenResponse{}, fmt.Errorf("spotify channel integration missing refresh token")
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", channelIntegration.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.spotifyTokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("create spotify refresh request: %w", err)
	}
	req.SetBasicAuth(*integration.ClientID, *integration.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("refresh spotify token: %w", err)
	}
	defer resp.Body.Close()

	var refreshData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := decodeJsonResponse(resp, &refreshData); err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("decode spotify refresh response: %w", err)
	}

	update := channelsintegrationsspotifyrepository.UpdateInput{
		AccessToken: &refreshData.AccessToken,
	}
	if refreshData.RefreshToken != "" {
		update.RefreshToken = &refreshData.RefreshToken
	}
	if err := c.spotifyIntegrationsRepo.Update(ctx, channelIntegration.ID, update); err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("persist spotify token: %w", err)
	}

	return tokens.TokenResponse{
		AccessToken: refreshData.AccessToken,
		Scopes:      channelIntegration.Scopes,
	}, nil
}

func (c *tokensImpl) requestNightbotChannelIntegrationToken(ctx context.Context, channelID string) (tokens.TokenResponse, error) {
	mu := c.newMutex("tokens-channel-integration-nightbot-" + channelID)
	mu.Lock()
	defer mu.Unlock()

	integration, err := c.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("get nightbot integration settings: %w", err)
	}
	if integration.ClientID == nil || integration.ClientSecret == nil {
		return tokens.TokenResponse{}, fmt.Errorf("nightbot integration missing client credentials")
	}

	channelIntegration, err := c.channelIntegrationsRepo.GetByChannelAndService(ctx, channelID, integrationsmodel.ServiceNightbot)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("get nightbot channel integration: %w", err)
	}
	if channelIntegration.RefreshToken == nil || *channelIntegration.RefreshToken == "" {
		return tokens.TokenResponse{}, fmt.Errorf("nightbot channel integration missing refresh token")
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", *integration.ClientID)
	formData.Set("client_secret", *integration.ClientSecret)
	formData.Set("refresh_token", *channelIntegration.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.nightbotTokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("create nightbot refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("refresh nightbot token: %w", err)
	}
	defer resp.Body.Close()

	var refreshData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := decodeJsonResponse(resp, &refreshData); err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("decode nightbot refresh response: %w", err)
	}

	updateInput := channelsintegrationsrepository.UpdateInput{
		Enabled:     boolPtr(true),
		AccessToken: &refreshData.AccessToken,
	}
	if refreshData.RefreshToken != "" {
		updateInput.RefreshToken = &refreshData.RefreshToken
	}

	if err := c.channelIntegrationsRepo.Update(ctx, channelIntegration.ID, updateInput); err != nil {
		return tokens.TokenResponse{}, fmt.Errorf("persist nightbot token: %w", err)
	}

	return tokens.TokenResponse{
		AccessToken: refreshData.AccessToken,
		ExpiresIn:   int32(refreshData.ExpiresIn),
	}, nil
}
