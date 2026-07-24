package vk

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type oauthRoundTripFunc func(*http.Request) (*http.Response, error)

func (f oauthRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestOAuthClientAuthorizationURLUsesDevAPIContract(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("authorization URL must not make an HTTP request")
		return nil, nil
	}))

	authorizationURL, err := client.AuthorizationURL("state-value", []string{"user_info", "video_list"})
	if err != nil {
		t.Fatalf("build authorization URL: %v", err)
	}

	parsed, err := url.Parse(authorizationURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if parsed.Scheme != "https" || parsed.Host != "auth.example.test" || parsed.Path != "/app/oauth2/authorize" {
		t.Fatalf("unexpected authorization endpoint: %s", parsed)
	}

	assertQueryValues(t, parsed.Query(), map[string]string{
		"response_type": "code",
		"client_id":     "client-id",
		"redirect_uri":  "https://twir.example.test/auth/vk/callback",
		"state":         "state-value",
		"scope":         "user_info,video_list",
	})
}

func TestNewOAuthClientRequiresClientSecretBeforeAnyRequest(t *testing.T) {
	requests := 0
	client, err := NewOAuthClient(OAuthClientOpts{
		ClientID:    "client-id",
		RedirectURL: "https://twir.example.test/auth/vk/callback",
		HTTPClient: &http.Client{Transport: oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requests++
			return nil, errors.New("unexpected HTTP request")
		})},
	})
	if err == nil {
		t.Fatal("expected blank client secret to be rejected")
	}
	if client != nil {
		t.Fatal("expected no client when client secret is blank")
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests, got %d", requests)
	}
}

func TestNewOAuthClientRejectsInvalidBaseURLs(t *testing.T) {
	if _, err := NewOAuthClient(OAuthClientOpts{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://twir.example.test/auth/vk/callback",
		APIBaseURL:   "not-a-url",
	}); err == nil {
		t.Fatal("expected invalid API base URL to be rejected")
	}

	if _, err := NewOAuthClient(OAuthClientOpts{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://twir.example.test/auth/vk/callback",
		AuthBaseURL:  "not-a-url",
	}); err == nil {
		t.Fatal("expected invalid auth base URL to be rejected")
	}
}

