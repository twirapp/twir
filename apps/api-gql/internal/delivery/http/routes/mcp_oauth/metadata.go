package mcp_oauth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
)

type protectedResourceMetadata struct {
	handler *Handler
}

type metadataOutput[T any] struct {
	AccessControlAllowOrigin string `header:"Access-Control-Allow-Origin"`
	CacheControl             string `header:"Cache-Control"`
	Body                     T
}

var _ httpbase.Route[*emptyInput, *metadataOutput[protectedResourceMetadataResponse]] = (*protectedResourceMetadata)(nil)

func newProtectedResourceMetadata(handler *Handler) *protectedResourceMetadata {
	return &protectedResourceMetadata{handler: handler}
}

func (*protectedResourceMetadata) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-protected-resource", Method: http.MethodGet, Path: "/.well-known/oauth-protected-resource", Tags: []string{"MCP OAuth"}, Summary: "OAuth protected resource metadata", DefaultStatus: http.StatusOK}
}

func (route *protectedResourceMetadata) Handler(_ context.Context, _ *emptyInput) (*metadataOutput[protectedResourceMetadataResponse], error) {
	return &metadataOutput[protectedResourceMetadataResponse]{AccessControlAllowOrigin: "*", CacheControl: "public, max-age=3600", Body: protectedResourceMetadataResponse{Resource: route.handler.origin + "/api/mcp", AuthorizationServers: []string{route.handler.origin}, ScopesSupported: []string{"read", "write"}, BearerMethodsSupported: []string{"header"}}}, nil
}

func (route *protectedResourceMetadata) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Responses = map[string]*huma.Response{"200": {Description: "OAuth protected resource metadata", Content: jsonContent[protectedResourceMetadataResponse](api)}}
	huma.Register(api, meta, route.Handler)
}

type authorizationServerMetadata struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *metadataOutput[authorizationServerMetadataResponse]] = (*authorizationServerMetadata)(nil)

func newAuthorizationServerMetadata(handler *Handler) *authorizationServerMetadata {
	return &authorizationServerMetadata{handler: handler}
}

func (*authorizationServerMetadata) GetMeta() huma.Operation {
	return huma.Operation{OperationID: "mcp-oauth-authorization-server", Method: http.MethodGet, Path: "/.well-known/oauth-authorization-server", Tags: []string{"MCP OAuth"}, Summary: "OAuth authorization server metadata", DefaultStatus: http.StatusOK}
}

func (route *authorizationServerMetadata) Handler(_ context.Context, _ *emptyInput) (*metadataOutput[authorizationServerMetadataResponse], error) {
	return &metadataOutput[authorizationServerMetadataResponse]{AccessControlAllowOrigin: "*", CacheControl: "public, max-age=3600", Body: authorizationServerMetadataResponse{Issuer: route.handler.origin, AuthorizationEndpoint: route.handler.origin + "/api/oauth/authorize", TokenEndpoint: route.handler.origin + "/api/oauth/token", RegistrationEndpoint: route.handler.origin + "/api/oauth/register", RevocationEndpoint: route.handler.origin + "/api/oauth/revoke", ScopesSupported: []string{"read", "write"}, ResponseTypesSupported: []string{"code"}, GrantTypesSupported: []string{"authorization_code", "refresh_token"}, TokenEndpointAuthMethodsSupported: []string{"none"}, CodeChallengeMethodsSupported: []string{"S256"}}}, nil
}

func (route *authorizationServerMetadata) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Responses = map[string]*huma.Response{"200": {Description: "OAuth authorization server metadata", Content: jsonContent[authorizationServerMetadataResponse](api)}}
	huma.Register(api, meta, route.Handler)
}

type protectedResourceMetadataResponse struct {
	Resource               string   `json:"resource"`
	AuthorizationServers   []string `json:"authorization_servers"`
	ScopesSupported        []string `json:"scopes_supported"`
	BearerMethodsSupported []string `json:"bearer_methods_supported"`
}

type authorizationServerMetadataResponse struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
}
