package nightbot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/oauth"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	integrations "github.com/twirapp/twir/libs/repositories/integrations"
)

const (
	Provider        = oauth.Provider("nightbot")
	defaultTokenURL = "https://api.nightbot.tv/oauth2/token"
	lockPrefix      = "oauth"
)

type TokenSource interface {
	Token(context.Context, string) (oauth.Credential, error)
}

type SourceOptions struct {
	Redis      *redis.Client
	HTTPClient *http.Client
	TokenURL   string
}

func (options SourceOptions) tokenURL() string {
	if options.TokenURL == "" {
		return defaultTokenURL
	}
	return options.TokenURL
}

func (options SourceOptions) httpClient() *http.Client {
	if options.HTTPClient == nil {
		return http.DefaultClient
	}
	return options.HTTPClient
}

func (options SourceOptions) locker() (*oauth.RedisLocker, error) {
	return oauth.NewRedisLocker(options.Redis, oauth.RedisLockerOptions{
		Prefix: lockPrefix, TTL: 30 * time.Second, RenewEvery: 10 * time.Second, Timeout: 5 * time.Second,
	})
}

func NewTokenSource(
	options SourceOptions,
	channelIntegrations channelsintegrations.Repository,
	integrationSettings integrations.Repository,
) (TokenSource, error) {
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create Nightbot token locker: %w", err)
	}
	runtime, err := oauth.NewRefreshRuntime(
		store{repository: channelIntegrations},
		refresher{client: options.httpClient(), tokenURL: options.tokenURL(), integrations: integrationSettings},
		locker,
		oauth.RuntimeOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("create Nightbot token runtime: %w", err)
	}
	return tokenSource{runtime: runtime}, nil
}

type tokenSource struct {
	runtime *oauth.RefreshRuntime
}

func (source tokenSource) Token(ctx context.Context, channelID string) (oauth.Credential, error) {
	return source.runtime.Refresh(ctx, oauth.CredentialKey{Provider: Provider, ID: oauth.CredentialID(channelID)})
}
