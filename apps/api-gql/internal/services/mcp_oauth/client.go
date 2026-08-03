package mcp_oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strings"

	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	repository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

type RegisterClientInput struct{ Metadata json.RawMessage }
type clientMetadata struct {
	ClientName              string   `json:"client_name,omitempty"`
	ClientURI               string   `json:"client_uri,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
}
type AuthorizeInput struct{ ClientID, RedirectURI, ResponseType, Scope, Resource, CodeChallenge, CodeChallengeMethod string }
type AuthorizationRequest struct {
	Client                               entity.Client
	RedirectURI, CodeChallenge, Resource string
	Scopes                               []entity.Scope
}

var privateScheme = regexp.MustCompile(`^[a-z][a-z0-9+.-]*$`)

func (s *Service) RegisterClient(ctx context.Context, input RegisterClientInput) (entity.Client, error) {
	metadata, scopes, err := normalizeMetadata(input.Metadata)
	if err != nil {
		return entity.NilClient, err
	}
	clientID, err := s.randomValue()
	if err != nil {
		return entity.NilClient, err
	}
	raw, err := json.Marshal(metadata)
	if err != nil {
		return entity.NilClient, err
	}
	client, err := s.repository.CreateClient(ctx, repository.CreateClientInput{ClientID: clientID, Metadata: raw, RedirectURIs: metadata.RedirectURIs, GrantTypes: metadata.GrantTypes, ResponseTypes: metadata.ResponseTypes, TokenEndpointAuthMethod: metadata.TokenEndpointAuthMethod, Scopes: scopes})
	if err != nil {
		return entity.NilClient, err
	}
	return client, nil
}
func (s *Service) GetClient(ctx context.Context, clientID string) (entity.Client, error) {
	client, err := s.repository.GetClient(ctx, clientID)
	if err != nil {
		return entity.NilClient, oauthError(ErrorInvalidClient, "unknown client")
	}
	return client, nil
}
func (s *Service) ValidateAuthorizeInput(ctx context.Context, input AuthorizeInput) (AuthorizationRequest, error) {
	client, err := s.GetClient(ctx, input.ClientID)
	if err != nil {
		return AuthorizationRequest{}, err
	}
	if input.ResponseType != "code" || !slices.Contains(client.RedirectURIs, input.RedirectURI) || input.Resource != s.resource || input.CodeChallengeMethod != "S256" || !validChallenge(input.CodeChallenge) {
		return AuthorizationRequest{}, oauthError(ErrorInvalidRequest, "invalid authorization request")
	}
	scopes := client.Scopes
	if input.Scope != "" {
		scopes, err = parseScopes(input.Scope)
	}
	if err != nil || !scopeSubset(scopes, client.Scopes) {
		return AuthorizationRequest{}, oauthError(ErrorInvalidScope, "requested scope is not permitted")
	}
	return AuthorizationRequest{Client: client, RedirectURI: input.RedirectURI, CodeChallenge: input.CodeChallenge, Resource: s.resource, Scopes: scopes}, nil
}
func normalizeMetadata(raw json.RawMessage) (clientMetadata, []entity.Scope, error) {
	if len(raw) == 0 || len(raw) > 16*1024 {
		return clientMetadata{}, nil, oauthError(ErrorInvalidClientMetadata, "invalid client metadata")
	}
	var metadata clientMetadata
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return clientMetadata{}, nil, oauthError(ErrorInvalidClientMetadata, "invalid client metadata")
	}
	metadata.ClientName = strings.TrimSpace(metadata.ClientName)
	if metadata.ClientName == "" {
		metadata.ClientName = "MCP client"
	}
	if len(metadata.GrantTypes) == 0 {
		metadata.GrantTypes = []string{"authorization_code", "refresh_token"}
	}
	if len(metadata.ResponseTypes) == 0 {
		metadata.ResponseTypes = []string{"code"}
	}
	if metadata.TokenEndpointAuthMethod == "" {
		metadata.TokenEndpointAuthMethod = "none"
	}
	if metadata.Scope == "" {
		metadata.Scope = "read write"
	}
	if len(metadata.ClientName) > 256 || len(metadata.RedirectURIs) == 0 || len(metadata.RedirectURIs) > 16 || metadata.TokenEndpointAuthMethod != "none" || !sameSet(metadata.GrantTypes, []string{"authorization_code", "refresh_token"}) || !sameSet(metadata.ResponseTypes, []string{"code"}) {
		return clientMetadata{}, nil, oauthError(ErrorInvalidClientMetadata, "unsupported client metadata")
	}
	if metadata.ClientURI != "" {
		normalized, err := normalizeClientURI(metadata.ClientURI)
		if err != nil {
			return clientMetadata{}, nil, err
		}
		metadata.ClientURI = normalized
	}
	for i, redirect := range metadata.RedirectURIs {
		normalized, err := normalizeRedirectURI(redirect)
		if err != nil {
			return clientMetadata{}, nil, err
		}
		for _, previous := range metadata.RedirectURIs[:i] {
			if previous == normalized {
				return clientMetadata{}, nil, oauthError(ErrorInvalidClientMetadata, "duplicate redirect URI")
			}
		}
		metadata.RedirectURIs[i] = normalized
	}
	scopes, err := parseScopes(metadata.Scope)
	if err != nil || !slices.Contains(scopes, entity.ScopeRead) {
		return clientMetadata{}, nil, oauthError(ErrorInvalidClientMetadata, "invalid client scope")
	}
	metadata.GrantTypes = []string{"authorization_code", "refresh_token"}
	metadata.ResponseTypes = []string{"code"}
	metadata.Scope = strings.Join(scopeStrings(scopes), " ")
	return metadata, scopes, nil
}
func normalizeRedirectURI(raw string) (string, error) {
	if len(raw) == 0 || len(raw) > 2048 || strings.Contains(raw, "*") {
		return "", oauthError(ErrorInvalidClientMetadata, "invalid redirect URI")
	}
	u, err := url.Parse(raw)
	if err != nil || !u.IsAbs() || u.Fragment != "" || u.User != nil {
		return "", oauthError(ErrorInvalidClientMetadata, "invalid redirect URI")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	switch u.Scheme {
	case "https":
		if u.Hostname() == "" {
			return "", oauthError(ErrorInvalidClientMetadata, "invalid redirect URI")
		}
	case "http":
		if !localHost(u.Hostname()) {
			return "", oauthError(ErrorInvalidClientMetadata, "invalid redirect URI")
		}
	default:
		if (u.Host == "" && u.Path == "") || u.Opaque != "" || !privateScheme.MatchString(u.Scheme) {
			return "", oauthError(ErrorInvalidClientMetadata, "invalid redirect URI")
		}
	}
	return u.String(), nil
}
func normalizeClientURI(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.Fragment != "" || u.User != nil {
		return "", oauthError(ErrorInvalidClientMetadata, "invalid client URI")
	}
	return u.String(), nil
}
func localHost(host string) bool {
	host = strings.ToLower(host)
	if host == "localhost" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.To4() != nil && ip.To4()[0] == 127
}
func sameSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]bool{}
	for _, value := range got {
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	return slices.EqualFunc(want, want, func(left, right string) bool { return seen[left] && seen[right] })
}
func parseScopes(raw string) ([]entity.Scope, error) {
	values := strings.Fields(raw)
	seen := map[entity.Scope]bool{}
	for _, value := range values {
		scope := entity.Scope(value)
		if (scope != entity.ScopeRead && scope != entity.ScopeWrite) || seen[scope] {
			return nil, oauthError(ErrorInvalidScope, "invalid scope")
		}
		seen[scope] = true
	}
	if len(seen) == 0 {
		return nil, oauthError(ErrorInvalidScope, "scope is required")
	}
	result := []entity.Scope{entity.ScopeRead}
	if seen[entity.ScopeWrite] {
		result = append(result, entity.ScopeWrite)
	}
	return result, nil
}
func scopeSubset(requested, allowed []entity.Scope) bool {
	for _, scope := range requested {
		if !slices.Contains(allowed, scope) {
			return false
		}
	}
	return true
}
func scopeStrings(scopes []entity.Scope) []string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = string(scope)
	}
	return result
}
func validChallenge(value string) bool {
	bytes, err := base64.RawURLEncoding.DecodeString(value)
	return err == nil && len(bytes) == sha256.Size
}
