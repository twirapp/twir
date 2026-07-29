package twitch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kvizyx/twitchy/helix"
	twitchoauth "github.com/kvizyx/twitchy/oauth"
	sharedoauth "github.com/twirapp/twir/libs/oauth"
)

type appTokenStore struct {
	mu     sync.RWMutex
	tokens map[sharedoauth.AppTokenKey]sharedoauth.AppToken
}

func newAppTokenStore() *appTokenStore {
	return &appTokenStore{tokens: make(map[sharedoauth.AppTokenKey]sharedoauth.AppToken)}
}

func (store *appTokenStore) LoadAppToken(_ context.Context, key sharedoauth.AppTokenKey) (sharedoauth.AppToken, error) {
	store.mu.RLock()
	token, found := store.tokens[key]
	store.mu.RUnlock()
	if !found {
		return sharedoauth.AppToken{}, sharedoauth.ErrAppTokenNotFound
	}
	return token, nil
}

func (store *appTokenStore) CommitAppToken(_ context.Context, key sharedoauth.AppTokenKey, token sharedoauth.AppToken) error {
	store.mu.Lock()
	store.tokens[key] = token
	store.mu.Unlock()
	return nil
}

type twitchAppTokenFetcher struct {
	client       *twitchoauth.Client
	clientID     string
	clientSecret string
	now          func() time.Time
}

func (fetcher twitchAppTokenFetcher) FetchAppToken(ctx context.Context, _ sharedoauth.AppTokenKey) (sharedoauth.AppToken, error) {
	pair, err := fetcher.client.ClientCredentials(ctx, twitchoauth.ClientCredentialsRequest{
		ClientID: fetcher.clientID, ClientSecret: fetcher.clientSecret,
	})
	if err != nil {
		return sharedoauth.AppToken{}, fmt.Errorf("request Twitch app credential: %w", err)
	}
	return sharedoauth.AppToken{AccessToken: pair.AccessToken, ObtainedAt: fetcher.now(), ExpiresIn: pair.ExpiresIn}, nil
}

type appTokenSourceAdapter struct {
	source   *sharedoauth.AppTokenSource
	key      sharedoauth.AppTokenKey
	clientID string
	now      func() time.Time
}

func (adapter appTokenSourceAdapter) Token(ctx context.Context) (helix.CredentialSnapshot, error) {
	token, err := adapter.source.Token(ctx, adapter.key)
	if err != nil {
		return helix.CredentialSnapshot{}, fmt.Errorf("load Twitch app credential: %w", err)
	}
	return helix.NewCredentialSnapshot(helix.Credential{
		AccessToken: token.AccessToken, ClientID: adapter.clientID, TokenClass: helix.TokenClassApp,
		ExpiresAt: token.ObtainedAt.Add(token.ExpiresIn),
	}), nil
}
