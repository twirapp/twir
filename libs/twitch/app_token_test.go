package twitch

import (
	"context"
	"testing"
	"time"

	sharedoauth "github.com/twirapp/twir/libs/oauth"
)

type appTokenFetcherFake struct{ token sharedoauth.AppToken }

func (fetcher appTokenFetcherFake) FetchAppToken(context.Context, sharedoauth.AppTokenKey) (sharedoauth.AppToken, error) {
	return fetcher.token, nil
}

type appTokenLockerFake struct{ lease appTokenLeaseFake }

func (locker appTokenLockerFake) AcquireAppToken(context.Context, sharedoauth.AppTokenKey) (sharedoauth.Lease, error) {
	return locker.lease, nil
}

type appTokenLeaseFake struct{ ctx context.Context }

func (lease appTokenLeaseFake) Context() context.Context      { return lease.ctx }
func (lease appTokenLeaseFake) Lost() <-chan struct{}         { return make(chan struct{}) }
func (lease appTokenLeaseFake) Release(context.Context) error { return nil }

func TestAppTokenSourceAdapter_returns_app_credential_without_user_registration(t *testing.T) {
	// Given
	now := time.Date(2026, time.July, 29, 12, 0, 0, 0, time.UTC)
	source, err := sharedoauth.NewAppTokenSource(sharedoauth.AppTokenDependencies{
		Store: newAppTokenStore(), Fetcher: appTokenFetcherFake{token: sharedoauth.AppToken{
			AccessToken: "app-token", ObtainedAt: now, ExpiresIn: time.Hour,
		}}, Locker: appTokenLockerFake{lease: appTokenLeaseFake{ctx: context.Background()}},
	}, sharedoauth.AppTokenSourceOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = source.Close() })
	adapter := appTokenSourceAdapter{source: source, key: sharedoauth.AppTokenKey{Provider: twitchProvider, ID: "client"}, clientID: "client", now: func() time.Time { return now }}

	// When
	credential, err := adapter.Token(context.Background())

	// Then
	if err != nil || credential.TokenClass() != "app" || credential.UserID() != "" {
		t.Fatalf("app source did not remain separate from user credentials: %v", err)
	}
}

func TestChannelBotIntent_is_exact_and_never_uses_chat_fallback(t *testing.T) {
	// Given
	channelID := "channel-123"

	// When
	intent := channelBotIntent(channelID)

	// Then
	if intent != "bot:channel:channel-123" {
		t.Fatalf("channel bot intent = %q", intent)
	}
}
