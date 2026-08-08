package mcp_oauth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	repository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func TestService_RegisterClient_normalizes_public_metadata(t *testing.T) {
	// Given
	repo := newFakeRepository()
	service := newTestService(t, repo, true, appentity.User{})
	metadata := json.RawMessage(`{"client_name":" Example ","redirect_uris":["https://client.example/callback","com.example.app:/oauth/callback"],"grant_types":["refresh_token","authorization_code"],"response_types":["code"],"token_endpoint_auth_method":"none","scope":"write read"}`)

	// When
	client, err := service.RegisterClient(context.Background(), RegisterClientInput{Metadata: metadata})

	// Then
	if err != nil {
		t.Fatalf("RegisterClient() error = %v", err)
	}
	if client.ClientID == "" || len(client.ClientID) != 43 {
		t.Fatalf("client ID = %q, want 32-byte base64url value", client.ClientID)
	}
	if !equalScopes(client.Scopes, entity.AllScopes()) {
		t.Fatalf("scopes = %v", client.Scopes)
	}
	var normalizedMetadata clientMetadata
	if err := json.Unmarshal(client.Metadata, &normalizedMetadata); err != nil {
		t.Fatalf("unmarshal normalized metadata: %v", err)
	}
	if normalizedMetadata.Scope != strings.Join(entity.ScopeStrings(entity.AllScopes()), " ") {
		t.Fatalf("metadata scope = %q, want canonical scopes", normalizedMetadata.Scope)
	}
	if client.RedirectURIs[1] != "com.example.app:/oauth/callback" {
		t.Fatalf("redirect URI = %q", client.RedirectURIs[1])
	}
}

func TestService_RegisterClient_rejects_unknown_or_malformed_scopes(t *testing.T) {
	for _, scope := range []string{"unknown", "commands", "commands:delete"} {
		t.Run(scope, func(t *testing.T) {
			// Given
			service := newTestService(t, newFakeRepository(), true, appentity.User{})
			metadata := json.RawMessage(`{"redirect_uris":["https://client.example/callback"],"scope":"` + scope + `"}`)

			// When
			_, err := service.RegisterClient(context.Background(), RegisterClientInput{Metadata: metadata})

			// Then
			requireOAuthCode(t, err, ErrorInvalidClientMetadata)
		})
	}
}

func TestService_RegisterClient_rejects_unsafe_public_redirects(t *testing.T) {
	for _, redirectURI := range []string{"http://example.com/callback", "https://user@example.com/callback", "https://example.com/callback#fragment", "https://*.example.com/callback"} {
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

func TestService_ExchangeAuthorizationCode_issues_hashed_credentials(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{ID: uuid.NewString()})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	challenge := s256Challenge(verifier)
	issued, err := service.CreateAuthorizationCode(ctx, CreateAuthorizationCodeInput{Authorize: AuthorizeInput{ClientID: client.ClientID, RedirectURI: client.RedirectURIs[0], ResponseType: "code", Scope: "write", Resource: "https://twir.example/api/mcp", CodeChallenge: challenge, CodeChallengeMethod: "S256"}, ChannelID: uuid.New(), ApprovingUserID: uuid.New()})
	if err != nil {
		t.Fatalf("CreateAuthorizationCode() error = %v", err)
	}

	// When
	tokens, err := service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{ClientID: client.ClientID, Code: issued.Code, RedirectURI: client.RedirectURIs[0], CodeVerifier: verifier, Resource: "https://twir.example/api/mcp"})

	// Then
	if err != nil {
		t.Fatalf("ExchangeAuthorizationCode() error = %v", err)
	}
	if len(tokens.AccessToken) != 43 || len(tokens.RefreshToken) != 43 {
		t.Fatalf("opaque token lengths = %d, %d", len(tokens.AccessToken), len(tokens.RefreshToken))
	}
	accessHash := entity.CredentialHash(sha256.Sum256([]byte(tokens.AccessToken)))
	if repo.lastCreated.AccessTokenHash != accessHash {
		t.Fatal("repository did not receive access-token hash")
	}
	if !equalScopes(repo.lastCreated.Scopes, entity.AllScopes()) {
		t.Fatalf("issued scopes = %v", repo.lastCreated.Scopes)
	}
}

