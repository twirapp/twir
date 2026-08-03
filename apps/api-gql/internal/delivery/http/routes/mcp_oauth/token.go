package mcp_oauth

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type tokenOptions struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *preflightOutput] = (*tokenOptions)(nil)

func newTokenOptions(handler *Handler) *tokenOptions {
	return &tokenOptions{handler: handler}
}

func (*tokenOptions) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-token-options",
		Method:        http.MethodOptions,
		Path:          "/oauth/token",
		DefaultStatus: http.StatusNoContent,
	}
}

func (*tokenOptions) Handler(context.Context, *emptyInput) (*preflightOutput, error) {
	return &preflightOutput{Status: http.StatusNoContent}, nil
}

func (route *tokenOptions) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Middlewares = huma.Middlewares{route.handler.publicPreflight}
	huma.Register(api, meta, route.Handler)
}

type tokenRoute struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *huma.StreamResponse] = (*tokenRoute)(nil)

func newToken(handler *Handler) *tokenRoute {
	return &tokenRoute{handler: handler}
}

func (*tokenRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-token",
		Method:        http.MethodPost,
		Path:          "/oauth/token",
		Tags:          []string{"MCP OAuth"},
		Summary:       "Issue OAuth tokens",
		DefaultStatus: http.StatusOK,
	}
}

func (route *tokenRoute) Handler(ctx context.Context, _ *emptyInput) (*huma.StreamResponse, error) {
	form, ok := parsedBody[url.Values](ctx, tokenBodyKey)
	if !ok {
		return rawStream(func(context huma.Context) {
			writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		}), nil
	}
	return rawStream(func(context huma.Context) { route.token(context, form) }), nil
}

func (route *tokenRoute) Register(api huma.API) {
	meta := route.GetMeta()
	meta.RequestBody = &huma.RequestBody{Content: formContent[oauthTokenForm](api)}
	meta.Responses = map[string]*huma.Response{
		"200": {Description: "OAuth token set", Content: jsonContent[tokenResponse](api)},
		"400": {Description: "Invalid token request", Content: jsonContent[oauthErrorResponse](api)},
		"401": {Description: "Invalid client", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	meta.Middlewares = huma.Middlewares{route.handler.publicCORS, noStore, requireContentType("application/x-www-form-urlencoded"), formBody(tokenBodyKey)}
	huma.Register(api, meta, route.Handler)
}

func noStore(context huma.Context, next func(huma.Context)) {
	context.SetHeader("Cache-Control", "no-store")
	context.SetHeader("Pragma", "no-cache")
	next(context)
}

func (route *tokenRoute) token(context huma.Context, form url.Values) {
	var tokens service.TokenSet
	var err error
	switch form.Get("grant_type") {
	case "authorization_code":
		if !required(form.Get("client_id"), form.Get("code"), form.Get("redirect_uri"), form.Get("code_verifier"), form.Get("resource")) {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "authorization code fields are required")
			return
		}
		tokens, err = route.handler.service.ExchangeAuthorizationCode(context.Context(), service.ExchangeAuthorizationCodeInput{
			ClientID:     form.Get("client_id"),
			Code:         form.Get("code"),
			RedirectURI:  form.Get("redirect_uri"),
			CodeVerifier: form.Get("code_verifier"),
			Resource:     form.Get("resource"),
		})
	case "refresh_token":
		if !required(form.Get("client_id"), form.Get("refresh_token"), form.Get("resource")) {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "refresh token fields are required")
			return
		}
		tokens, err = route.handler.service.Refresh(context.Context(), service.RefreshInput{
			ClientID:     form.Get("client_id"),
			RefreshToken: form.Get("refresh_token"),
			Scope:        form.Get("scope"),
			Resource:     form.Get("resource"),
		})
	default:
		writeOAuthError(context, http.StatusBadRequest, "unsupported_grant_type", "unsupported grant type")
		return
	}
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	writeJSON(context, http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    int(time.Until(tokens.AccessExpiresAt).Seconds()),
		RefreshToken: tokens.RefreshToken,
		Scope:        strings.Join(scopeNames(tokens), " "),
	})
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type oauthTokenForm struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Resource     string `json:"resource"`
}

func required(values ...string) bool {
	return !slices.Contains(values, "")
}

func scopeNames(tokens service.TokenSet) []string {
	scopes := make([]string, len(tokens.Scopes))
	for index, scope := range tokens.Scopes {
		scopes[index] = string(scope)
	}
	return scopes
}
