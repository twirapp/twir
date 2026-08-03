package mcp_oauth

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

func TestService_RegisterClient_defaults_optional_metadata_and_preserves_clientURI(t *testing.T) {
	// Given
	service := newTestService(t, newFakeRepository(), true, appentity.User{})
	metadata := json.RawMessage(`{"redirect_uris":["example.app://callback"],"client_uri":"https://client.example"}`)

	// When
	client, err := service.RegisterClient(context.Background(), RegisterClientInput{Metadata: metadata})

	// Then
	if err != nil {
		t.Fatal(err)
	}
	if string(client.Metadata) == "" || !equalScopes(client.Scopes, testClient().Scopes) {
		t.Fatalf("metadata = %s, scopes = %v", client.Metadata, client.Scopes)
	}
}

func TestService_ValidateAuthorizeInput_uses_client_scopes_when_scope_omitted(t *testing.T) {
	// Given
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{})

	// When
	request, err := service.ValidateAuthorizeInput(context.Background(), AuthorizeInput{ClientID: client.ClientID, RedirectURI: client.RedirectURIs[0], ResponseType: "code", Resource: "https://twir.example/api/mcp", CodeChallenge: s256Challenge("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"), CodeChallengeMethod: "S256"})

	// Then
	if err != nil || !equalScopes(request.Scopes, client.Scopes) {
		t.Fatalf("request = %#v, err = %v", request, err)
	}
}

func TestService_ValidateAuthorizeInput_normalizes_legacy_client_scopes(t *testing.T) {
	// Given
	repo := newFakeRepository()
	client := testClient()
	client.Scopes = []entity.Scope{entity.ScopeRead}
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{})

	// When
	request, err := service.ValidateAuthorizeInput(context.Background(), AuthorizeInput{ClientID: client.ClientID, RedirectURI: client.RedirectURIs[0], ResponseType: "code", Scope: "commands:read", Resource: "https://twir.example/api/mcp", CodeChallenge: s256Challenge("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"), CodeChallengeMethod: "S256"})

	// Then
	if err != nil {
		t.Fatalf("ValidateAuthorizeInput() error = %v", err)
	}
	if !equalScopes(request.Client.Scopes, canonicalReadScopes()) {
		t.Fatalf("client scopes = %v, want canonical read scopes", request.Client.Scopes)
	}
	if !equalScopes(request.Scopes, []entity.Scope{"commands:read"}) {
		t.Fatalf("requested scopes = %v", request.Scopes)
	}
}

func TestService_ConsentAndTokenOperations_require_resource(t *testing.T) {
	// Given
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, false, appentity.User{})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"

	// When
	_, err := service.CreateAuthorizationCode(context.Background(), CreateAuthorizationCodeInput{Authorize: AuthorizeInput{ClientID: client.ClientID, RedirectURI: client.RedirectURIs[0], ResponseType: "code", Scope: "read", Resource: "https://twir.example/api/mcp", CodeChallenge: s256Challenge(verifier), CodeChallengeMethod: "S256"}, ChannelID: uuid.New(), ApprovingUserID: uuid.New()})

	// Then
	requireOAuthCode(t, err, ErrorAccessDenied)
}
