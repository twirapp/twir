package mcp_oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

func TestHandler_revoke_binds_token_to_public_client(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	form := url.Values{"client_id": {"client"}, "token": {"token"}}

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/revoke", strings.NewReader(form.Encode()), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	// Then
	if response.Code != http.StatusOK || handler.service.revoked.ClientID != "client" || handler.service.revoked.Token != "token" {
		t.Fatalf("revoke response = %d, input = %#v", response.Code, handler.service.revoked)
	}
}

func TestHandler_authorize_redirects_validation_errors_only_to_trusted_callback(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	query := url.Values{"client_id": {"client"}, "redirect_uri": {"https://client.example/callback"}, "state": {"opaque-state"}}
	handler.service.validateErr = &service.OAuthError{Code: service.ErrorInvalidScope, Description: "requested scope is not permitted"}

	// When
	response := serve(handler.router(), http.MethodGet, "/oauth/authorize?"+query.Encode(), nil, nil)

	// Then
	location, err := url.Parse(response.Header().Get("Location"))
	if err != nil || response.Code != http.StatusFound || location.String()[:len("https://client.example/callback")] != "https://client.example/callback" || location.Query().Get("error") != "invalid_scope" || location.Query().Get("error_description") != "requested scope is not permitted" || location.Query().Get("state") != "opaque-state" {
		t.Fatalf("authorize response = %d %q", response.Code, response.Header().Get("Location"))
	}
}

func TestHandler_authorize_never_redirects_unknown_client_or_unregistered_uri(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	unknown := url.Values{"client_id": {"unknown"}, "redirect_uri": {"https://client.example/callback"}, "state": {"state"}}
	unregistered := url.Values{"client_id": {"client"}, "redirect_uri": {"https://attacker.example/callback"}, "state": {"state"}}
	handler.service.getClientErr = &service.OAuthError{Code: service.ErrorInvalidClient, Description: "unknown client"}

	// When
	unknownResponse := serve(handler.router(), http.MethodGet, "/oauth/authorize?"+unknown.Encode(), nil, nil)
	handler.service.getClientErr = nil
	unregisteredResponse := serve(handler.router(), http.MethodGet, "/oauth/authorize?"+unregistered.Encode(), nil, nil)

	// Then
	if unknownResponse.Code != http.StatusUnauthorized || unknownResponse.Header().Get("Location") != "" || unregisteredResponse.Code != http.StatusBadRequest || unregisteredResponse.Header().Get("Location") != "" {
		t.Fatalf("authorize responses = %d %q, %d %q", unknownResponse.Code, unknownResponse.Header().Get("Location"), unregisteredResponse.Code, unregisteredResponse.Header().Get("Location"))
	}
}

func TestHandler_authorize_redirect_validation_matches_loopback_dynamic_port(t *testing.T) {
	tests := []struct {
		name          string
		registeredURI string
		requestedURI  string
		wantRedirect  bool
	}{
		{name: "allows changed port", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/callback?source=mcp", wantRedirect: true},
		{name: "rejects changed loopback host", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.2:49321/callback?source=mcp"},
		{name: "rejects changed path", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/other?source=mcp"},
		{name: "rejects changed query", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/callback?source=other"},
		{name: "rejects changed HTTPS port", registeredURI: "https://client.example:3000/callback", requestedURI: "https://client.example:49321/callback"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			handler := newTestHandler(t)
			handler.service.client.RedirectURIs = []string{test.registeredURI}
			handler.service.validateErr = &service.OAuthError{Code: service.ErrorInvalidScope, Description: "requested scope is not permitted"}
			query := url.Values{"client_id": {"client"}, "redirect_uri": {test.requestedURI}, "state": {"opaque-state"}}

			// When
			response := serve(handler.router(), http.MethodGet, "/oauth/authorize?"+query.Encode(), nil, nil)

			// Then
			if !test.wantRedirect {
				if response.Code != http.StatusBadRequest || response.Header().Get("Location") != "" {
					t.Fatalf("authorize response = %d %q", response.Code, response.Header().Get("Location"))
				}
				return
			}
			location, err := url.Parse(response.Header().Get("Location"))
			if err != nil || response.Code != http.StatusFound || location.Scheme != "http" || location.Host != "127.0.0.1:49321" || location.Path != "/callback" || location.Query().Get("source") != "mcp" || location.Query().Get("error") != "invalid_scope" || location.Query().Get("state") != "opaque-state" {
				t.Fatalf("authorize response = %d %q", response.Code, response.Header().Get("Location"))
			}
		})
	}
}

func TestHandler_consent_rejects_dashboard_changed_since_display(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	displayedChannelID := handler.sessions.dashboardID
	handler.sessions.dashboardID = uuid.New()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","access_level":"read"}`, displayedChannelID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	if response.Code != http.StatusConflict || handler.service.created.ChannelID != uuid.Nil {
		t.Fatalf("consent response = %d, created = %#v", response.Code, handler.service.created)
	}
	if _, exists := handler.sessions.attempts["attempt"]; !exists {
		t.Fatal("dashboard mismatch consumed authorization attempt")
	}
}

func TestHandler_public_oauth_endpoints_support_preflight_without_credentials(t *testing.T) {
	for _, path := range []string{"/oauth/register", "/oauth/token", "/oauth/revoke"} {
		t.Run(path, func(t *testing.T) {
			// Given
			handler := newTestHandler(t)

			// When
			response := serve(handler.router(), http.MethodOptions, path, nil, map[string]string{"Origin": "https://client.example", "Access-Control-Request-Method": "POST"})

			// Then
			if response.Code != http.StatusNoContent || response.Header().Get("Access-Control-Allow-Origin") != "*" || response.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS" || response.Header().Get("Access-Control-Allow-Headers") != "Content-Type" || response.Header().Get("Access-Control-Max-Age") != "600" || response.Header().Get("Access-Control-Allow-Credentials") != "" {
				t.Fatalf("preflight response = %d, headers = %#v", response.Code, response.Header())
			}
		})
	}
}
