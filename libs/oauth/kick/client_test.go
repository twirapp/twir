package kick

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/oauth"
)

func TestAppFetcherFetchesKickAppToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth/token" {
			t.Fatalf("path = %s", request.URL.Path)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"access_token":"app-access","expires_in":3600}`))
	}))
	defer server.Close()

	fetcher := appFetcher{factory: func() (AppClient, error) {
		return gokick.NewClient(&gokick.ClientOptions{ClientID: "client", ClientSecret: "secret", AuthBaseURL: server.URL})
	}}

	token, err := fetcher.FetchAppToken(context.Background(), oauth.AppTokenKey{Provider: Provider, ID: appCredentialID})
	if err != nil {
		t.Fatalf("fetch app token: %v", err)
	}
	if token.AccessToken != "app-access" || token.ExpiresIn != time.Hour {
		t.Fatalf("token = %#v", token)
	}
}

func TestRefresherPreservesOmittedRefreshTokenAndParsesScopes(t *testing.T) {
	refresher := newRefresher(func() (Client, error) {
		return fakeClient{response: gokick.TokenResponse{AccessToken: "rotated-access", ExpiresIn: 3600, Scope: "channel:read  channel:write"}}, nil
	})

	result, err := refresher.Refresh(context.Background(), oauth.Credential{RefreshToken: "original-refresh"})
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if result.RefreshToken != nil {
		t.Fatalf("refresh token = %q, want omitted", *result.RefreshToken)
	}
	if len(result.Scopes) != 2 || result.Scopes[0] != "channel:read" || result.Scopes[1] != "channel:write" {
		t.Fatalf("scopes = %#v", result.Scopes)
	}
}

func TestRefresherReturnsRedactedProviderError(t *testing.T) {
	refresher := newRefresher(func() (Client, error) {
		return fakeClient{err: providerStatusError{status: http.StatusUnauthorized}}, nil
	})

	_, err := refresher.Refresh(context.Background(), oauth.Credential{RefreshToken: "redacted"})
	var providerError ProviderError
	if !errors.As(err, &providerError) || providerError.StatusCode != http.StatusUnauthorized {
		t.Fatalf("error = %v", err)
	}
}

type fakeClient struct {
	response gokick.TokenResponse
	err      error
}

func (client fakeClient) RefreshToken(context.Context, string) (gokick.TokenResponse, error) {
	return client.response, client.err
}

type providerStatusError struct{ status int }

func (error providerStatusError) Error() string { return "provider response" }
func (error providerStatusError) Code() int     { return error.status }
