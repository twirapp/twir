package vkvideo

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestProviderBuildsDevAPIAuthorizationURL(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("authorization URL must not make an HTTP request")
		return nil, nil
	}))

	rawURL := provider.GetAuthURL("state-value", "ignored-pkce-challenge")
	if !strings.HasPrefix(rawURL, "https://auth.example.test/app/oauth2/authorize?") {
		t.Fatalf("unexpected VK Video authorization URL: %s", rawURL)
	}
	if !strings.Contains(rawURL, "response_type=code") {
		t.Fatalf("unexpected VK Video authorization URL: %s", rawURL)
	}
	if strings.Contains(rawURL, "scope=") {
		t.Fatalf("VK Video DevAPI current user requires no scope, got: %s", rawURL)
	}
	if strings.Contains(rawURL, "code_challenge") {
		t.Fatalf("VK Video DevAPI authorization must not use PKCE: %s", rawURL)
	}
}

func TestProviderExchangeCodeUsesBasicAuthAndMapsTokens(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/oauth/server/token" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("client-id:client-secret"))
		if auth := req.Header.Get("Authorization"); auth != wantAuth {
			t.Errorf("authorization header = %q, want %q", auth, wantAuth)
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("parse request form: %v", err)
		}
		if got := req.PostForm.Get("code"); got != "authorization-code" {
			t.Errorf("code = %q, want authorization-code", got)
		}
		if got := req.PostForm.Get("grant_type"); got != "authorization_code" {
			t.Errorf("grant_type = %q, want authorization_code", got)
		}

		return providerResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"token_type":"Bearer"}`), nil
	}))

	tokens, err := provider.ExchangeCode(context.Background(), appplatform.ExchangeCodeInput{
		Code: "authorization-code",
	})
	if err != nil {
		t.Fatalf("exchange VK Video code: %v", err)
	}
	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" || tokens.ExpiresIn != 3600 {
		t.Fatalf("unexpected platform tokens: %#v", tokens)
	}
	if tokens.DeviceID != "" {
		t.Fatalf("VK Video DevAPI tokens must not carry a device ID: %#v", tokens)
	}
}

func TestProviderRefreshTokenSendsRefreshGrant(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/oauth/server/token" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("parse request form: %v", err)
		}
		if got := req.PostForm.Get("grant_type"); got != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", got)
		}
		if got := req.PostForm.Get("refresh_token"); got != "refresh-token" {
			t.Errorf("refresh_token = %q, want refresh-token", got)
		}

		return providerResponse(http.StatusOK, `{"access_token":"new-access-token","refresh_token":"new-refresh-token","expires_in":7200,"token_type":"Bearer"}`), nil
	}))

	tokens, err := provider.RefreshToken(context.Background(), appplatform.RefreshTokenInput{
		RefreshToken: "refresh-token",
	})
	if err != nil {
		t.Fatalf("refresh VK Video token: %v", err)
	}
	if tokens.AccessToken != "new-access-token" || tokens.RefreshToken != "new-refresh-token" || tokens.ExpiresIn != 7200 {
		t.Fatalf("unexpected platform tokens: %#v", tokens)
	}
}

func TestProviderGetUserMapsDevAPIProfile(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/current_user" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
			t.Errorf("authorization header = %q, want Bearer access-token", auth)
		}
		return providerResponse(http.StatusOK, `{"data":{"user":{"id":123,"nick":"ada","avatar_url":"https://cdn.example.test/avatar.jpg"}}}`), nil
	}))

	user, err := provider.GetUser(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get VK Video user: %v", err)
	}
	if user.ID != "123" || user.Login != "ada" || user.DisplayName != "ada" || user.Avatar != "https://cdn.example.test/avatar.jpg" {
		t.Fatalf("unexpected platform user: %#v", user)
	}
}

func newTestProvider(t *testing.T, transport http.RoundTripper) *Provider {
	t.Helper()

	client, err := vk.NewOAuthClient(vk.OAuthClientOpts{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://twir.example.test/auth/vk/callback",
		APIBaseURL:   "https://api.example.test",
		AuthBaseURL:  "https://auth.example.test",
		HTTPClient:   &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("create VK Video OAuth client: %v", err)
	}

	return &Provider{client: client}
}

func providerResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
