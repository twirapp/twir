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
	"reflect"
	"strings"
	"testing"
	"time"

	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
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

func TestHandler_discovery_metadata_lists_exact_canonical_scopes(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	wantScopes := []string{
		"commands:read", "commands:edit", "timers:read", "timers:edit",
		"files:read", "files:edit", "games:read", "games:edit",
		"song_requests:read", "song_requests:edit", "moderation:read", "moderation:edit",
		"overlays:read", "overlays:edit", "integrations:read", "integrations:edit",
		"events:read", "events:edit", "rewards:read", "rewards:edit",
		"giveaways:read", "giveaways:edit", "greetings:read", "greetings:edit",
		"notifications:read", "notifications:edit", "alerts:read", "alerts:edit",
		"secrets:read", "secrets:edit", "storage:read", "storage:edit",
		"pastes:read", "pastes:edit", "short_urls:read", "short_urls:edit",
		"dashboard:read", "dashboard:edit", "variables:read", "variables:edit",
		"quotes:read", "quotes:edit", "keywords:read", "keywords:edit",
	}

	for _, path := range []string{
		"/.well-known/oauth-protected-resource",
		"/.well-known/oauth-authorization-server",
	} {
		t.Run(path, func(t *testing.T) {
			// When
			response := serve(handler.router(), http.MethodGet, path, nil, nil)

			// Then
			if response.Code != http.StatusOK {
				t.Fatalf("metadata response = %d", response.Code)
			}
			var metadata struct {
				ScopesSupported []string `json:"scopes_supported"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &metadata); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(metadata.ScopesSupported, wantScopes) {
				t.Fatalf("scopes_supported = %v, want %v", metadata.ScopesSupported, wantScopes)
			}
		})
	}
}

func TestHandler_scope_catalog_returns_public_canonical_groups(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	wantScopes := []scopeCatalogItemResponse{
		{Group: "commands", Name: "Commands", Description: "View and manage custom commands, groups, and role cooldowns", Actions: []string{"read", "edit"}},
		{Group: "timers", Name: "Timers", Description: "View and manage chat timers", Actions: []string{"read", "edit"}},
		{Group: "files", Name: "Files", Description: "View and manage uploaded image and audio files", Actions: []string{"read", "edit"}},
		{Group: "games", Name: "Games", Description: "View and manage channel games and their settings", Actions: []string{"read", "edit"}},
		{Group: "song_requests", Name: "Song Requests", Description: "View and manage the song request queue and playback", Actions: []string{"read", "edit"}},
		{Group: "moderation", Name: "Moderation", Description: "View and manage channel moderation rules and chat wall settings", Actions: []string{"read", "edit"}},
		{Group: "overlays", Name: "Overlays", Description: "View and manage custom and built-in overlay settings", Actions: []string{"read", "edit"}},
		{Group: "integrations", Name: "Integrations", Description: "View and manage connected third-party integrations", Actions: []string{"read", "edit"}},
		{Group: "events", Name: "Events", Description: "View and manage channel automation events and operations", Actions: []string{"read", "edit"}},
		{Group: "rewards", Name: "Rewards", Description: "View and manage Twitch custom rewards", Actions: []string{"read", "edit"}},
		{Group: "giveaways", Name: "Giveaways", Description: "View and manage channel giveaways", Actions: []string{"read", "edit"}},
		{Group: "greetings", Name: "Greetings", Description: "View and manage channel greetings", Actions: []string{"read", "edit"}},
		{Group: "notifications", Name: "Notifications", Description: "View and manage channel notifications", Actions: []string{"read", "edit"}},
		{Group: "alerts", Name: "Alerts", Description: "View and manage channel alerts and their bindings", Actions: []string{"read", "edit"}},
		{Group: "secrets", Name: "Secrets", Description: "View and manage encrypted channel secrets", Actions: []string{"read", "edit"}},
		{Group: "storage", Name: "Storage", Description: "View and manage channel JSON storage entries", Actions: []string{"read", "edit"}},
		{Group: "pastes", Name: "Pastes", Description: "View and manage owned pastes", Actions: []string{"read", "edit"}},
		{Group: "short_urls", Name: "Short URLs", Description: "View and manage owned short URLs and their settings", Actions: []string{"read", "edit"}},
		{Group: "dashboard", Name: "Dashboard", Description: "View and manage channel dashboard settings and statistics", Actions: []string{"read", "edit"}},
		{Group: "variables", Name: "Variables", Description: "View and manage channel custom variables", Actions: []string{"read", "edit"}},
		{Group: "quotes", Name: "Quotes", Description: "View and manage channel quotes", Actions: []string{"read", "edit"}},
		{Group: "keywords", Name: "Keywords", Description: "View and manage keyword triggers", Actions: []string{"read", "edit"}},
	}

	// When
	response := serve(handler.router(), http.MethodGet, "/oauth/scopes", nil, nil)

	// Then
	if response.Code != http.StatusOK || response.Header().Get("Access-Control-Allow-Origin") != "*" ||
		response.Header().Get("Cache-Control") != "public, max-age=3600" {
		t.Fatalf("scope catalog response = %d, headers = %#v", response.Code, response.Header())
	}
	var body struct {
		Scopes []scopeCatalogItemResponse `json:"scopes"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body.Scopes, wantScopes) {
		t.Fatalf("scope catalog = %#v, want %#v", body.Scopes, wantScopes)
	}
}

func TestHandler_consent_get_groups_requested_scopes_in_catalog_order(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = authsessions.MCPOAuthAttempt{
		ClientID:        "client",
		RequestedScopes: []string{"files:edit", "commands:read", "timers:read", "commands:edit"},
		CSRFToken:       "csrf",
		ExpiresAt:       time.Now().Add(time.Minute),
	}

	// When
	response := serve(handler.router(), http.MethodGet, "/oauth/consent?attempt=attempt", nil, nil)

	// Then
	if response.Code != http.StatusOK {
		t.Fatalf("consent response = %d %s", response.Code, response.Body.String())
	}
	var body struct {
		RequestedScopes []struct {
			Group       string   `json:"group"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Actions     []string `json:"actions"`
		} `json:"requested_scopes"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	want := []struct {
		group   string
		name    string
		desc    string
		actions []string
	}{
		{group: "commands", name: "Commands", desc: "View and manage custom commands, groups, and role cooldowns", actions: []string{"read", "edit"}},
		{group: "timers", name: "Timers", desc: "View and manage chat timers", actions: []string{"read"}},
		{group: "files", name: "Files", desc: "View and manage uploaded image and audio files", actions: []string{"read", "edit"}},
	}
	if len(body.RequestedScopes) != len(want) {
		t.Fatalf("requested scopes = %#v, want %#v", body.RequestedScopes, want)
	}
	for index, scope := range body.RequestedScopes {
		if scope.Group != want[index].group || scope.Name != want[index].name || scope.Description != want[index].desc || !reflect.DeepEqual(scope.Actions, want[index].actions) {
			t.Fatalf("requested scope[%d] = %#v, want %#v", index, scope, want[index])
		}
	}
}

func TestHandler_consent_get_expands_legacy_read_to_all_groups(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	legacy := testAttempt()
	legacy.RequestedScopes = []string{"read"}
	handler.sessions.attempts["attempt"] = legacy

	// When
	response := serve(handler.router(), http.MethodGet, "/oauth/consent?attempt=attempt", nil, nil)

	// Then
	var body struct {
		RequestedScopes []struct {
			Group   string   `json:"group"`
			Actions []string `json:"actions"`
		} `json:"requested_scopes"`
	}
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &body) != nil {
		t.Fatalf("consent response = %d %s", response.Code, response.Body.String())
	}
	if len(body.RequestedScopes) != 22 {
		t.Fatalf("requested scope count = %d, want 22", len(body.RequestedScopes))
	}
	for _, scope := range body.RequestedScopes {
		if len(scope.Actions) != 1 || scope.Actions[0] != "read" {
			t.Fatalf("legacy requested scope = %#v", scope)
		}
	}
}

type scopeCatalogItemResponse struct {
	Group       string   `json:"group"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
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
	if !reflect.DeepEqual(attempt.RequestedScopes, entity.ScopeStrings(entity.AllScopes())) {
		t.Fatalf("requested scopes = %v, want canonical scopes", attempt.RequestedScopes)
	}
}

func TestHandler_consent_requires_origin_csrf_permission_and_is_one_use(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	attemptID := "attempt"
	handler.sessions.attempts[attemptID] = testAttempt()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["commands:read","commands:edit"]}`, handler.sessions.dashboardID)

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
	createdScopeMatches := handler.service.created.Authorize.Scope == "commands:read commands:edit"
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
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"wrong","decision":"approve","approved_scopes":["commands:read"]}`, handler.sessions.dashboardID)

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
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["commands:read"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body),
		map[string]string{"Content-Type": "application/json", "Origin": "https://twir.example"})

	// Then
	if response.Code != http.StatusForbidden || !strings.Contains(response.Body.String(), "access_denied") {
		t.Fatalf("permission response = %d %s", response.Code, response.Body.String())
	}
}

func TestHandler_consent_approve_issues_only_canonical_approved_subset(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	attempt := testAttempt()
	attempt.RequestedScopes = []string{"timers:read", "commands:edit"}
	handler.sessions.attempts["attempt"] = attempt
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["commands:edit"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusOK || handler.service.created.Authorize.Scope != "commands:read commands:edit" {
		t.Fatalf("approval response = %d %s, scope = %q", response.Code, response.Body.String(), handler.service.created.Authorize.Scope)
	}
}

func TestHandler_consent_approve_rejects_empty_scopes_without_consuming_attempt(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":[]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid_scope") {
		t.Fatalf("approval response = %d %s", response.Code, response.Body.String())
	}
	if _, exists := handler.sessions.attempts["attempt"]; !exists {
		t.Fatal("empty approval consumed authorization attempt")
	}
}

func TestHandler_consent_approve_rejects_out_of_request_scope_without_consuming_attempt(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	attempt := testAttempt()
	attempt.RequestedScopes = []string{"commands:read"}
	handler.sessions.attempts["attempt"] = attempt
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["timers:read"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid_scope") {
		t.Fatalf("approval response = %d %s", response.Code, response.Body.String())
	}
	if _, exists := handler.sessions.attempts["attempt"]; !exists {
		t.Fatal("out-of-request approval consumed authorization attempt")
	}
}

func TestHandler_consent_approve_rejects_read_to_edit_promotion(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	attempt := testAttempt()
	attempt.RequestedScopes = []string{"commands:read"}
	handler.sessions.attempts["attempt"] = attempt
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["commands:edit"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid_scope") {
		t.Fatalf("approval response = %d %s", response.Code, response.Body.String())
	}
}

func TestHandler_consent_deny_rejects_approved_scopes(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	handler.sessions.attempts["attempt"] = testAttempt()
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"deny","approved_scopes":["commands:read"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid_request") {
		t.Fatalf("denial response = %d %s", response.Code, response.Body.String())
	}
	if _, exists := handler.sessions.attempts["attempt"]; !exists {
		t.Fatal("denial with approved scopes consumed authorization attempt")
	}
}

func TestHandler_consent_legacy_approval_is_canonicalized(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	legacy := testAttempt()
	legacy.RequestedScopes = []string{"read", "write"}
	handler.sessions.attempts["attempt"] = legacy
	body := fmt.Sprintf(`{"attempt":"attempt","channel_id":"%s","csrf_token":"csrf","decision":"approve","approved_scopes":["write"]}`, handler.sessions.dashboardID)

	// When
	response := serve(handler.router(), http.MethodPost, "/oauth/consent", strings.NewReader(body), map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://twir.example",
	})

	// Then
	if response.Code != http.StatusOK || handler.service.created.Authorize.Scope != strings.Join(allCanonicalScopeStrings(), " ") {
		t.Fatalf("legacy approval response = %d %s, scope = %q", response.Code, response.Body.String(), handler.service.created.Authorize.Scope)
	}
}

func allCanonicalScopeStrings() []string {
	parsed, _ := entity.ParseScopes("write")
	return entity.ScopeStrings(parsed)
}

func testAttempt() authsessions.MCPOAuthAttempt {
	return authsessions.MCPOAuthAttempt{
		ClientID:        "client",
		RedirectURI:     "https://client.example/callback",
		ClientState:     "client-state",
		CodeChallenge:   "challenge",
		RequestedScopes: []string{"commands:read", "commands:edit"},
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
