package mcp_oauth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

const testMCPResource = "https://twir.example/api/mcp"

func TestService_ValidateAuthorizeInput_defaults_omitted_resource_and_rejects_mismatch(t *testing.T) {
	// Given
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"

	for _, test := range []struct {
		name     string
		resource string
		wantCode ErrorCode
	}{
		{name: "omitted resource"},
		{name: "explicit mismatch", resource: "https://other.example/mcp", wantCode: ErrorInvalidRequest},
		{name: "whitespace is an explicit mismatch", resource: " ", wantCode: ErrorInvalidRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			// When
			request, err := service.ValidateAuthorizeInput(context.Background(), AuthorizeInput{
				ClientID:            client.ClientID,
				RedirectURI:         client.RedirectURIs[0],
				ResponseType:        "code",
				Resource:            test.resource,
				CodeChallenge:       s256Challenge(verifier),
				CodeChallengeMethod: "S256",
			})

			// Then
			if test.wantCode != "" {
				requireOAuthCode(t, err, test.wantCode)
				return
			}
			if err != nil || request.Resource != testMCPResource {
				t.Fatalf("request = %#v, err = %v", request, err)
			}
		})
	}
}

func TestService_AuthorizationCodeAndExchange_default_omitted_resource_and_persist_canonical_resource(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{ID: uuid.NewString()})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	authorize := AuthorizeInput{
		ClientID:            client.ClientID,
		RedirectURI:         client.RedirectURIs[0],
		ResponseType:        "code",
		Resource:            "",
		CodeChallenge:       s256Challenge(verifier),
		CodeChallengeMethod: "S256",
	}

	// When
	issued, err := service.CreateAuthorizationCode(ctx, CreateAuthorizationCodeInput{
		Authorize:       authorize,
		ChannelID:       uuid.New(),
		ApprovingUserID: uuid.New(),
	})

	// Then
	if err != nil {
		t.Fatalf("CreateAuthorizationCode() error = %v", err)
	}
	if len(repo.codes) != 1 {
		t.Fatalf("stored authorization codes = %d, want 1", len(repo.codes))
	}
	for _, code := range repo.codes {
		if code.Resource != testMCPResource {
			t.Fatalf("stored authorization code resource = %q, want %q", code.Resource, testMCPResource)
		}
	}

	tokens, err := service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{
		ClientID:     client.ClientID,
		Code:         issued.Code,
		RedirectURI:  client.RedirectURIs[0],
		CodeVerifier: verifier,
		Resource:     "",
	})
	if err != nil {
		t.Fatalf("ExchangeAuthorizationCode() error = %v", err)
	}
	if tokens.Resource != testMCPResource || repo.lastCreated.Resource != testMCPResource {
		t.Fatalf("token resources = %q, stored = %q, want %q", tokens.Resource, repo.lastCreated.Resource, testMCPResource)
	}
}

func TestService_ExchangeAuthorizationCode_rejects_explicit_resource_mismatch(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{ID: uuid.NewString()})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	issued, err := service.CreateAuthorizationCode(ctx, CreateAuthorizationCodeInput{
		Authorize: AuthorizeInput{
			ClientID:            client.ClientID,
			RedirectURI:         client.RedirectURIs[0],
			ResponseType:        "code",
			Resource:            testMCPResource,
			CodeChallenge:       s256Challenge(verifier),
			CodeChallengeMethod: "S256",
		},
		ChannelID:       uuid.New(),
		ApprovingUserID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("CreateAuthorizationCode() error = %v", err)
	}

	// When
	_, err = service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{
		ClientID:     client.ClientID,
		Code:         issued.Code,
		RedirectURI:  client.RedirectURIs[0],
		CodeVerifier: verifier,
		Resource:     "https://other.example/mcp",
	})

	// Then
	requireOAuthCode(t, err, ErrorInvalidGrant)
}

func TestService_Refresh_defaults_omitted_resource_and_persists_canonical_resource(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	refresh := "refresh-token"
	token := testToken(client.ClientID, "access-token", refresh, entity.AllScopes())
	repo.putToken(token)
	service := newTestService(t, repo, true, appentity.User{ID: token.UserID.String()})

	// When
	refreshed, err := service.Refresh(ctx, RefreshInput{ClientID: client.ClientID, RefreshToken: refresh})

	// Then
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if refreshed.Resource != testMCPResource || repo.lastCreated.Resource != testMCPResource {
		t.Fatalf("token resources = %q, stored = %q, want %q", refreshed.Resource, repo.lastCreated.Resource, testMCPResource)
	}
}

func TestService_Refresh_rejects_explicit_resource_mismatch(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	refresh := "refresh-token"
	token := testToken(client.ClientID, "access-token", refresh, entity.AllScopes())
	repo.putToken(token)
	service := newTestService(t, repo, true, appentity.User{ID: token.UserID.String()})

	// When
	_, err := service.Refresh(ctx, RefreshInput{ClientID: client.ClientID, RefreshToken: refresh, Resource: "https://other.example/mcp"})

	// Then
	requireOAuthCode(t, err, ErrorInvalidGrant)
	if repo.rotations != 0 {
		t.Fatalf("refresh rotations = %d, want 0", repo.rotations)
	}
}
