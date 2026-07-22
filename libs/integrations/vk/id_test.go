package vk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type idRoundTripFunc func(*http.Request) (*http.Response, error)

func (f idRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestIDClientAuthorizationURLUsesVKIDPKCEContract(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("authorization URL must not make an HTTP request")
		return nil, nil
	}))

	authorizationURL, err := client.AuthorizationURL("state-value", "pkce-challenge")
	if err != nil {
		t.Fatalf("build authorization URL: %v", err)
	}

	parsed, err := url.Parse(authorizationURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if parsed.Scheme != "https" || parsed.Host != "id.example.test" || parsed.Path != "/authorize" {
		t.Fatalf("unexpected authorization endpoint: %s", parsed)
	}

	assertFormValues(t, parsed.Query(), map[string]string{
		"response_type":         "code",
		"client_id":             "client-id",
		"redirect_uri":          "https://twir.example.test/auth/vk/callback",
		"state":                 "state-value",
		"code_challenge":        "pkce-challenge",
		"code_challenge_method": "S256",
		"scope":                 "vkid.personal_info",
	})
}

func TestIDClientExchangeCodeSendsVKIDForm(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertIDEndpoint(t, req, "/oauth2/auth")
		form := parseRequestForm(t, req)
		assertFormValues(t, form, map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     "client-id",
			"code":          "authorization-code",
			"code_verifier": "pkce-verifier",
			"redirect_uri":  "https://twir.example.test/auth/vk/callback",
			"device_id":     "device-id",
			"service_token": "service-token",
		})
		if form.Has("client_secret") {
			t.Error("VK ID form must not send client_secret")
		}

		return idResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"scope":"vkid.personal_info"}`), nil
	}))

	tokens, err := client.ExchangeCode(context.Background(), IDExchangeCodeInput{
		Code:         "authorization-code",
		CodeVerifier: "pkce-verifier",
		DeviceID:     "device-id",
	})
	if err != nil {
		t.Fatalf("exchange VK ID code: %v", err)
	}
	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" || tokens.ExpiresIn != 3600 {
		t.Fatalf("unexpected token response: %#v", tokens)
	}
	if got, want := tokens.Scopes, []string{"vkid.personal_info"}; !equalStrings(got, want) {
		t.Fatalf("token scopes = %#v, want %#v", got, want)
	}
}

func TestIDClientExchangeCodeRequiresDeviceIDBeforeRequest(t *testing.T) {
	requests := 0
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return nil, errors.New("unexpected HTTP request")
	}))

	_, err := client.ExchangeCode(context.Background(), IDExchangeCodeInput{
		Code:         "authorization-code",
		CodeVerifier: "pkce-verifier",
	})
	if !errors.Is(err, ErrDeviceIDRequired) {
		t.Fatalf("expected missing device ID error, got %v", err)
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests, got %d", requests)
	}
}

func TestIDClientRefreshTokenSendsVKIDForm(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertIDEndpoint(t, req, "/oauth2/auth")
		form := parseRequestForm(t, req)
		assertFormValues(t, form, map[string]string{
			"grant_type":    "refresh_token",
			"client_id":     "client-id",
			"refresh_token": "refresh-token",
			"device_id":     "device-id",
			"service_token": "service-token",
		})
		if form.Has("client_secret") {
			t.Error("VK ID form must not send client_secret")
		}

		return idResponse(http.StatusOK, `{"access_token":"new-access-token","expires_in":7200,"scope":"vkid.personal_info"}`), nil
	}))

	tokens, err := client.RefreshToken(context.Background(), IDRefreshTokenInput{
		RefreshToken: "refresh-token",
		DeviceID:     "device-id",
	})
	if err != nil {
		t.Fatalf("refresh VK ID token: %v", err)
	}
	if tokens.AccessToken != "new-access-token" || tokens.RefreshToken != "" || tokens.ExpiresIn != 7200 {
		t.Fatalf("unexpected token response: %#v", tokens)
	}
}

func TestIDClientRefreshTokenRequiresDeviceIDBeforeRequest(t *testing.T) {
	requests := 0
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return nil, errors.New("unexpected HTTP request")
	}))

	_, err := client.RefreshToken(context.Background(), IDRefreshTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, ErrDeviceIDRequired) {
		t.Fatalf("expected missing device ID error, got %v", err)
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests, got %d", requests)
	}
}

func TestIDClientUserInfoMapsVKIDProfile(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertIDEndpoint(t, req, "/oauth2/user_info")
		form := parseRequestForm(t, req)
		assertFormValues(t, form, map[string]string{
			"client_id":    "client-id",
			"access_token": "access-token",
		})
		if len(form) != 2 {
			t.Errorf("profile request has unexpected fields: %#v", form)
		}

		return idResponse(http.StatusOK, `{"user":{"user_id":"123","first_name":"Ada","last_name":"Lovelace","avatar":"https://cdn.example.test/avatar.jpg"}}`), nil
	}))

	user, err := client.UserInfo(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("get VK ID profile: %v", err)
	}
	if user.ID != "123" || user.FirstName != "Ada" || user.LastName != "Lovelace" || user.Avatar != "https://cdn.example.test/avatar.jpg" {
		t.Fatalf("unexpected VK ID profile: %#v", user)
	}
}

func TestIDClientReturnsTypedProviderError(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return idResponse(http.StatusBadRequest, `{"error":"invalid_grant","error_description":"authorization code is invalid","state":"request-state"}`), nil
	}))

	_, err := client.ExchangeCode(context.Background(), IDExchangeCodeInput{
		Code:         "authorization-code",
		CodeVerifier: "pkce-verifier",
		DeviceID:     "device-id",
	})

	var providerErr *ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T (%v)", err, err)
	}
	if providerErr.Code != "invalid_grant" || providerErr.Description != "authorization code is invalid" || providerErr.State != "request-state" || providerErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected provider error: %#v", providerErr)
	}
}

func TestIDClientReturnsTypedErrorForNon2xxResponse(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return idResponse(http.StatusInternalServerError, `{"message":"unexpected"}`), nil
	}))

	_, err := client.UserInfo(context.Background(), "access-token")

	var providerErr *ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T (%v)", err, err)
	}
	if providerErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("provider error status = %d, want %d", providerErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestIDClientRejectsMalformedJSON(t *testing.T) {
	client := newTestIDClient(t, idRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return idResponse(http.StatusOK, `{malformed`), nil
	}))

	_, err := client.UserInfo(context.Background(), "access-token")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}

	var providerErr *ProviderError
	if errors.As(err, &providerErr) {
		t.Fatalf("malformed JSON must not become a provider error: %#v", providerErr)
	}
}

func newTestIDClient(t *testing.T, transport http.RoundTripper) *IDClient {
	t.Helper()

	client, err := NewIDClient(IDClientOpts{
		ClientID:     "client-id",
		ServiceToken: "service-token",
		RedirectURL:  "https://twir.example.test/auth/vk/callback",
		APIBaseURL:   "https://id.example.test",
		HTTPClient:   &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("create VK ID client: %v", err)
	}

	return client
}

func idResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func assertIDEndpoint(t *testing.T, req *http.Request, path string) {
	t.Helper()
	if req.Method != http.MethodPost {
		t.Errorf("request method = %s, want POST", req.Method)
	}
	if req.URL.Scheme != "https" || req.URL.Host != "id.example.test" || req.URL.Path != path {
		t.Errorf("unexpected endpoint: %s", req.URL)
	}
	if contentType := req.Header.Get("Content-Type"); contentType != "application/x-www-form-urlencoded" {
		t.Errorf("content type = %q, want application/x-www-form-urlencoded", contentType)
	}
}

func parseRequestForm(t *testing.T, req *http.Request) url.Values {
	t.Helper()
	if err := req.ParseForm(); err != nil {
		t.Fatalf("parse request form: %v", err)
	}
	return req.PostForm
}

func assertFormValues(t *testing.T, form url.Values, want map[string]string) {
	t.Helper()
	for key, expected := range want {
		if got := form.Get(key); got != expected {
			t.Errorf("form %s = %q, want %q", key, got, expected)
		}
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
