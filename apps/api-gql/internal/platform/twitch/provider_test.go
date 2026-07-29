package twitchplatform

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/twirapp/twir/apps/api-gql/internal/platform"
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

func TestExchangeCodeUsesTwitchyOAuthClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth2/token" {
			t.Fatalf("request path = %q, want /oauth2/token", request.URL.Path)
		}
		if err := request.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if request.Form.Get("grant_type") != "authorization_code" || request.Form.Get("code") != "authorization-code" || request.Form.Get("client_id") != "client-id" || request.Form.Get("client_secret") != "client-secret" || request.Form.Get("redirect_uri") != "https://twir.example.test/login" {
			t.Fatalf("token request form = %#v", request.Form)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"scope":["channel:bot","channel:manage:raids"],"token_type":"bearer"}`))
	}))
	defer server.Close()

	provider := New(Opts{Config: cfg.Config{
		SiteBaseUrl:        "https://twir.example.test",
		TwitchClientId:     "client-id",
		TwitchClientSecret: "client-secret",
		TwitchMockEnabled:  true,
		TwitchMockAuthUrl:  server.URL,
	}})

	tokens, err := provider.ExchangeCode(context.Background(), platform.ExchangeCodeInput{Code: "authorization-code"})
	if err != nil {
		t.Fatalf("exchange code: %v", err)
	}
	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" || tokens.ExpiresIn != 3600 || !reflect.DeepEqual(tokens.Scopes, []string{"channel:bot", "channel:manage:raids"}) {
		t.Fatalf("exchange tokens = %#v", tokens)
	}
}

func TestGetUserUsesMockHelixEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/helix/users" {
			t.Fatalf("request path = %q, want /helix/users", request.URL.Path)
		}
		if request.Header.Get("Authorization") != "Bearer access-token" || request.Header.Get("Client-Id") != "client-id" {
			t.Fatalf("request headers = %#v", request.Header)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"data":[{"id":"123","login":"streamer","display_name":"Streamer","profile_image_url":"https://cdn.example.test/avatar.png"}]}`))
	}))
	defer server.Close()

	provider := New(Opts{Config: cfg.Config{
		TwitchClientId:    "client-id",
		TwitchMockEnabled: true,
		TwitchMockApiUrl:  server.URL + "/helix",
	}})

	user, err := provider.GetUser(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if user.ID != "123" || user.Login != "streamer" || user.DisplayName != "Streamer" || user.Avatar != "https://cdn.example.test/avatar.png" {
		t.Fatalf("platform user = %#v", user)
	}
}

func TestRefreshTokenWrapsTwitchyOAuthErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer server.Close()

	provider := New(Opts{Config: cfg.Config{
		TwitchClientId:     "client-id",
		TwitchClientSecret: "client-secret",
		TwitchMockEnabled:  true,
		TwitchMockAuthUrl:  server.URL,
	}})

	_, err := provider.RefreshToken(context.Background(), platform.RefreshTokenInput{RefreshToken: "refresh-token"})
	if err == nil || !strings.Contains(err.Error(), "refresh Twitch user access token") {
		t.Fatalf("refresh error = %v", err)
	}
}