func TestService_ExchangeAuthorizationCode_normalizes_legacy_scopes(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	codeValue := "legacy-authorization-code"
	codeHash := credentialHash(codeValue)
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	repo.codes[codeHash] = entity.AuthorizationCode{
		CodeHash:      codeHash,
		ClientID:      client.ClientID,
		ChannelID:     uuid.New(),
		UserID:        uuid.New(),
		RedirectURI:   client.RedirectURIs[0],
		PKCEChallenge: s256Challenge(verifier),
		Scopes:        []entity.Scope{entity.ScopeRead},
		Resource:      "https://twir.example/api/mcp",
		ExpiresAt:     time.Date(2026, 8, 3, 12, 5, 0, 0, time.UTC),
	}
	service := newTestService(t, repo, true, appentity.User{})

	// When
	tokens, err := service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{ClientID: client.ClientID, Code: codeValue, RedirectURI: client.RedirectURIs[0], CodeVerifier: verifier, Resource: "https://twir.example/api/mcp"})

	// Then
	if err != nil {
		t.Fatalf("ExchangeAuthorizationCode() error = %v", err)
	}
	if !equalScopes(tokens.Scopes, canonicalReadScopes()) || !equalScopes(repo.lastCreated.Scopes, canonicalReadScopes()) {
		t.Fatalf("issued scopes = %v, stored scopes = %v", tokens.Scopes, repo.lastCreated.Scopes)
	}
}

func TestService_AuthorizationCode_binds_runtime_loopback_redirect_URI(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	registeredURI := "http://127.0.0.1:3000/callback?source=mcp"
	requestedURI := "http://127.0.0.1:49321/callback?source=mcp"
	client.RedirectURIs = []string{registeredURI}
	repo.clients[client.ClientID] = client
	service := newTestService(t, repo, true, appentity.User{ID: uuid.NewString()})
	verifier := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~abcdefghijk"
	authorize := AuthorizeInput{ClientID: client.ClientID, RedirectURI: requestedURI, ResponseType: "code", Scope: "read", Resource: "https://twir.example/api/mcp", CodeChallenge: s256Challenge(verifier), CodeChallengeMethod: "S256"}
	issueCode := func() (IssuedAuthorizationCode, error) {
		return service.CreateAuthorizationCode(ctx, CreateAuthorizationCodeInput{Authorize: authorize, ChannelID: uuid.New(), ApprovingUserID: uuid.New()})
	}

	// When
	issued, err := issueCode()

	// Then
	if err != nil {
		t.Fatalf("CreateAuthorizationCode() error = %v", err)
	}
	for _, code := range repo.codes {
		if code.RedirectURI != requestedURI {
			t.Fatalf("stored redirect URI = %q, want %q", code.RedirectURI, requestedURI)
		}
	}
	if _, err := service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{ClientID: client.ClientID, Code: issued.Code, RedirectURI: requestedURI, CodeVerifier: verifier, Resource: "https://twir.example/api/mcp"}); err != nil {
		t.Fatalf("ExchangeAuthorizationCode() with runtime redirect URI error = %v", err)
	}

	issued, err = issueCode()
	if err != nil {
		t.Fatalf("CreateAuthorizationCode() second issue error = %v", err)
	}
	if _, err := service.ExchangeAuthorizationCode(ctx, ExchangeAuthorizationCodeInput{ClientID: client.ClientID, Code: issued.Code, RedirectURI: registeredURI, CodeVerifier: verifier, Resource: "https://twir.example/api/mcp"}); err == nil {
		t.Fatal("ExchangeAuthorizationCode() accepted registered URI instead of runtime redirect URI")
	} else {
		requireOAuthCode(t, err, ErrorInvalidGrant)
	}
}

