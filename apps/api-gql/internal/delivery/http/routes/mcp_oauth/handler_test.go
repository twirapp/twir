package mcp_oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

func TestHandler_metadata_has_discovery_headers_and_contract(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	router := handler.router()

	// When
	response := serve(router, http.MethodGet, "/.well-known/oauth-authorization-server", nil, nil)
	duplicate := serve(router, http.MethodGet, "/api/.well-known/oauth-authorization-server", nil, nil)

	// Then
	allowOrigin := response.Header().Get("Access-Control-Allow-Origin")
	cacheControl := response.Header().Get("Cache-Control")
	if response.Code != http.StatusOK || allowOrigin != "*" || cacheControl != "public, max-age=3600" || duplicate.Code != http.StatusNotFound {
		t.Fatalf("metadata response = %d, headers = %#v", response.Code, response.Header())
	}
	var metadata struct {
		Issuer                string   `json:"issuer"`
		AuthorizationEndpoint string   `json:"authorization_endpoint"`
		CodeChallengeMethods  []string `json:"code_challenge_methods_supported"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &metadata); err != nil {
		t.Fatal(err)
	}
	issuerMatches := metadata.Issuer == "https://twir.example"
	authorizationEndpointMatches := metadata.AuthorizationEndpoint == "https://twir.example/api/oauth/authorize"
	codeChallengeMethodsMatch := len(metadata.CodeChallengeMethods) == 1 && metadata.CodeChallengeMethods[0] == "S256"
	if !issuerMatches || !authorizationEndpointMatches || !codeChallengeMethodsMatch {
		t.Fatalf("metadata = %#v", metadata)
	}
}

func TestHandler_register_returns_normalized_public_client_without_secret(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	payload := `{"client_name":"Example","redirect_uris":["https://client.example/callback"],"client_secret":"forbidden"}`

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/register", strings.NewReader(payload), map[string]string{"Content-Type": "application/json"})

	// Then
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid_client_metadata") {
		t.Fatalf("registration response = %d %s", response.Code, response.Body.String())
	}
}

func TestHandler_authorize_stores_sensitive_values_server_side_and_redirects_with_opaque_attempt(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	query := url.Values{
		"response_type":         {"code"},
		"client_id":             {"client"},
		"redirect_uri":          {"https://client.example/callback"},
		"scope":                 {"read write"},
		"state":                 {"client-state"},
		"code_challenge":        {"challenge"},
		"code_challenge_method": {"S256"},
		"resource":              {"https://twir.example/api/mcp"},
	}

	// When
	response := serve(handler.router(), http.MethodGet, "/oauth/authorize?"+query.Encode(), nil, nil)

	// Then
	location, err := url.Parse(response.Header().Get("Location"))
	contentType := response.Header().Get("Content-Type")
	locationHeader := response.Header().Get("Location")
	foundBody := strings.Contains(response.Body.String(), "Found")
	attemptPresent := location.Query().Get("attempt") != ""
	clientStateAbsent := !strings.Contains(locationHeader, "client-state")
	challengeAbsent := !strings.Contains(locationHeader, "challenge")
	responseFound := response.Code == http.StatusFound
	contentTypeMatches := contentType == "text/html; charset=utf-8"
	if err != nil || !responseFound || !contentTypeMatches || !foundBody || !attemptPresent || !clientStateAbsent || !challengeAbsent {
		t.Fatalf("authorize response = %d %q", response.Code, response.Header().Get("Location"))
	}
	attempt, err := handler.sessions.GetMCPOAuthAttempt(context.Background(), location.Query().Get("attempt"))
	if err != nil || attempt.ClientState != "client-state" || attempt.CodeChallenge != "challenge" || attempt.RedirectURI != "https://client.example/callback" {
		t.Fatalf("attempt = %#v, err = %v", attempt, err)
	}
}

func TestHandler_consent_requires_origin_csrf_permission_and_is_one_use(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	attemptID := "attempt"
	handler.sessions.attempts[attemptID] = testAttempt()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","access_level":"write"}`, handler.sessions.dashboardID)

	// When
	badOrigin := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://evil.example"})
	allowed := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})
	replay := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	if badOrigin.Code != http.StatusForbidden || allowed.Code != http.StatusOK || replay.Code != http.StatusNotFound {
		t.Fatalf("consent statuses = %d, %d, %d", badOrigin.Code, allowed.Code, replay.Code)
	}
	createdScopeMatches := handler.service.created.Authorize.Scope == "read write"
	createdChannelMatches := handler.service.created.ChannelID == handler.sessions.dashboardID
	createdApproverMatches := handler.service.created.ApprovingUserID == handler.sessions.userID
	if !createdScopeMatches || !createdChannelMatches || !createdApproverMatches {
		t.Fatalf("created code input = %#v", handler.service.created)
	}
}

func TestHandler_consent_denial_redirects_to_validated_callback(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"deny"}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	bodyText := response.Body.String()
	if response.Code != http.StatusOK || !strings.Contains(bodyText, "access_denied") || !strings.Contains(bodyText, "client-state") {
		t.Fatalf("denial response = %d %s", response.Code, response.Body.String())
	}
}

func TestHandler_consent_rejects_missing_browser_session_and_invalid_csrf(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	handler.sessions.userErr = errors.New("not signed in")
	unauthorized := serve(handler.router(), http.MethodGet, "/oauth/consent?attempt=attempt", nil, nil)
	handler.sessions.userErr = nil
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"wrong","decision":"approve","access_level":"read"}`, handler.sessions.dashboardID)

	// When
	csrf := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	if unauthorized.Code != http.StatusUnauthorized || csrf.Code != http.StatusForbidden {
		t.Fatalf("consent statuses = %d, %d", unauthorized.Code, csrf.Code)
	}
	if _, exists := handler.sessions.attempts["attempt"]; !exists {
		t.Fatal("invalid CSRF consumed authorization attempt")
	}
}

func TestHandler_consent_returns_access_denied_when_service_permission_is_lost(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	handler.service.err = &service.OAuthError{Code: service.ErrorAccessDenied, Description: "permission is required"}
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","access_level":"read"}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	if response.Code != http.StatusForbidden || !strings.Contains(response.Body.String(), "access_denied") {
		t.Fatalf("permission response = %d %s", response.Code, response.Body.String())
	}
}

func testAttempt() authsessions.MCPOAuthAttempt {
	return authsessions.MCPOAuthAttempt{
		ClientID:        "client",
		RedirectURI:     "https://client.example/callback",
		ClientState:     "client-state",
		CodeChallenge:   "challenge",
		RequestedScopes: []string{"read", "write"},
		Resource:        "https://twir.example/api/mcp",
		CSRFToken:       "csrf",
		ExpiresAt:       time.Now().Add(time.Minute),
	}
}

func serve(handler http.Handler, method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}
