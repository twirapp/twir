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