func TestOAuthClientExchangeCodeSendsDevAPIForm(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertTokenEndpoint(t, req)
		form := parseRequestForm(t, req)
		assertExactFormValues(t, form, map[string]string{
			"grant_type":   "authorization_code",
			"code":         "authorization-code",
			"redirect_uri": "https://twir.example.test/auth/vk/callback",
		})
		return oauthResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"token_type":"Bearer"}`), nil
	}))

	tokens, err := client.ExchangeCode(context.Background(), "authorization-code")
	if err != nil {
		t.Fatalf("exchange VK Video code: %v", err)
	}
	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" || tokens.ExpiresIn != 3600 {
		t.Fatalf("unexpected token response: %#v", tokens)
	}
}

func TestOAuthClientExchangeCodeRejectsMissingAccessToken(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusOK, `{"refresh_token":"refresh-token","expires_in":3600,"token_type":"Bearer"}`), nil
	}))

	if _, err := client.ExchangeCode(context.Background(), "authorization-code"); err == nil {
		t.Fatal("expected token exchange without access_token to fail")
	}
}

func TestOAuthClientRefreshTokenSendsDevAPIForm(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertTokenEndpoint(t, req)
		form := parseRequestForm(t, req)
		assertExactFormValues(t, form, map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": "refresh-token",
			"redirect_uri":  "https://twir.example.test/auth/vk/callback",
		})
		return oauthResponse(http.StatusOK, `{"access_token":"new-access-token","expires_in":7200,"token_type":"Bearer"}`), nil
	}))

	tokens, err := client.RefreshToken(context.Background(), "refresh-token")
	if err != nil {
		t.Fatalf("refresh VK Video token: %v", err)
	}
	if tokens.AccessToken != "new-access-token" || tokens.RefreshToken != "" || tokens.ExpiresIn != 7200 {
		t.Fatalf("unexpected token response: %#v", tokens)
	}
}

func TestOAuthClientRefreshTokenRejectsMissingAccessToken(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusOK, `{"refresh_token":"refresh-token","expires_in":3600,"token_type":"Bearer"}`), nil
	}))

	if _, err := client.RefreshToken(context.Background(), "refresh-token"); err == nil {
		t.Fatal("expected token refresh without access_token to fail")
	}
}

func TestOAuthClientCurrentUserMapsDevAPIProfile(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Errorf("request method = %s, want GET", req.Method)
		}
		if req.URL.Scheme != "https" || req.URL.Host != "devapi.example.test" || req.URL.Path != "/v1/current_user" {
			t.Errorf("unexpected endpoint: %s", req.URL)
		}
		if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
			t.Errorf("authorization header = %q, want Bearer access-token", auth)
		}

		return oauthResponse(http.StatusOK, `{"data":{"user":{"id":123,"nick":"ada","avatar_url":"https://cdn.example.test/avatar.jpg","is_streamer":true},"channel":{"url":"ada"}}}`), nil
	}))

	user, err := client.CurrentUser(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get VK Video profile: %v", err)
	}
	if user.ID != "123" || user.Nick != "ada" || user.Avatar != "https://cdn.example.test/avatar.jpg" {
		t.Fatalf("unexpected VK Video profile: %#v", user)
	}
}

func TestOAuthClientCurrentUserAcceptsStringID(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusOK, `{"data":{"user":{"id":"123","nick":"ada","avatar_url":""}}}`), nil
	}))

	user, err := client.CurrentUser(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get VK Video profile: %v", err)
	}
	if user.ID != "123" {
		t.Fatalf("unexpected VK Video profile: %#v", user)
	}
}

func TestOAuthClientReturnsTypedProviderError(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusBadRequest, `{"error":"invalid_grant","error_description":"authorization code is invalid","state":"request-state"}`), nil
	}))

	_, err := client.ExchangeCode(context.Background(), "authorization-code")

	var providerErr *ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T (%v)", err, err)
	}
	if providerErr.Code != "invalid_grant" || providerErr.Description != "authorization code is invalid" || providerErr.State != "request-state" || providerErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected provider error: %#v", providerErr)
	}
}

func TestOAuthClientReturnsTypedErrorForNon2xxResponse(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusInternalServerError, `{"message":"unexpected"}`), nil
	}))

	_, err := client.CurrentUser(context.Background(), "access-token")

	var providerErr *ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T (%v)", err, err)
	}
	if providerErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("provider error status = %d, want %d", providerErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestOAuthClientRejectsMalformedJSON(t *testing.T) {
	client := newTestOAuthClient(t, oauthRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return oauthResponse(http.StatusOK, `{malformed`), nil
	}))

	_, err := client.CurrentUser(context.Background(), "access-token")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}

	var providerErr *ProviderError
	if errors.As(err, &providerErr) {
		t.Fatalf("malformed JSON must not become a provider error: %#v", providerErr)
	}
}

func newTestOAuthClient(t *testing.T, transport http.RoundTripper) *OAuthClient {
	t.Helper()

	client, err := NewOAuthClient(OAuthClientOpts{
		ClientID:      "client-id",
		ClientSecret:  "client-secret",
		RedirectURL:   "https://twir.example.test/auth/vk/callback",
		APIBaseURL:    "https://api.example.test",
		AuthBaseURL:   "https://auth.example.test",
		DevAPIBaseURL: "https://devapi.example.test",
		HTTPClient:    &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("create VK Video OAuth client: %v", err)
	}

	return client
}

func oauthResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func assertTokenEndpoint(t *testing.T, req *http.Request) {
	t.Helper()
	if req.Method != http.MethodPost {
		t.Errorf("request method = %s, want POST", req.Method)
	}
	if req.URL.Scheme != "https" || req.URL.Host != "api.example.test" || req.URL.Path != "/oauth/server/token" {
		t.Errorf("unexpected endpoint: %s", req.URL)
	}
	if contentType := req.Header.Get("Content-Type"); contentType != "application/x-www-form-urlencoded" {
		t.Errorf("content type = %q, want application/x-www-form-urlencoded", contentType)
	}

	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("client-id:client-secret"))
	if auth := req.Header.Get("Authorization"); auth != wantAuth {
		t.Errorf("authorization header = %q, want %q", auth, wantAuth)
	}
}

func parseRequestForm(t *testing.T, req *http.Request) url.Values {
	t.Helper()
	if err := req.ParseForm(); err != nil {
		t.Fatalf("parse request form: %v", err)
	}
	return req.PostForm
}

func assertQueryValues(t *testing.T, query url.Values, want map[string]string) {
	t.Helper()
	for key, expected := range want {
		if got := query.Get(key); got != expected {
			t.Errorf("query %s = %q, want %q", key, got, expected)
		}
	}
}

func assertExactFormValues(t *testing.T, form url.Values, want map[string]string) {
	t.Helper()
	if len(form) != len(want) {
		t.Errorf("form field count = %d, want %d: %#v", len(form), len(want), form)
	}

	for key, values := range form {
		expected, ok := want[key]
		if !ok {
			t.Errorf("form contains unexpected field %q", key)
			continue
		}
		if len(values) != 1 {
			t.Errorf("form %s has %d values, want 1", key, len(values))
			continue
		}
		if values[0] != expected {
			t.Errorf("form %s = %q, want %q", key, values[0], expected)
		}
	}

	for key, expected := range want {
		if _, ok := form[key]; !ok {
			t.Errorf("form is missing %s=%q", key, expected)
		}
	}
}
