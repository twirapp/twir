package mcp_oauth

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

func TestHandler_token_passes_pkce_and_resource_and_disables_caching(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {"client"},
		"code":          {"one-use-code"},
		"redirect_uri":  {"https://client.example/callback"},
		"code_verifier": {"verifier"},
		"resource":      {"https://twir.example/api/mcp"},
	}

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	// Then
	cacheControl := response.Header().Get("Cache-Control")
	exchangedCodeVerifierMatches := handler.service.exchanged.CodeVerifier == "verifier"
	exchangedResourceMatches := handler.service.exchanged.Resource == "https://twir.example/api/mcp"
	if response.Code != http.StatusOK || cacheControl != "no-store" || !exchangedCodeVerifierMatches || !exchangedResourceMatches {
		t.Fatalf("token response = %d, headers = %#v, exchange = %#v", response.Code, response.Header(), handler.service.exchanged)
	}
}

func TestHandler_token_accepts_omitted_resource_for_authorization_code_and_refresh_grants(t *testing.T) {
	for _, test := range []struct {
		name      string
		grantType string
		form      url.Values
		resource  func(*fakeService) string
	}{
		{
			name:      "authorization code",
			grantType: "authorization_code",
			form: url.Values{
				"grant_type":    {"authorization_code"},
				"client_id":     {"client"},
				"code":          {"one-use-code"},
				"redirect_uri":  {"https://client.example/callback"},
				"code_verifier": {"verifier"},
			},
			resource: func(service *fakeService) string { return service.exchanged.Resource },
		},
		{
			name:      "refresh token",
			grantType: "refresh_token",
			form: url.Values{
				"grant_type":    {"refresh_token"},
				"client_id":     {"client"},
				"refresh_token": {"refresh"},
			},
			resource: func(service *fakeService) string { return service.refreshed.Resource },
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			handler := newTestHandler(t)

			// When
			response := serve(handler.router(), http.MethodPost, "/oauth/token", strings.NewReader(test.form.Encode()),
				map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

			// Then
			if response.Code != http.StatusOK || test.resource(handler.service) != "" {
				t.Fatalf("%s token response = %d, input = %#v", test.grantType, response.Code, handler.service)
			}
		})
	}
}

func TestHandler_token_rejects_client_secret_and_refresh_reuse(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	secretForm := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {"client"},
		"refresh_token": {"refresh"},
		"resource":      {"https://twir.example/api/mcp"},
		"client_secret": {"forbidden"},
	}
	reuseForm := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {"client"},
		"refresh_token": {"reused"},
		"resource":      {"https://twir.example/api/mcp"},
	}
	handler.service.err = &service.OAuthError{Code: service.ErrorInvalidGrant, Description: "invalid refresh token"}

	// When
	secret := serve(handler.router(), http.MethodPost, "/oauth/token", strings.NewReader(secretForm.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	reuse := serve(handler.router(), http.MethodPost, "/oauth/token", strings.NewReader(reuseForm.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	// Then
	if secret.Code != http.StatusUnauthorized || reuse.Code != http.StatusBadRequest || !strings.Contains(reuse.Body.String(), "invalid_grant") {
		t.Fatalf("token statuses = %d, %d: %s", secret.Code, reuse.Code, reuse.Body.String())
	}
}

func TestHandler_revoke_is_idempotent_for_unknown_token(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	form := url.Values{"client_id": {"client"}, "token": {"unknown"}}

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/revoke", strings.NewReader(form.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	// Then
	cacheControl := response.Header().Get("Cache-Control")
	allowOrigin := response.Header().Get("Access-Control-Allow-Origin")
	if response.Code != http.StatusOK || cacheControl != "no-store" || allowOrigin != "*" || handler.service.revocations != 1 {
		t.Fatalf("revoke response = %d, headers = %#v, calls = %d", response.Code, response.Header(), handler.service.revocations)
	}
}
