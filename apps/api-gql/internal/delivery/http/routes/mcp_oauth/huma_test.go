package mcp_oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
)

func TestHandler_registers_exact_public_huma_routes_without_security(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	wantRoutes := map[string]struct{}{
		http.MethodGet + " /.well-known/oauth-protected-resource":   {},
		http.MethodGet + " /.well-known/oauth-authorization-server": {},
		http.MethodOptions + " /oauth/register":                     {},
		http.MethodPost + " /oauth/register":                        {},
		http.MethodGet + " /oauth/authorize":                        {},
		http.MethodGet + " /oauth/consent":                          {},
		http.MethodPost + " /oauth/consent":                         {},
		http.MethodOptions + " /oauth/token":                        {},
		http.MethodPost + " /oauth/token":                           {},
		http.MethodOptions + " /oauth/revoke":                       {},
		http.MethodPost + " /oauth/revoke":                          {},
	}

	// When
	router := handler.router()
	for _, route := range handler.routes() {
		meta := route.GetMeta()
		if _, ok := wantRoutes[meta.Method+" "+meta.Path]; !ok {
			t.Fatalf("unexpected route = %s %s", meta.Method, meta.Path)
		}
		if len(meta.Security) != 0 {
			t.Fatalf("%s %s security = %#v", meta.Method, meta.Path, meta.Security)
		}
	}

	// Then
	registered := make(map[string]struct{})
	for _, route := range router.Routes() {
		registered[route.Method+" "+route.Path] = struct{}{}
	}
	for wantRoute := range wantRoutes {
		if _, ok := registered[wantRoute]; !ok {
			t.Fatalf("missing Gin-Huma route = %s", wantRoute)
		}
	}
}

func TestHandler_register_preserves_json_transport_failures_and_cors(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		status      int
	}{
		{name: "accepts normalized JSON media type", body: `{"client_name":"Example"}`, contentType: "application/json; charset=utf-8", status: http.StatusCreated},
		{name: "rejects malformed JSON", body: `{"client_name":`, contentType: "application/json", status: http.StatusBadRequest},
		{name: "rejects trailing JSON", body: `{"client_name":"Example"} {}`, contentType: "application/json", status: http.StatusBadRequest},
		{name: "rejects JSON over 16KiB", body: `{"client_name":"` + strings.Repeat("x", requestBodyLimit) + `"}`, contentType: "application/json", status: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			handler := newTestHandler(t)

			// When
			response := serve(handler.router(), http.MethodPost, "/oauth/register", strings.NewReader(test.body), map[string]string{"Content-Type": test.contentType})

			// Then
			if response.Code != test.status || response.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Fatalf("register response = %d, headers = %#v", response.Code, response.Header())
			}
			if test.status == http.StatusBadRequest {
				var errorResponse oauthErrorResponse
				if err := json.Unmarshal(response.Body.Bytes(), &errorResponse); err != nil || errorResponse.Code == "" || errorResponse.Description == "" {
					t.Fatalf("error response = %s, err = %v", response.Body.String(), err)
				}
			}
		})
	}
}

func TestHandler_token_preserves_form_auth_rejection_and_cache_headers(t *testing.T) {
	tests := []struct {
		name          string
		form          url.Values
		headers       map[string]string
		status        int
		wantChallenge bool
	}{
		{name: "accepts normalized form media type", form: tokenForm(), headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"}, status: http.StatusOK},
		{name: "rejects client secret without challenge", form: addClientSecret(tokenForm()), headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, status: http.StatusUnauthorized},
		{name: "rejects authorization with challenge", form: tokenForm(), headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Authorization": "Basic Y2xpZW50OnNlY3JldA=="}, status: http.StatusUnauthorized, wantChallenge: true},
		{name: "rejects form over 16KiB", form: url.Values{"client_id": {strings.Repeat("x", requestBodyLimit+1)}}, headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, status: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			handler := newTestHandler(t)

			// When
			response := serve(handler.router(), http.MethodPost, "/oauth/token", strings.NewReader(test.form.Encode()), test.headers)

			// Then
			if response.Code != test.status || response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Pragma") != "no-cache" || response.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Fatalf("token response = %d, headers = %#v", response.Code, response.Header())
			}
			if test.wantChallenge && response.Header().Get("WWW-Authenticate") != `Basic realm="oauth"` {
				t.Fatalf("WWW-Authenticate = %q", response.Header().Get("WWW-Authenticate"))
			}
			if !test.wantChallenge && response.Header().Get("WWW-Authenticate") != "" {
				t.Fatalf("unexpected WWW-Authenticate = %q", response.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestHandler_register_limiter_preserves_gin_client_ip_and_short_circuits(t *testing.T) {
	// Given
	responses := make([]rate_limiter.LeakyResponse, 21)
	for index := range responses[:20] {
		responses[index] = rate_limiter.LeakyResponse{Success: true, RemainingTokens: 20 - index, ResetAt: time.Unix(100, 0)}
	}
	responses[20] = rate_limiter.LeakyResponse{Success: false, ResetAt: time.Unix(100, 0)}
	limiter := &fakeRegisterRateLimiter{responses: responses}
	handler := newTestHandlerWithRateLimiter(t, limiter)
	body := `{"client_name":"Example"}`

	// When
	for range 20 {
		request := httptest.NewRequest(http.MethodPost, "/oauth/register", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.RemoteAddr = "198.51.100.11:1234"
		response := httptest.NewRecorder()
		handler.router().ServeHTTP(response, request)
		if response.Code != http.StatusCreated {
			t.Fatalf("successful register status = %d", response.Code)
		}
	}
	request := httptest.NewRequest(http.MethodPost, "/oauth/register", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "198.51.100.11:1234"
	response := httptest.NewRecorder()
	handler.router().ServeHTTP(response, request)

	// Then
	if response.Code != http.StatusTooManyRequests || response.Body.Len() != 0 || response.Header().Get("Access-Control-Allow-Origin") != "*" || response.Header().Get("X-Rate-Limit-Limit") != "20" {
		t.Fatalf("limited response = %d, headers = %#v, body = %q", response.Code, response.Header(), response.Body.String())
	}
	if handler.service.registrations != 20 || len(limiter.options) != 21 {
		t.Fatalf("registrations = %d, limiter calls = %d", handler.service.registrations, len(limiter.options))
	}
	options := limiter.options[0]
	if options.KeyPrefix != "api-gql:ratelimiter:198.51.100.11:mcp-oauth-register" || options.MaximumCapacity != 20 || options.WindowSeconds != 60 {
		t.Fatalf("limiter options = %#v", options)
	}
}

func tokenForm() url.Values {
	return url.Values{"grant_type": {"authorization_code"}, "client_id": {"client"}, "code": {"one-use-code"}, "redirect_uri": {"https://client.example/callback"}, "code_verifier": {"verifier"}, "resource": {"https://twir.example/api/mcp"}}
}

func addClientSecret(form url.Values) url.Values {
	form.Set("client_secret", "forbidden")
	return form
}
