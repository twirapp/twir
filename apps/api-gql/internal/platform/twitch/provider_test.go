package twitchplatform

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
)

func TestGetAuthURLUsesLegacyProductionScopes(t *testing.T) {
	provider := New(Opts{Config: cfg.Config{
		SiteBaseUrl:    "https://twir.example.test",
		TwitchClientId: "client-id",
	}})

	authorizeURL, err := url.Parse(provider.GetAuthURL("state-value", ""))
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}

	wantScopes := []string{
		"moderation:read",
		"channel:manage:broadcast",
		"channel:read:redemptions",
		"channel:manage:redemptions",
		"moderator:read:chatters",
		"moderator:manage:shoutouts",
		"moderator:manage:banned_users",
		"channel:read:vips",
		"channel:manage:vips",
		"channel:manage:moderators",
		"moderator:read:followers",
		"moderator:manage:chat_settings",
		"channel:read:polls",
		"channel:manage:polls",
		"channel:read:predictions",
		"channel:manage:predictions",
		"channel:read:subscriptions",
		"channel:moderate",
		"user:read:follows",
		"channel:bot",
		"channel:manage:raids",
	}
	if got := strings.Fields(authorizeURL.Query().Get("scope")); !reflect.DeepEqual(got, wantScopes) {
		t.Fatalf("production scopes = %#v, want %#v", got, wantScopes)
	}
}

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
	if scopes, ok := authorizeURL.Query()["scope"]; !ok || len(scopes) != 1 || scopes[0] != "" {
		t.Fatalf("mock scopes = %#v, want [\"\"]", scopes)
	}
}
