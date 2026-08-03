package mcp_oauth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

func TestHandler_generated_openapi_exposes_concrete_oauth_contract(t *testing.T) {
	handler := newTestHandler(t)
	_ = handler.router()
	openapi := handler.api.OpenAPI()

	if len(openapi.Security) != 0 {
		t.Fatalf("top-level security = %#v", openapi.Security)
	}

	type operationCheck struct {
		path   string
		method string
	}

	for _, check := range []operationCheck{
		{path: "/.well-known/oauth-protected-resource", method: http.MethodGet},
		{path: "/.well-known/oauth-authorization-server", method: http.MethodGet},
		{path: "/oauth/register", method: http.MethodOptions},
		{path: "/oauth/register", method: http.MethodPost},
		{path: "/oauth/authorize", method: http.MethodGet},
		{path: "/oauth/consent", method: http.MethodGet},
		{path: "/oauth/consent", method: http.MethodPost},
		{path: "/oauth/token", method: http.MethodOptions},
		{path: "/oauth/token", method: http.MethodPost},
		{path: "/oauth/revoke", method: http.MethodOptions},
		{path: "/oauth/revoke", method: http.MethodPost},
	} {
		operation := mustOperation(t, openapi, check.path, check.method)
		if len(operation.Security) != 0 {
			t.Fatalf("%s %s security = %#v", check.method, check.path, operation.Security)
		}
	}

	for _, check := range []responseSetCheck{
		{path: "/oauth/authorize", method: http.MethodGet, statuses: []string{"302", "400", "401", "500"}},
		{path: "/oauth/consent", method: http.MethodGet, statuses: []string{"200", "401", "403", "404", "410", "500"}},
	} {
		assertResponseStatuses(t, openapi, check)
	}

	assertParameters(t, openapi, "/oauth/authorize", http.MethodGet, "client_id", "redirect_uri", "response_type", "scope", "state", "resource", "code_challenge", "code_challenge_method")
	assertParameters(t, openapi, "/oauth/consent", http.MethodGet, "attempt")

	assertJSONRequestSchemaProps(t, openapi, "/oauth/register", http.MethodPost, "application/json", "client_name", "client_uri", "redirect_uris", "grant_types", "response_types", "token_endpoint_auth_method", "scope")
	assertJSONRequestSchemaProps(t, openapi, "/oauth/consent", http.MethodPost, "application/json", "attempt", "channel_id", "csrf_token", "decision", "access_level")
	assertFormRequestSchemaProps(t, openapi, "/oauth/token", http.MethodPost, "application/x-www-form-urlencoded", "grant_type", "client_id", "code", "redirect_uri", "code_verifier", "refresh_token", "scope", "resource")
	assertFormRequestSchemaProps(t, openapi, "/oauth/revoke", http.MethodPost, "application/x-www-form-urlencoded", "client_id", "token")

	assertJSONResponseSchemaProps(t, openapi, "/.well-known/oauth-protected-resource", http.MethodGet, http.StatusOK, "resource", "authorization_servers")
	assertJSONResponseSchemaProps(t, openapi, "/.well-known/oauth-authorization-server", http.MethodGet, http.StatusOK, "issuer", "token_endpoint")
	assertJSONResponseSchemaProps(t, openapi, "/oauth/register", http.MethodPost, http.StatusCreated, "client_id", "client_id_issued_at", "client_name")
	assertJSONResponseSchemaProps(t, openapi, "/oauth/consent", http.MethodGet, http.StatusOK, "client", "channel_id", "requested_scopes", "access_levels", "csrf_token")
	assertJSONResponseSchemaProps(t, openapi, "/oauth/consent", http.MethodPost, http.StatusOK, "redirect_to")
	assertJSONResponseSchemaProps(t, openapi, "/oauth/token", http.MethodPost, http.StatusOK, "access_token", "token_type", "expires_in", "refresh_token", "scope")

	for _, check := range []struct {
		path   string
		method string
		status int
	}{
		{path: "/oauth/authorize", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/oauth/authorize", method: http.MethodGet, status: http.StatusUnauthorized},
		{path: "/oauth/authorize", method: http.MethodGet, status: http.StatusInternalServerError},
		{path: "/oauth/consent", method: http.MethodGet, status: http.StatusUnauthorized},
		{path: "/oauth/consent", method: http.MethodGet, status: http.StatusForbidden},
		{path: "/oauth/consent", method: http.MethodGet, status: http.StatusNotFound},
		{path: "/oauth/consent", method: http.MethodGet, status: http.StatusGone},
		{path: "/oauth/consent", method: http.MethodGet, status: http.StatusInternalServerError},
	} {
		assertJSONResponseSchemaProps(t, openapi, check.path, check.method, check.status, "error", "error_description")
	}

	assertNoContentResponse(t, openapi, "/oauth/register", http.MethodOptions, http.StatusNoContent)
	assertNoContentResponse(t, openapi, "/oauth/token", http.MethodOptions, http.StatusNoContent)
	assertNoContentResponse(t, openapi, "/oauth/revoke", http.MethodOptions, http.StatusNoContent)

	assertAuthorizeRedirect(t, openapi, "/oauth/authorize", http.MethodGet)
}

func mustOperation(t *testing.T, openapi *huma.OpenAPI, path string, method string) *huma.Operation {
	t.Helper()
	item, ok := openapi.Paths[path]
	if !ok || item == nil {
		t.Fatalf("missing path %s", path)
	}

	var operation *huma.Operation
	switch method {
	case http.MethodGet:
		operation = item.Get
	case http.MethodPost:
		operation = item.Post
	case http.MethodOptions:
		operation = item.Options
	default:
		t.Fatalf("unsupported method %s", method)
	}

	if operation == nil {
		t.Fatalf("missing operation %s %s", method, path)
	}

	return operation
}

func assertParameters(t *testing.T, openapi *huma.OpenAPI, path string, method string, want ...string) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	got := map[string]struct{}{}
	for _, param := range operation.Parameters {
		got[param.Name] = struct{}{}
	}
	for _, name := range want {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s %s missing query parameter %q", method, path, name)
		}
	}
}

