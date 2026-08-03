package mcp_oauth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

const maxStateLength = 1024

type authorizeRoute struct {
	handler *Handler
}

type authorizeInput struct {
	ClientID            string `query:"client_id"`
	RedirectURI         string `query:"redirect_uri"`
	ResponseType        string `query:"response_type"`
	Scope               string `query:"scope"`
	State               string `query:"state"`
	Resource            string `query:"resource"`
	CodeChallenge       string `query:"code_challenge"`
	CodeChallengeMethod string `query:"code_challenge_method"`
}

var _ httpbase.Route[*authorizeInput, *huma.StreamResponse] = (*authorizeRoute)(nil)

func newAuthorize(handler *Handler) *authorizeRoute {
	return &authorizeRoute{handler: handler}
}

func (*authorizeRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-authorize",
		Method:        http.MethodGet,
		Path:          "/oauth/authorize",
		Tags:          []string{"MCP OAuth"},
		Summary:       "Begin OAuth authorization",
		DefaultStatus: http.StatusFound,
		Responses: map[string]*huma.Response{
			"302": {
				Description: "Redirect to consent or the validated client callback",
				Headers: map[string]*huma.Header{
					"Location": {
						Schema: &huma.Schema{Type: huma.TypeString},
					},
				},
			},
		},
	}
}

func (route *authorizeRoute) Handler(_ context.Context, input *authorizeInput) (*huma.StreamResponse, error) {
	return rawStream(func(context huma.Context) { route.authorize(context, *input) }), nil
}

func (route *authorizeRoute) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Responses = map[string]*huma.Response{
		"302": {
			Description: "Redirect to consent or the validated client callback",
			Headers: map[string]*huma.Header{
				"Location": {
					Schema: &huma.Schema{Type: huma.TypeString},
				},
			},
		},
		"400": {Description: "Invalid request", Content: jsonContent[oauthErrorResponse](api)},
		"401": {Description: "Unauthorized", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	huma.Register(api, meta, route.Handler)
}

func (route *authorizeRoute) authorize(context huma.Context, input authorizeInput) {
	state := input.State
	if state == "" || len(state) > maxStateLength {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "state is required")
		return
	}
	authorizeInput := service.AuthorizeInput{
		ClientID:            input.ClientID,
		RedirectURI:         input.RedirectURI,
		ResponseType:        input.ResponseType,
		Scope:               input.Scope,
		Resource:            input.Resource,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
	}
	client, err := route.handler.service.GetClient(context.Context(), authorizeInput.ClientID)
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	if !service.MatchesRegisteredRedirectURI(client.RedirectURIs, authorizeInput.RedirectURI) {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid redirect URI")
		return
	}
	request, err := route.handler.service.ValidateAuthorizeInput(context.Context(), authorizeInput)
	if err != nil {
		route.writeAuthorizeErrorRedirect(context, authorizeInput.RedirectURI, state, err)
		return
	}
	attemptID, err := randomValue(route.handler.random)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	csrfToken, err := randomValue(route.handler.random)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	scopes := make([]string, len(request.Scopes))
	for index, scope := range request.Scopes {
		scopes[index] = string(scope)
	}
	attempt := authsessions.MCPOAuthAttempt{
		ClientID:        request.Client.ClientID,
		RedirectURI:     request.RedirectURI,
		ClientState:     state,
		CodeChallenge:   request.CodeChallenge,
		RequestedScopes: scopes,
		Resource:        request.Resource,
		CSRFToken:       csrfToken,
		ExpiresAt:       time.Now().Add(10 * time.Minute),
	}
	if err := route.handler.sessions.SetMCPOAuthAttempt(context.Context(), attemptID, attempt); err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	redirect, _ := url.Parse(route.handler.origin + "/dashboard/mcp/authorize")
	redirect.RawQuery = url.Values{"attempt": {attemptID}}.Encode()
	writeRedirect(context, redirect.String())
}

func (route *authorizeRoute) writeAuthorizeErrorRedirect(context huma.Context, redirectURI, state string, err error) {
	code, description := "server_error", "server error"
	if oauthError, ok := errors.AsType[*service.OAuthError](err); ok {
		code, description = string(oauthError.Code), oauthError.Description
	}
	redirect, redirectErr := oauthRedirect(redirectURI, url.Values{"error": {code}, "error_description": {description}, "state": {state}})
	if redirectErr != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	writeRedirect(context, redirect)
}