func TestService_Refresh_rejects_scope_elevation_and_revokes_reuse(t *testing.T) {
	for _, test := range []struct {
		name  string
		scope string
	}{
		{name: "legacy write alias", scope: "write"},
		{name: "new group", scope: "timers:read"},
		{name: "read to edit", scope: "commands:edit"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			ctx := context.Background()
			repo := newFakeRepository()
			client := testClient()
			repo.clients[client.ClientID] = client
			refresh := "refresh-token"
			token := testToken(client.ClientID, "access-token", refresh, []entity.Scope{"commands:read"})
			repo.putToken(token)
			service := newTestService(t, repo, true, appentity.User{ID: token.UserID.String()})

			// When
			_, err := service.Refresh(ctx, RefreshInput{ClientID: client.ClientID, RefreshToken: refresh, Scope: test.scope, Resource: "https://twir.example/api/mcp"})

			// Then
			requireOAuthCode(t, err, ErrorInvalidScope)
			if repo.rotations != 0 {
				t.Fatal("scope elevation attempted token rotation")
			}
		})
	}
}

func TestService_Refresh_narrows_scopes_per_group(t *testing.T) {
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
	refreshed, err := service.Refresh(ctx, RefreshInput{ClientID: client.ClientID, RefreshToken: refresh, Scope: "commands:read timers:edit", Resource: "https://twir.example/api/mcp"})

	// Then
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	want := []entity.Scope{"commands:read", "timers:read", "timers:edit"}
	if !equalScopes(refreshed.Scopes, want) || !equalScopes(repo.lastCreated.Scopes, want) {
		t.Fatalf("refreshed scopes = %v, stored scopes = %v, want %v", refreshed.Scopes, repo.lastCreated.Scopes, want)
	}
}

func TestService_VerifyAccessToken_revokes_family_when_permission_is_lost(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	token := testToken(client.ClientID, "access-token", "refresh-token", []entity.Scope{entity.ScopeRead})
	repo.putToken(token)
	service := newTestService(t, repo, false, appentity.User{ID: token.UserID.String()})

	// When
	_, err := service.VerifyAccessToken(ctx, "access-token")

	// Then
	requireOAuthCode(t, err, ErrorInvalidToken)
	if repo.tokensByAccess[token.AccessTokenHash].RevokedAt == nil {
		t.Fatal("permission loss did not revoke token family")
	}
}

func TestService_VerifyAccessToken_normalizes_legacy_scopes(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	token := testToken(client.ClientID, "access-token", "refresh-token", []entity.Scope{entity.ScopeRead})
	repo.putToken(token)
	service := newTestService(t, repo, true, appentity.User{ID: token.UserID.String()})

	// When
	grant, err := service.VerifyAccessToken(ctx, "access-token")

	// Then
	if err != nil {
		t.Fatalf("VerifyAccessToken() error = %v", err)
	}
	if !equalScopes(grant.Scopes, canonicalReadScopes()) {
		t.Fatalf("grant scopes = %v, want canonical read scopes", grant.Scopes)
	}
}

func TestService_VerifyAccessToken_rejects_invalid_stored_scopes(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := newFakeRepository()
	client := testClient()
	repo.clients[client.ClientID] = client
	token := testToken(client.ClientID, "access-token", "refresh-token", []entity.Scope{"commands:delete"})
	repo.putToken(token)
	service := newTestService(t, repo, true, appentity.User{ID: token.UserID.String()})

	// When
	_, err := service.VerifyAccessToken(ctx, "access-token")

	// Then
	requireOAuthCode(t, err, ErrorInvalidToken)
}

func newTestService(t *testing.T, repo *fakeRepository, allowed bool, user appentity.User) *Service {
	t.Helper()
	channelID := uuid.New()
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	service, err := New(Dependencies{Repository: repo, Users: fakeUsers{user: user}, Channels: fakeChannels{channel: channelentity.Channel{ID: channelID}}, DashboardAccess: fakeDashboardAccess{allowed: allowed}, SiteBaseURL: "https://twir.example/path", Clock: fixedClock{now: time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)}, Random: bytes.NewReader(bytes.Repeat([]byte{1, 2, 3, 4}, 64))})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return service
}

