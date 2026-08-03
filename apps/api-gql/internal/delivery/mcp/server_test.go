package mcp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appserver "github.com/twirapp/twir/apps/api-gql/internal/server"
	oauth "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

const protectedResourceMetadataURL = "https://twir.example/.well-known/oauth-protected-resource/api/mcp"
const bearerChallenge = `Bearer resource_metadata="` + protectedResourceMetadataURL + `", scope="read write"`
const corsAllowMethods = "POST, GET, OPTIONS"
const corsAllowHeaders = "Authorization, Content-Type, Accept, MCP-Protocol-Version, Last-Event-ID"
const corsExposeHeaders = "MCP-Session-Id, MCP-Protocol-Version, WWW-Authenticate"

func TestNewServerRegistersToolSchemas(t *testing.T) {
	handler := &Handler{}
	server := handler.newServer(scope{Channel: channelentity.Channel{ID: uuid.New()}})
	if server == nil {
		t.Fatal("expected MCP server")
	}
}

func TestRegisterServesOnlyInternalMCP(t *testing.T) {
	verifier := &mcpAccessVerifier{}
	handler := newTestMCPHandler(verifier)
	router := &appserver.Server{Engine: gin.New()}
	Register(router, handler)

	internal := httptest.NewRecorder()
	router.ServeHTTP(internal, httptest.NewRequest(http.MethodPost, "/mcp", nil))
	public := httptest.NewRecorder()
	router.ServeHTTP(public, httptest.NewRequest(http.MethodPost, "/api/mcp", nil))

	if internal.Code != http.StatusUnauthorized || public.Code != http.StatusNotFound {
		t.Fatalf("route statuses = %d, %d", internal.Code, public.Code)
	}
}

func TestHandlerAllowsCORSPreflightWithoutBearer(t *testing.T) {
	verifier := &mcpAccessVerifier{}
	handler := newTestMCPHandler(verifier)
	request := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
	request.Header.Set("Origin", "https://client.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent || verifier.calls != 0 {
		t.Fatalf("response = %d, verifier calls = %d", response.Code, verifier.calls)
	}
	assertMCPResponseCORS(t, response.Header())
}

func TestHandlerRejectsNonBearerCredentials(t *testing.T) {
	tests := []struct {
		name           string
		authorizations []string
		apiKey         string
		wantCalls      int
	}{
		{name: "missing bearer"},
		{name: "API key alone", apiKey: "legacy-api-key"},
		{name: "other scheme", authorizations: []string{"Basic token"}},
		{name: "empty bearer", authorizations: []string{"Bearer "}},
		{name: "multiple credentials", authorizations: []string{"Bearer first Bearer second"}},
		{name: "multiple authorization headers", authorizations: []string{"Bearer first", "Bearer second"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &mcpAccessVerifier{}
			handler := newTestMCPHandler(verifier)
			request := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
			for _, authorization := range test.authorizations {
				request.Header.Add("Authorization", authorization)
			}
			request.Header.Set("Api-Key", test.apiKey)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			assertUnauthorized(t, response, test.wantCalls, verifier)
		})
	}
}

func TestHandlerDoesNotLeakAccessVerificationFailures(t *testing.T) {
	for _, name := range []string{"invalid", "expired", "revoked", "wrong resource", "permission lost"} {
		t.Run(name, func(t *testing.T) {
			verifier := &mcpAccessVerifier{err: errors.New(name + " access token")}
			handler := newTestMCPHandler(verifier)
			request := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
			request.Header.Set("Authorization", "Bearer opaque-token")
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			assertUnauthorized(t, response, 1, verifier)
			if strings.Contains(response.Body.String(), verifier.err.Error()) {
				t.Fatalf("response leaked verifier error: %s", response.Body.String())
			}
		})
	}
}

func TestHandlerBindsVerifiedGrantToRequestScope(t *testing.T) {
	grant := testAuthorizedGrant(entity.ScopeRead)
	verifier := &mcpAccessVerifier{grant: grant}
	handler := newTestMCPHandler(verifier)
	handler.transport = http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		requestScope, ok := request.Context().Value(contextKey{}).(scope)
		if !ok {
			t.Fatal("request scope was not set")
		}
		if requestScope.Channel.ID != grant.Channel.ID {
			t.Fatalf("channel ID = %s, want %s", requestScope.Channel.ID, grant.Channel.ID)
		}
		if requestScope.ActorID != grant.ApprovingUserID.String() {
			t.Fatalf("actor ID = %s, want %s", requestScope.ActorID, grant.ApprovingUserID)
		}
		if !requestScope.AccessScopes.allowsTool("list_commands") || requestScope.AccessScopes.allowsTool("create_command") {
			t.Fatalf("access scopes = %#v", requestScope.AccessScopes)
		}
	})
	request := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
	request.Header.Set("Authorization", "Bearer opaque-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK || verifier.token != "opaque-token" {
		t.Fatalf("response = %d, verified token = %q", response.Code, verifier.token)
	}
	assertMCPResponseCORS(t, response.Header())
}

