package kick

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/oauth"
	kickbots "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
)

const (
	appCredentialID = oauth.CredentialID("application")
	lockPrefix      = "oauth"
)

type AppTokenSource interface {
	Token(context.Context) (oauth.AppToken, error)
}

type UserTokenSource interface {
	Token(context.Context, uuid.UUID) (oauth.Credential, error)
}

type DefaultBotTokenSource interface {
	Token(context.Context) (oauth.Credential, error)
}

type SourceOptions struct {
	ClientID      string
	ClientSecret  string
	Redis         *goredis.Client
	CipherKey     string
	ClientFactory ClientFactory
	AppClientFactory AppClientFactory
}

func (options SourceOptions) appClientFactory() AppClientFactory {
	if options.AppClientFactory != nil {
		return options.AppClientFactory
	}
	return func() (AppClient, error) {
		return gokick.NewClient(&gokick.ClientOptions{ClientID: options.ClientID, ClientSecret: options.ClientSecret})
	}
}

func (options SourceOptions) clientFactory() ClientFactory {
	if options.ClientFactory != nil {
		return options.ClientFactory
	}
	return func() (Client, error) {
		return gokick.NewClient(&gokick.ClientOptions{ClientID: options.ClientID, ClientSecret: options.ClientSecret})
	}
}

func (options SourceOptions) locker() (*oauth.RedisLocker, error) {
	return oauth.NewRedisLocker(options.Redis, oauth.RedisLockerOptions{Prefix: lockPrefix, TTL: 30 * time.Second, RenewEvery: 10 * time.Second, Timeout: 5 * time.Second})
}

func NewAppTokenSource(options SourceOptions) (AppTokenSource, error) {
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create Kick app token locker: %w", err)
	}
	source, err := oauth.NewAppTokenSource(oauth.AppTokenDependencies{Store: &memoryAppStore{}, Fetcher: appFetcher{factory: options.appClientFactory()}, Locker: locker}, oauth.AppTokenSourceOptions{Skew: time.Minute})
	if err != nil {
		return nil, fmt.Errorf("create Kick app token source: %w", err)
	}
	return appSource{source: source}, nil
}

func NewUserTokenSource(options SourceOptions, repository tokens.Repository) (UserTokenSource, error) {
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create Kick user token locker: %w", err)
	}
	runtime, err := oauth.NewRefreshRuntime(userStore{repository: repository, cipherKey: options.CipherKey}, newRefresher(options.clientFactory()), locker, oauth.RuntimeOptions{Skew: time.Minute})
	if err != nil {
		return nil, fmt.Errorf("create Kick user token runtime: %w", err)
	}
	return userSource{runtime: runtime}, nil
}

func NewDefaultBotTokenSource(options SourceOptions, repository kickbots.Repository) (DefaultBotTokenSource, error) {
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create Kick bot token locker: %w", err)
	}
	runtime, err := oauth.NewRefreshRuntime(botStore{repository: repository, cipherKey: options.CipherKey}, newRefresher(options.clientFactory()), locker, oauth.RuntimeOptions{Skew: time.Minute})
	if err != nil {
		return nil, fmt.Errorf("create Kick bot token runtime: %w", err)
	}
	return botSource{runtime: runtime, repository: repository}, nil
}

type appSource struct{ source *oauth.AppTokenSource }

func (source appSource) Token(ctx context.Context) (oauth.AppToken, error) {
	return source.source.Token(ctx, oauth.AppTokenKey{Provider: Provider, ID: appCredentialID})
}

type userSource struct{ runtime *oauth.RefreshRuntime }

func (source userSource) Token(ctx context.Context, userID uuid.UUID) (oauth.Credential, error) {
	return source.runtime.Refresh(ctx, oauth.CredentialKey{Provider: Provider, ID: oauth.CredentialID(userID.String())})
}

type botSource struct {
	runtime    *oauth.RefreshRuntime
	repository kickbots.Repository
}

func (source botSource) Token(ctx context.Context) (oauth.Credential, error) {
	bot, err := source.repository.GetDefault(ctx)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load default Kick bot: %w", err)
	}
	return source.runtime.Refresh(ctx, oauth.CredentialKey{Provider: Provider, ID: oauth.CredentialID(bot.ID.String())})
}

type memoryAppStore struct {
	mu    sync.RWMutex
	token oauth.AppToken
	found bool
}

func (store *memoryAppStore) LoadAppToken(_ context.Context, _ oauth.AppTokenKey) (oauth.AppToken, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if !store.found {
		return oauth.AppToken{}, oauth.ErrAppTokenNotFound
	}
	return store.token, nil
}

func (store *memoryAppStore) CommitAppToken(_ context.Context, _ oauth.AppTokenKey, token oauth.AppToken) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.token = token
	store.found = true
	return nil
}
