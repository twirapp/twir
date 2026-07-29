package twitch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kvizyx/twitchy/helix"
	twitchoauth "github.com/kvizyx/twitchy/oauth"
	"github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	sharedoauth "github.com/twirapp/twir/libs/oauth"
	"github.com/twirapp/twir/libs/repositories/tokens"
	tokenspgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"
)

const (
	twitchProvider      sharedoauth.Provider = "twitch"
	twitchRefreshPrefix                      = "twir:twitch:refresh"
)

var (
	ErrCredentialNotFound      = errors.New("Twitch credential not found")
	ErrChannelBotConflict      = errors.New("Twitch channel already has a registered bot")
	ErrChannelBotNotRegistered = errors.New("Twitch channel bot is not registered")
	ErrRegistrationConflict    = errors.New("Twitch credential registration conflicts with an existing identity")
)

type credentialKind uint8

const (
	broadcasterCredential credentialKind = iota + 1
	botCredential
)

type registeredCredential struct {
	kind    credentialKind
	intents map[helix.Intent]struct{}
}

type twitchRuntime struct {
	rootClient *helix.Client
	appClient  *helix.Client
	registry   *twitchoauth.CoordinatedRegistry
	tokens     tokens.Repository
	cipherKey  string
	now        func() time.Time

	appSource   *sharedoauth.AppTokenSource
	redisClient *redis.Client
	pool        *pgxpool.Pool

	registrationMu sync.Mutex
	registrations  sync.Map
	channelBots    map[string]string
}

var packageRuntime struct {
	mu    sync.Mutex
	once  sync.Once
	value *twitchRuntime
	err   error
}

func getRuntime(ctx context.Context, config cfg.Config) (*twitchRuntime, error) {
	packageRuntime.mu.Lock()
	defer packageRuntime.mu.Unlock()
	packageRuntime.once.Do(func() {
		packageRuntime.value, packageRuntime.err = newTwitchRuntime(ctx, config)
	})
	return packageRuntime.value, packageRuntime.err
}

func closeRuntime() error {
	packageRuntime.mu.Lock()
	defer packageRuntime.mu.Unlock()
	if packageRuntime.value == nil {
		return packageRuntime.err
	}
	err := packageRuntime.value.close()
	packageRuntime.value = nil
	packageRuntime.err = nil
	packageRuntime.once = sync.Once{}
	return err
}

func newTwitchRuntime(ctx context.Context, config cfg.Config) (*twitchRuntime, error) {
	pool, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	redisOptions, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("parse Redis URL: %w", err)
	}
	redisClient := redis.NewClient(redisOptions)
	runtime, err := newTwitchRuntimeWithDependencies(ctx, config, tokenspgx.NewFx(pool), redisClient, pool, time.Now)
	if err != nil {
		pool.Close()
		_ = redisClient.Close()
		return nil, err
	}
	return runtime, nil
}

