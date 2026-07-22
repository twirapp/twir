package vkvideo

import (
	"context"
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

func TestProviderBuildsVKIDAuthorizationURL(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("authorization URL must not make an HTTP request")
		return nil, nil
	}))

	url := provider.GetAuthURL("state-value", "pkce-challenge")
	if !strings.Contains(url, "response_type=code") || !strings.Contains(url, "scope=vkid.personal_info") {
		t.Fatalf("unexpected VK ID authorization URL: %s", url)
	}
}

func TestProviderExchangeCodeMapsTokensAndDeviceID(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/oauth2/auth" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("parse request form: %v", err)
		}
		if got := req.PostForm.Get("device_id"); got != "device-id" {
			t.Errorf("device_id = %q, want device-id", got)
		}

		return providerResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"scope":"vkid.personal_info"}`), nil
	}))

	tokens, err := provider.ExchangeCode(context.Background(), appplatform.ExchangeCodeInput{
		Code:         "authorization-code",
		CodeVerifier: "pkce-verifier",
		DeviceID:     "device-id",
	})
	if err != nil {
		t.Fatalf("exchange VK ID code: %v", err)
	}
	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" || tokens.ExpiresIn != 3600 || tokens.DeviceID != "device-id" {
		t.Fatalf("unexpected platform tokens: %#v", tokens)
	}
}

func TestProviderGetUserMapsVKIDProfileWithoutLogin(t *testing.T) {
	provider := newTestProvider(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/oauth2/user_info" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		return providerResponse(http.StatusOK, `{"user":{"user_id":"123","first_name":"Ada","last_name":"Lovelace","avatar":"https://cdn.example.test/avatar.jpg"}}`), nil
	}))

	user, err := provider.GetUser(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get VK ID user: %v", err)
	}
	if user.ID != "123" || user.Login != "" || user.DisplayName != "Ada Lovelace" || user.Avatar != "https://cdn.example.test/avatar.jpg" {
		t.Fatalf("unexpected platform user: %#v", user)
	}
}

func newTestProvider(t *testing.T, transport http.RoundTripper) *Provider {
	t.Helper()

	client, err := vk.NewIDClient(vk.IDClientOpts{
		ClientID:     "client-id",
		ServiceToken: "service-token",
		RedirectURL:  "https://twir.example.test/auth/vk/callback",
		APIBaseURL:   "https://id.example.test",
		HTTPClient:   &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("create VK ID client: %v", err)
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
