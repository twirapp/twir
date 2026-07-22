package twitchplatform

import (
	"net/url"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
)

func TestGetAuthURLUsesMockAuthorizationEndpoint(t *testing.T) {
	provider := New(Opts{Config: cfg.Config{
		SiteBaseUrl:        "https://twir.example.test",
		TwitchClientId:     "client-id",
		TwitchMockEnabled:  true,
		TwitchMockAuthUrl:  "https://twitch-mock.example.test",
		TwitchMockApiUrl:   "https://twitch-mock.example.test/helix",
		TwitchClientSecret: "client-secret",
	}})

	authorizeURL, err := url.Parse(provider.GetAuthURL("state-value", ""))
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if got, want := authorizeURL.Scheme+"://"+authorizeURL.Host+authorizeURL.Path, "https://twitch-mock.example.test/oauth2/authorize"; got != want {
		t.Fatalf("authorization endpoint = %q, want %q", got, want)
	}
	if got, want := authorizeURL.Query().Get("state"), "state-value"; got != want {
		t.Fatalf("state = %q, want %q", got, want)
	}
	if got, want := authorizeURL.Query().Get("redirect_uri"), "https://twir.example.test/login"; got != want {
		t.Fatalf("redirect URI = %q, want %q", got, want)
	}
}