func testClient() entity.Client {
	return entity.Client{ClientID: "client", RedirectURIs: []string{"https://client.example/callback"}, GrantTypes: []string{"authorization_code", "refresh_token"}, ResponseTypes: []string{"code"}, TokenEndpointAuthMethod: "none", Scopes: entity.AllScopes()}
}

func testToken(clientID, access, refresh string, scopes []entity.Scope) entity.Token {
	accessHash := entity.CredentialHash(sha256.Sum256([]byte(access)))
	refreshHash := entity.CredentialHash(sha256.Sum256([]byte(refresh)))
	return entity.Token{ID: uuid.New(), FamilyID: uuid.New(), ClientID: clientID, ChannelID: uuid.New(), UserID: uuid.New(), AccessTokenHash: accessHash, RefreshTokenHash: refreshHash, Scopes: scopes, Resource: "https://twir.example/api/mcp", AccessExpiresAt: time.Date(2026, 8, 3, 13, 0, 0, 0, time.UTC), RefreshExpiresAt: time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)}
}

func equalScopes(left, right []entity.Scope) bool {
	return slices.Equal(left, right)
}

func canonicalReadScopes() []entity.Scope {
	return slices.DeleteFunc(entity.AllScopes(), func(scope entity.Scope) bool {
		return !strings.HasSuffix(string(scope), ":read")
	})
}

func requireOAuthCode(t *testing.T, err error, want ErrorCode) {
	t.Helper()
	var oauthError *OAuthError
	if !errors.As(err, &oauthError) || oauthError.Code != want {
		t.Fatalf("error = %v, want OAuth %s", err, want)
	}
}

type fixedClock struct{ now time.Time }

func (clock fixedClock) Now() time.Time { return clock.now }

type fakeUsers struct{ user appentity.User }

func (users fakeUsers) GetByID(context.Context, string) (appentity.User, error) {
	return users.user, nil
}

type fakeChannels struct{ channel channelentity.Channel }

func (channels fakeChannels) GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return channels.channel, nil
}

type fakeDashboardAccess struct{ allowed bool }

func (access fakeDashboardAccess) CanAccess(context.Context, DashboardSubject, uuid.UUID, string) (bool, error) {
	return access.allowed, nil
}

var _ repository.Repository = (*fakeRepository)(nil)

