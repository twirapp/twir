package mcp_oauth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type revokeOptions struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *preflightOutput] = (*revokeOptions)(nil)

func newRevokeOptions(handler *Handler) *revokeOptions {
	return &revokeOptions{handler: handler}
}

func (*revokeOptions) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-revoke-options", Method: http.MethodOptions, Path: "/oauth/revoke", DefaultStatus: http.StatusNoContent}
}

func (*revokeOptions) Handler(context.Context, *emptyInput) (*preflightOutput, error) {
	return &preflightOutput{Status: http.StatusNoContent}, nil
}

func (route *revokeOptions) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Middlewares = huma.Middlewares{route.handler.publicPreflight}
	huma.Register(api, meta, route.Handler)
}

type revokeRoute struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *huma.StreamResponse] = (*revokeRoute)(nil)

func newRevoke(handler *Handler) *revokeRoute {
	return &revokeRoute{handler: handler}
}

func (*revokeRoute) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-revoke", Method: http.MethodPost, Path: "/oauth/revoke", Tags: []string{"MCP OAuth"}, Summary: "Revoke an OAuth token", DefaultStatus: http.StatusOK}
}

func (route *revokeRoute) Handler(ctx context.Context, _ *emptyInput) (*huma.StreamResponse, error) {
	form, ok := parsedBody[url.Values](ctx, revokeBodyKey)
	if !ok {
		return rawStream(func(context huma.Context) {
			writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		}), nil
	}
	return rawStream(func(context huma.Context) { route.revoke(context, form) }), nil
}

func (route *revokeRoute) Register(api huma.API) {
	meta := route.GetMeta()
	meta.RequestBody = &huma.RequestBody{Content: formContent[revokeForm](api)}
	meta.Responses = map[string]*huma.Response{
		"200": {Description: "Token revoked"},
		"400": {Description: "Invalid revoke request", Content: jsonContent[oauthErrorResponse](api)},
		"401": {Description: "Invalid client", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	meta.Middlewares = huma.Middlewares{route.handler.publicCORS, noStore, requireContentType("application/x-www-form-urlencoded"), formBody(revokeBodyKey)}
	huma.Register(api, meta, route.Handler)
}

func (route *revokeRoute) revoke(context huma.Context, form url.Values) {
	clientID := form.Get("client_id")
	token := form.Get("token")
	if clientID == "" || token == "" {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "client_id and token are required")
		return
	}
	if _, err := route.handler.service.GetClient(context.Context(), clientID); err != nil {
		writeServiceError(context, err, false)
		return
	}
	if err := route.handler.service.Revoke(context.Context(), service.RevokeInput{ClientID: clientID, Token: token}); err != nil {
		writeServiceError(context, err, false)
		return
	}
	context.SetStatus(http.StatusOK)
}

type revokeForm struct {
	ClientID string `json:"client_id"`
	Token    string `json:"token"`
}
