package mcp_oauth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type testHandler struct {
	*Handler
	service  *fakeService
	sessions *fakeSessions
}

func newTestHandler(t *testing.T) testHandler {
	t.Helper()
	service := &fakeService{client: entity.Client{ClientID: "client", Metadata: []byte(`{"client_name":"Example","client_uri":"https://client.example","redirect_uris":["https://client.example/callback"],"grant_types":["authorization_code","refresh_token"],"response_types":["code"],"token_endpoint_auth_method":"none","scope":"read write"}`), RedirectURIs: []string{"https://client.example/callback"}, Scopes: []entity.Scope{entity.ScopeRead, entity.ScopeWrite}, CreatedAt: time.Unix(1, 0)}, authorized: service.AuthorizationRequest{Client: entity.Client{ClientID: "client", Scopes: []entity.Scope{entity.ScopeRead, entity.ScopeWrite}}, RedirectURI: "https://client.example/callback", CodeChallenge: "challenge", Resource: "https://twir.example/api/mcp", Scopes: []entity.Scope{entity.ScopeRead, entity.ScopeWrite}}, code: service.IssuedAuthorizationCode{Code: "issued-code"}, tokens: service.TokenSet{AccessToken: "access", RefreshToken: "refresh", TokenType: "Bearer", Scopes: []entity.Scope{entity.ScopeRead}, AccessExpiresAt: time.Now().Add(time.Hour)}}
	sessions := &fakeSessions{attempts: map[string]authsessions.MCPOAuthAttempt{}, userID: uuid.New(), dashboardID: uuid.New()}
	handler, err := New(Dependencies{Service: service, Sessions: sessions, SiteBaseURL: "https://twir.example", RegisterRateLimit: func(c *gin.Context) { c.Next() }})
	if err != nil {
		t.Fatal(err)
	}
	return testHandler{Handler: handler, service: service, sessions: sessions}
}

func (handler testHandler) router() *gin.Engine {
	router := gin.New()
	handler.Register(router)
	return router
}

type fakeSessions struct {
	attempts     map[string]authsessions.MCPOAuthAttempt
	userID       uuid.UUID
	dashboardID  uuid.UUID
	userErr      error
	dashboardErr error
}

func (sessions *fakeSessions) SetMCPOAuthAttempt(_ context.Context, id string, attempt authsessions.MCPOAuthAttempt) error {
	sessions.attempts[id] = attempt
	return nil
}
func (sessions *fakeSessions) GetMCPOAuthAttempt(_ context.Context, id string) (authsessions.MCPOAuthAttempt, error) {
	attempt, ok := sessions.attempts[id]
	if !ok {
		return authsessions.MCPOAuthAttempt{}, authsessions.ErrMCPOAuthAttemptNotFound
	}
	return attempt, nil
}
func (sessions *fakeSessions) DeleteMCPOAuthAttempt(_ context.Context, id string) error {
	if _, ok := sessions.attempts[id]; !ok {
		return authsessions.ErrMCPOAuthAttemptNotFound
	}
	delete(sessions.attempts, id)
	return nil
}
func (sessions *fakeSessions) GetInternalUserID(context.Context) (uuid.UUID, error) {
	if sessions.userErr != nil {
		return uuid.Nil, sessions.userErr
	}
	return sessions.userID, nil
}
func (sessions *fakeSessions) GetSelectedDashboard(context.Context) (string, error) {
	if sessions.dashboardErr != nil {
		return "", sessions.dashboardErr
	}
	return sessions.dashboardID.String(), nil
}

type fakeService struct {
	client       entity.Client
	authorized   service.AuthorizationRequest
	code         service.IssuedAuthorizationCode
	tokens       service.TokenSet
	err          error
	getClientErr error
	validateErr  error
	created      service.CreateAuthorizationCodeInput
	exchanged    service.ExchangeAuthorizationCodeInput
	refreshed    service.RefreshInput
	revoked      service.RevokeInput
	revocations  int
}

func (service *fakeService) RegisterClient(context.Context, service.RegisterClientInput) (entity.Client, error) {
	return service.client, service.err
}
func (service *fakeService) GetClient(context.Context, string) (entity.Client, error) {
	return service.client, service.getClientErr
}
func (service *fakeService) ValidateAuthorizeInput(context.Context, service.AuthorizeInput) (service.AuthorizationRequest, error) {
	return service.authorized, service.validateErr
}
func (service *fakeService) CreateAuthorizationCode(_ context.Context, input service.CreateAuthorizationCodeInput) (service.IssuedAuthorizationCode, error) {
	service.created = input
	return service.code, service.err
}
func (service *fakeService) ExchangeAuthorizationCode(_ context.Context, input service.ExchangeAuthorizationCodeInput) (service.TokenSet, error) {
	service.exchanged = input
	return service.tokens, service.err
}
func (service *fakeService) Refresh(_ context.Context, input service.RefreshInput) (service.TokenSet, error) {
	service.refreshed = input
	return service.tokens, service.err
}
func (service *fakeService) Revoke(_ context.Context, input service.RevokeInput) error {
	service.revocations++
	service.revoked = input
	return service.err
}

var _ error = errors.New("")