type fakeRepository struct {
	clients         map[string]entity.Client
	codes           map[entity.CredentialHash]entity.AuthorizationCode
	tokensByAccess  map[entity.CredentialHash]entity.Token
	tokensByRefresh map[entity.CredentialHash]entity.Token
	lastCreated     repository.CreateTokenInput
	rotations       int
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{clients: map[string]entity.Client{}, codes: map[entity.CredentialHash]entity.AuthorizationCode{}, tokensByAccess: map[entity.CredentialHash]entity.Token{}, tokensByRefresh: map[entity.CredentialHash]entity.Token{}}
}
func (repo *fakeRepository) CreateClient(_ context.Context, input repository.CreateClientInput) (entity.Client, error) {
	client := entity.Client{ClientID: input.ClientID, Metadata: input.Metadata, RedirectURIs: input.RedirectURIs, GrantTypes: input.GrantTypes, ResponseTypes: input.ResponseTypes, TokenEndpointAuthMethod: input.TokenEndpointAuthMethod, Scopes: input.Scopes}
	repo.clients[client.ClientID] = client
	return client, nil
}
func (repo *fakeRepository) GetClient(_ context.Context, id string) (entity.Client, error) {
	client, ok := repo.clients[id]
	if !ok {
		return entity.NilClient, repository.ErrClientNotFound
	}
	return client, nil
}
func (repo *fakeRepository) CreateAuthorizationCode(_ context.Context, input repository.CreateAuthorizationCodeInput) (entity.AuthorizationCode, error) {
	code := entity.AuthorizationCode{CodeHash: input.CodeHash, ClientID: input.ClientID, ChannelID: input.ChannelID, UserID: input.UserID, RedirectURI: input.RedirectURI, PKCEChallenge: input.PKCEChallenge, Scopes: input.Scopes, Resource: input.Resource, ExpiresAt: input.ExpiresAt}
	repo.codes[code.CodeHash] = code
	return code, nil
}
func (repo *fakeRepository) ConsumeAuthorizationCode(_ context.Context, hash entity.CredentialHash) (entity.AuthorizationCode, error) {
	code, ok := repo.codes[hash]
	if !ok {
		return entity.NilAuthorizationCode, repository.ErrAuthorizationCodeNotFound
	}
	delete(repo.codes, hash)
	return code, nil
}
func (repo *fakeRepository) CreateToken(_ context.Context, input repository.CreateTokenInput) (entity.Token, error) {
	repo.lastCreated = input
	token := entity.Token{ID: uuid.New(), FamilyID: uuid.New(), ClientID: input.ClientID, ChannelID: input.ChannelID, UserID: input.UserID, AccessTokenHash: input.AccessTokenHash, RefreshTokenHash: input.RefreshTokenHash, Scopes: input.Scopes, Resource: input.Resource, AccessExpiresAt: input.AccessExpiresAt, RefreshExpiresAt: input.RefreshExpiresAt}
	repo.putToken(token)
	return token, nil
}
func (repo *fakeRepository) GetTokenByAccessTokenHash(_ context.Context, hash entity.CredentialHash) (entity.Token, error) {
	token, ok := repo.tokensByAccess[hash]
	if !ok {
		return entity.NilToken, repository.ErrAccessTokenNotFound
	}
	return token, nil
}
func (repo *fakeRepository) GetTokenByRefreshTokenHash(_ context.Context, hash entity.CredentialHash) (entity.Token, error) {
	token, ok := repo.tokensByRefresh[hash]
	if !ok {
		return entity.NilToken, repository.ErrRefreshTokenNotFound
	}
	return token, nil
}
func (repo *fakeRepository) RotateRefreshToken(ctx context.Context, input repository.RotateRefreshTokenInput) (entity.Token, error) {
	repo.rotations++
	current, err := repo.GetTokenByRefreshTokenHash(ctx, input.PresentedRefreshTokenHash)
	if err != nil {
		return entity.NilToken, err
	}
	if current.RevokedAt != nil || current.ReplacedByID != nil {
		_ = repo.RevokeToken(ctx, current.ClientID, input.PresentedRefreshTokenHash)
		return entity.NilToken, &repository.RefreshTokenReuseError{FamilyID: current.FamilyID}
	}
	next, err := repo.CreateToken(ctx, repository.CreateTokenInput{ClientID: current.ClientID, ChannelID: current.ChannelID, UserID: current.UserID, AccessTokenHash: input.NextAccessTokenHash, RefreshTokenHash: input.NextRefreshTokenHash, Scopes: input.Scopes, Resource: current.Resource, AccessExpiresAt: input.AccessExpiresAt, RefreshExpiresAt: input.RefreshExpiresAt})
	if err != nil {
		return entity.NilToken, err
	}
	current.ReplacedByID = &next.ID
	repo.putToken(current)
	return next, nil
}
func (repo *fakeRepository) RevokeToken(_ context.Context, clientID string, hash entity.CredentialHash) error {
	token, ok := repo.tokensByAccess[hash]
	if !ok {
		token, ok = repo.tokensByRefresh[hash]
	}
	if !ok {
		return nil
	}
	if clientID != "" && token.ClientID != clientID {
		return nil
	}
	now := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	for _, candidate := range repo.tokensByAccess {
		if candidate.FamilyID == token.FamilyID {
			candidate.RevokedAt = &now
			repo.putToken(candidate)
		}
	}
	return nil
}
func (repo *fakeRepository) putToken(token entity.Token) {
	repo.tokensByAccess[token.AccessTokenHash] = token
	repo.tokensByRefresh[token.RefreshTokenHash] = token
}