func TestHandlerRejectsInvalidGrantScopes(t *testing.T) {
	for _, scopes := range [][]entity.Scope{{entity.ScopeWrite}, {entity.ScopeRead, "unknown"}} {
		t.Run(strings.Join(scopeNames(scopes), ","), func(t *testing.T) {
			verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(scopes...)}
			handler := newTestMCPHandler(verifier)
			request := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
			request.Header.Set("Authorization", "Bearer opaque-token")
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			assertUnauthorized(t, response, 1, verifier)
		})
	}
}

func TestHandlerReadGrantKeepsStreamableHTTPStatelessAndRejectsWrites(t *testing.T) {
	verifier := &mcpAccessVerifier{grant: testAuthorizedGrant(entity.ScopeRead)}
	handler := newTestMCPHandler(verifier)

	list := serveMCPRequest(handler, `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)
	if list.Code != http.StatusOK || list.Header().Get("Mcp-Session-Id") != "" {
		t.Fatalf("list response = %d, session ID = %q", list.Code, list.Header().Get("Mcp-Session-Id"))
	}
	if !strings.Contains(list.Body.String(), "list_commands") || !strings.Contains(list.Body.String(), "get_command") || strings.Contains(list.Body.String(), "get_secret") || strings.Contains(list.Body.String(), "create_command") {
		t.Fatalf("unexpected read tool surface: %s", list.Body.String())
	}

	secret := serveMCPRequest(handler, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_secret","arguments":{"id":"00000000-0000-0000-0000-000000000000"}}}`)
	if secret.Code != http.StatusOK || !strings.Contains(secret.Body.String(), "unknown tool") {
		t.Fatalf("secret response = %d %s", secret.Code, secret.Body.String())
	}

	write := serveMCPRequest(handler, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"create_command","arguments":{}}}`)
	if write.Code != http.StatusOK || !strings.Contains(write.Body.String(), "unknown tool") {
		t.Fatalf("write response = %d %s", write.Code, write.Body.String())
	}

	get := httptest.NewRequest(http.MethodGet, "/api/mcp", nil)
	get.Header.Set("Authorization", "Bearer opaque-token")
	get.Header.Set("Accept", "application/json, text/event-stream")
	getResponse := httptest.NewRecorder()
	handler.ServeHTTP(getResponse, get)
	if getResponse.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET response = %d, want %d", getResponse.Code, http.StatusMethodNotAllowed)
	}
}

func newTestMCPHandler(verifier *mcpAccessVerifier) *Handler {
	return New(Deps{AccessTokenVerifier: verifier})
}

func testAuthorizedGrant(scopes ...entity.Scope) oauth.AuthorizedGrant {
	return oauth.AuthorizedGrant{Channel: channelentity.Channel{ID: uuid.New()}, ApprovingUserID: uuid.New(), Scopes: scopes}
}

func serveMCPRequest(handler *Handler, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/mcp", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer opaque-token")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertUnauthorized(t *testing.T, response *httptest.ResponseRecorder, wantCalls int, verifier *mcpAccessVerifier) {
	t.Helper()
	if response.Code != http.StatusUnauthorized || response.Header().Get("WWW-Authenticate") != bearerChallenge {
		t.Fatalf("response = %d, challenge = %q", response.Code, response.Header().Get("WWW-Authenticate"))
	}
	assertMCPResponseCORS(t, response.Header())
	if verifier.calls != wantCalls || response.Body.String() != "unauthorized\n" {
		t.Fatalf("verifier calls = %d, response = %q", verifier.calls, response.Body.String())
	}
}

func assertMCPResponseCORS(t *testing.T, headers http.Header) {
	t.Helper()
	if headers.Get("Access-Control-Allow-Origin") != "*" ||
		headers.Get("Access-Control-Allow-Methods") != corsAllowMethods ||
		headers.Get("Access-Control-Allow-Headers") != corsAllowHeaders ||
		headers.Get("Access-Control-Expose-Headers") != corsExposeHeaders {
		t.Fatalf("CORS headers = %#v", headers)
	}
}

func scopeNames(scopes []entity.Scope) []string {
	names := make([]string, len(scopes))
	for index, scope := range scopes {
		names[index] = string(scope)
	}
	return names
}

type mcpAccessVerifier struct {
	grant oauth.AuthorizedGrant
	err   error
	token string
	calls int
}

func (verifier *mcpAccessVerifier) VerifyAccessToken(_ context.Context, token string) (oauth.AuthorizedGrant, error) {
	verifier.calls++
	verifier.token = token
	return verifier.grant, verifier.err
}

func (*mcpAccessVerifier) ProtectedResourceMetadataURL() string {
	return protectedResourceMetadataURL
}
