package mcp_oauth

import (
	"context"
	"encoding/json"
	"testing"

	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func TestService_RegisterClient_rejects_dangerous_redirect_schemes(t *testing.T) {
	for _, redirectURI := range []string{"javascript:alert(1)", "data:text/html,alert(1)", "file:///etc/passwd", "vbscript:msgbox(1)", "blob:https://client.example/callback", "about:blank"} {
		t.Run(redirectURI, func(t *testing.T) {
			// Given
			service := newTestService(t, newFakeRepository(), true, appentity.User{})
			metadata := json.RawMessage(`{"redirect_uris":["` + redirectURI + `"],"grant_types":["authorization_code","refresh_token"],"response_types":["code"],"token_endpoint_auth_method":"none","scope":"read"}`)

			// When
			_, err := service.RegisterClient(context.Background(), RegisterClientInput{Metadata: metadata})

			// Then
			requireOAuthCode(t, err, ErrorInvalidClientMetadata)
		})
	}
}

func TestService_RegisterClient_allows_reverse_dns_private_use_redirect_scheme(t *testing.T) {
	// Given
	service := newTestService(t, newFakeRepository(), true, appentity.User{})
	metadata := json.RawMessage(`{"redirect_uris":["com.example.app:/oauth/callback"],"grant_types":["authorization_code","refresh_token"],"response_types":["code"],"token_endpoint_auth_method":"none","scope":"read"}`)

	// When
	client, err := service.RegisterClient(context.Background(), RegisterClientInput{Metadata: metadata})

	// Then
	if err != nil {
		t.Fatalf("RegisterClient() error = %v", err)
	}
	if client.RedirectURIs[0] != "com.example.app:/oauth/callback" {
		t.Fatalf("redirect URI = %q, want %q", client.RedirectURIs[0], "com.example.app:/oauth/callback")
	}
}

func TestService_ValidateAuthorizeInput_matches_loopback_redirect_with_dynamic_port(t *testing.T) {
	const challenge = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	tests := []struct {
		name          string
		registeredURI string
		requestedURI  string
		wantError     bool
	}{
		{name: "allows changed port", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/callback?source=mcp"},
		{name: "allows localhost changed port", registeredURI: "http://localhost:3000/callback?source=mcp", requestedURI: "http://localhost:49321/callback?source=mcp"},
		{name: "allows IPv6 loopback changed port", registeredURI: "http://[::1]:3000/callback?source=mcp", requestedURI: "http://[::1]:49321/callback?source=mcp"},
		{name: "rejects changed loopback host", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.2:49321/callback?source=mcp", wantError: true},
		{name: "rejects changed path", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/other?source=mcp", wantError: true},
		{name: "rejects changed query", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/callback?source=other", wantError: true},
		{name: "rejects user info", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://client@127.0.0.1:49321/callback?source=mcp", wantError: true},
		{name: "rejects fragment", registeredURI: "http://127.0.0.1:3000/callback?source=mcp", requestedURI: "http://127.0.0.1:49321/callback?source=mcp#fragment", wantError: true},
		{name: "rejects changed HTTPS port", registeredURI: "https://client.example:3000/callback", requestedURI: "https://client.example:49321/callback", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			repo := newFakeRepository()
			client := testClient()
			client.RedirectURIs = []string{test.registeredURI}
			repo.clients[client.ClientID] = client
			service := newTestService(t, repo, true, appentity.User{})

			// When
			request, err := service.ValidateAuthorizeInput(context.Background(), AuthorizeInput{ClientID: client.ClientID, RedirectURI: test.requestedURI, ResponseType: "code", Resource: "https://twir.example/api/mcp", CodeChallenge: s256Challenge(challenge), CodeChallengeMethod: "S256"})

			// Then
			if test.wantError {
				requireOAuthCode(t, err, ErrorInvalidRequest)
				return
			}
			if err != nil || request.RedirectURI != test.requestedURI {
				t.Fatalf("ValidateAuthorizeInput() request = %#v, error = %v", request, err)
			}
		})
	}
}
