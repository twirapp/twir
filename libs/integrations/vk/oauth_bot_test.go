package vk

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestOAuthClientExchangesBotCodeWithBotRedirectURL(t *testing.T) {
	// Given
	const botRedirectURL = "https://twir.example.test/api/auth/vk-video/bot-callback"
	client, err := NewOAuthClient(OAuthClientOpts{
		ClientID:      "client-id",
		ClientSecret:  "client-secret",
		RedirectURL:   botRedirectURL,
		APIBaseURL:    "https://api.example.test",
		AuthBaseURL:   "https://auth.example.test",
		DevAPIBaseURL: "https://devapi.example.test",
		HTTPClient: &http.Client{Transport: oauthRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			assertTokenEndpoint(t, req)
			assertExactFormValues(t, parseRequestForm(t, req), map[string]string{
				"grant_type":   "authorization_code",
				"code":         "bot-authorization-code",
				"redirect_uri": botRedirectURL,
			})
			return oauthResponse(http.StatusOK, `{"access_token":"bot-access-token"}`), nil
		})},
	})
	if err != nil {
		t.Fatalf("create bot OAuth client: %v", err)
	}

	// When
	authorizationURL, err := client.AuthorizationURL("bot-state", nil)
	if err != nil {
		t.Fatalf("build bot authorization URL: %v", err)
	}
	tokens, err := client.ExchangeCode(context.Background(), "bot-authorization-code")
	if err != nil {
		t.Fatalf("exchange bot authorization code: %v", err)
	}

	// Then
	parsedAuthorizationURL, err := url.Parse(authorizationURL)
	if err != nil {
		t.Fatalf("parse bot authorization URL: %v", err)
	}
	if got := parsedAuthorizationURL.Query().Get("redirect_uri"); got != botRedirectURL {
		t.Errorf("bot authorization redirect_uri = %q, want %q", got, botRedirectURL)
	}
	if parsedAuthorizationURL.Query().Has("scope") {
		t.Errorf("bot authorization URL must not guess a scope: %s", authorizationURL)
	}
	if tokens.AccessToken != "bot-access-token" {
		t.Errorf("bot access token = %q", tokens.AccessToken)
	}
}