func newTwitchRuntimeWithDependencies(
	ctx context.Context,
	config cfg.Config,
	tokenRepository tokens.Repository,
	redisClient *redis.Client,
	pool *pgxpool.Pool,
	now func() time.Time,
) (*twitchRuntime, error) {
	if tokenRepository == nil || redisClient == nil || now == nil {
		return nil, fmt.Errorf("initialize Twitch runtime: missing dependency")
	}

	httpClient := twitchHTTPClient(config)
	oauthClient, err := newTwitchOAuthClient(config, httpClient)
	if err != nil {
		return nil, fmt.Errorf("create Twitch OAuth client: %w", err)
	}
	coordinator, err := twitchoauth.NewRedisRefreshCoordinator(redisClient, func(userID string) string {
		return twitchRefreshPrefix + ":" + userID
	})
	if err != nil {
		return nil, fmt.Errorf("create Twitch refresh coordinator: %w", err)
	}
	registry, err := twitchoauth.NewCoordinatedRegistry(oauthClient, coordinator)
	if err != nil {
		return nil, fmt.Errorf("create Twitch credential registry: %w", err)
	}
	locker, err := sharedoauth.NewRedisLocker(redisClient, sharedoauth.RedisLockerOptions{
		Prefix: twitchRefreshPrefix, TTL: 30 * time.Second, RenewEvery: 10 * time.Second, Timeout: time.Second,
	})
	if err != nil {
		_ = registry.Close()
		return nil, fmt.Errorf("create Twitch app-token locker: %w", err)
	}
	appSource, err := sharedoauth.NewAppTokenSource(sharedoauth.AppTokenDependencies{
		Store:   newAppTokenStore(),
		Fetcher: twitchAppTokenFetcher{client: oauthClient, clientID: config.TwitchClientId, clientSecret: config.TwitchClientSecret, now: now},
		Locker:  locker,
	}, sharedoauth.AppTokenSourceOptions{Skew: time.Minute})
	if err != nil {
		_ = registry.Close()
		return nil, fmt.Errorf("create Twitch app-token source: %w", err)
	}
	rootClient, err := newRootHelixClient(config, httpClient, registry)
	if err != nil {
		_ = appSource.Close()
		_ = registry.Close()
		return nil, fmt.Errorf("create Twitch credential client: %w", err)
	}
	appClient, err := newAppHelixClient(config, httpClient, appSource, now)
	if err != nil {
		_ = appSource.Close()
		_ = registry.Close()
		return nil, fmt.Errorf("create Twitch app client: %w", err)
	}

	return &twitchRuntime{
		rootClient: rootClient, appClient: appClient, registry: registry, tokens: tokenRepository,
		cipherKey: config.TokensCipherKey, now: now, appSource: appSource, redisClient: redisClient,
		pool: pool, channelBots: make(map[string]string),
	}, nil
}

func (runtime *twitchRuntime) close() error {
	registryErr := runtime.registry.Close()
	appSourceErr := runtime.appSource.Close()
	redisErr := runtime.redisClient.Close()
	if runtime.pool != nil {
		runtime.pool.Close()
	}
	return errors.Join(registryErr, appSourceErr, redisErr)
}

func twitchHTTPClient(config cfg.Config) *http.Client {
	httpClient := createHttpClient()
	if config.TwitchMockEnabled {
		httpClient = &http.Client{Transport: NewMockRoundTripper(httpClient.Transport, config)}
	}
	return httpClient
}

func newTwitchOAuthClient(config cfg.Config, httpClient *http.Client) (*twitchoauth.Client, error) {
	options := []twitchoauth.Option{twitchoauth.WithHTTPClient(httpClient)}
	if config.TwitchMockEnabled {
		options = append(options, twitchoauth.WithBaseURL(config.TwitchMockAuthUrl))
	}
	return twitchoauth.New(options...)
}

func newRootHelixClient(config cfg.Config, httpClient *http.Client, registry *twitchoauth.CoordinatedRegistry) (*helix.Client, error) {
	options := append(commonHelixOptions(config, httpClient), helix.WithCredentialResolver(registry))
	return helix.New(options...)
}

func newAppHelixClient(config cfg.Config, httpClient *http.Client, source *sharedoauth.AppTokenSource, now func() time.Time) (*helix.Client, error) {
	adapter := appTokenSourceAdapter{
		source:   source,
		key:      sharedoauth.AppTokenKey{Provider: twitchProvider, ID: sharedoauth.CredentialID(config.TwitchClientId)},
		clientID: config.TwitchClientId,
		now:      now,
	}
	options := append(commonHelixOptions(config, httpClient), helix.WithTokenSource(adapter))
	return helix.New(options...)
}

func commonHelixOptions(config cfg.Config, httpClient *http.Client) []helix.Option {
	options := []helix.Option{
		helix.WithHTTPClient(httpClient),
		helix.WithRateLimitPolicy(helix.RateLimitPolicy{Wait: true, MaxWait: 10 * time.Minute}),
	}
	if config.TwitchMockEnabled {
		options = append(options, helix.WithBaseURL(config.TwitchMockApiUrl))
	}
	return options
}