func assertJSONRequestSchemaProps(t *testing.T, openapi *huma.OpenAPI, path string, method string, contentType string, want ...string) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	if operation.RequestBody == nil {
		t.Fatalf("%s %s missing request body", method, path)
	}
	assertSchemaProps(t, openapi, operation.RequestBody.Content[contentType].Schema, want...)
}

func assertFormRequestSchemaProps(t *testing.T, openapi *huma.OpenAPI, path string, method string, contentType string, want ...string) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	if operation.RequestBody == nil {
		t.Fatalf("%s %s missing request body", method, path)
	}
	assertSchemaProps(t, openapi, operation.RequestBody.Content[contentType].Schema, want...)
}

func assertJSONResponseSchemaProps(t *testing.T, openapi *huma.OpenAPI, path string, method string, status int, want ...string) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	response := operation.Responses[fmt.Sprint(status)]
	if response == nil {
		t.Fatalf("%s %s missing response %d", method, path, status)
	}
	content := response.Content["application/json"]
	if content == nil {
		t.Fatalf("%s %s response %d missing json content", method, path, status)
	}
	assertSchemaProps(t, openapi, content.Schema, want...)
}

func assertNoContentResponse(t *testing.T, openapi *huma.OpenAPI, path string, method string, status int) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	response := operation.Responses[fmt.Sprint(status)]
	if response == nil {
		t.Fatalf("%s %s missing response %d", method, path, status)
	}
}

func assertAuthorizeRedirect(t *testing.T, openapi *huma.OpenAPI, path string, method string) {
	t.Helper()
	operation := mustOperation(t, openapi, path, method)
	response := operation.Responses["302"]
	if response == nil {
		t.Fatalf("%s %s missing 302 response", method, path)
	}
	header := response.Headers["Location"]
	if header == nil || header.Schema == nil || header.Schema.Type != huma.TypeString {
		t.Fatalf("%s %s location header schema = %#v", method, path, header)
	}
	if len(response.Content) != 0 {
		t.Fatalf("%s %s unexpected 302 content = %#v", method, path, response.Content)
	}
}

func assertSchemaProps(t *testing.T, openapi *huma.OpenAPI, schema *huma.Schema, want ...string) {
	t.Helper()
	resolved := resolveSchema(t, openapi, schema)
	for _, name := range want {
		if _, ok := resolved.Properties[name]; !ok {
			t.Fatalf("schema %q missing property %q; got keys %v", schemaLabel(schema), name, schemaKeys(resolved))
		}
	}
}

func resolveSchema(t *testing.T, openapi *huma.OpenAPI, schema *huma.Schema) *huma.Schema {
	t.Helper()
	if schema == nil {
		t.Fatal("schema is nil")
	}
	if schema.Ref == "" {
		return schema
	}
	refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
	resolved := openapi.Components.Schemas.Map()[refName]
	if resolved == nil {
		t.Fatalf("missing component schema %s", refName)
	}
	return resolved
}

func schemaLabel(schema *huma.Schema) string {
	if schema == nil {
		return "<nil>"
	}
	if schema.Ref != "" {
		return schema.Ref
	}
	return schema.Type
}

func schemaKeys(schema *huma.Schema) []string {
	keys := make([]string, 0, len(schema.Properties))
	for key := range schema.Properties {
		keys = append(keys, key)
	}
	return keys
}
