package mcp_oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type registeredClientMetadata struct {
	ClientName              string   `json:"client_name"`
	ClientURI               string   `json:"client_uri,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
}

type registerOptions struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *preflightOutput] = (*registerOptions)(nil)

func newRegisterOptions(handler *Handler) *registerOptions {
	return &registerOptions{handler: handler}
}

func (route *registerOptions) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-register-options", Method: http.MethodOptions, Path: "/oauth/register", DefaultStatus: http.StatusNoContent}
}

func (*registerOptions) Handler(context.Context, *emptyInput) (*preflightOutput, error) {
	return &preflightOutput{Status: http.StatusNoContent}, nil
}

func (route *registerOptions) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Middlewares = huma.Middlewares{route.handler.publicPreflight}
	huma.Register(api, meta, route.Handler)
}

type registerClient struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *huma.StreamResponse] = (*registerClient)(nil)

func newRegisterClient(handler *Handler) *registerClient {
	return &registerClient{handler: handler}
}

func (*registerClient) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-register", Method: http.MethodPost, Path: "/oauth/register", Tags: []string{"MCP OAuth"}, Summary: "Register an OAuth client", DefaultStatus: http.StatusCreated}
}

func (route *registerClient) Handler(ctx context.Context, _ *emptyInput) (*huma.StreamResponse, error) {
	metadata, ok := parsedBody[json.RawMessage](ctx, registerBodyKey)
	if !ok {
		return rawStream(func(context huma.Context) {
			writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		}), nil
	}
	return rawStream(func(context huma.Context) { route.register(context, metadata) }), nil
}

func (route *registerClient) Register(api huma.API) {
	meta := route.GetMeta()
	meta.RequestBody = &huma.RequestBody{Content: jsonContent[registeredClientMetadata](api)}
	meta.Responses = map[string]*huma.Response{
		"201": {Description: "Registered client", Content: jsonContent[registerClientResponse](api)},
		"400": {Description: "Invalid client metadata", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	meta.Middlewares = huma.Middlewares{route.handler.publicCORS, route.handler.registerRateLimit, requireContentType("application/json"), jsonBody[json.RawMessage](registerBodyKey)}
	huma.Register(api, meta, route.Handler)
}

func (handler *Handler) registerRateLimit(context huma.Context, next func(huma.Context)) {
	if handler.registerRateLimiter == nil {
		next(context)
		return
	}

	response, err := handler.registerRateLimiter.Use(
		context.Context(),
		&rate_limiter.LeakyOptions{
			KeyPrefix:       fmt.Sprintf("api-gql:ratelimiter:%s:mcp-oauth-register", humagin.Unwrap(context).ClientIP()),
			MaximumCapacity: 20,
			WindowSeconds:   60,
		},
		1,
	)
	if err != nil {
		context.SetStatus(http.StatusInternalServerError)
		return
	}

	context.SetHeader("X-Rate-Limit-Bucket", "mcp-oauth-register")
	context.SetHeader("X-Rate-Limit-Limit", "20")
	context.SetHeader("X-Rate-Limit-Remaining", fmt.Sprint(response.RemainingTokens))
	context.SetHeader("X-Rate-Limit-Reset", fmt.Sprint(response.ResetAt.Unix()))
	if !response.Success {
		context.SetStatus(http.StatusTooManyRequests)
		return
	}

	next(context)
}

func (route *registerClient) register(context huma.Context, metadata json.RawMessage) {
	var secret struct {
		ClientSecret json.RawMessage `json:"client_secret"`
	}
	if json.Unmarshal(metadata, &secret) != nil || secret.ClientSecret != nil {
		writeOAuthError(context, http.StatusBadRequest, "invalid_client_metadata", "client secrets are not supported")
		return
	}
	client, err := route.handler.service.RegisterClient(context.Context(), service.RegisterClientInput{Metadata: metadata})
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	var response registeredClientMetadata
	if err := json.Unmarshal(client.Metadata, &response); err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	writeJSON(context, http.StatusCreated, registerClientResponse{ClientID: client.ClientID, ClientIssuedAt: client.CreatedAt.Unix(), ClientName: response.ClientName, ClientURI: response.ClientURI, RedirectURIs: response.RedirectURIs, GrantTypes: response.GrantTypes, ResponseTypes: response.ResponseTypes, TokenEndpointAuthMethod: response.TokenEndpointAuthMethod, Scope: response.Scope})
}

type registerClientResponse struct {
	ClientID                string   `json:"client_id"`
	ClientIssuedAt          int64    `json:"client_id_issued_at"`
	ClientName              string   `json:"client_name"`
	ClientURI               string   `json:"client_uri,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
}
