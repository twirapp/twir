package mcp_oauth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

const requestBodyLimit = 16 * 1024

type oauthErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

type parsedBodyKey string

const (
	registerBodyKey parsedBodyKey = "register"
	consentBodyKey  parsedBodyKey = "consent"
	tokenBodyKey    parsedBodyKey = "token"
	revokeBodyKey   parsedBodyKey = "revoke"
)

func rawStream(handler func(huma.Context)) *huma.StreamResponse {
	return &huma.StreamResponse{Body: handler}
}

func schemaFor[T any](api huma.API) *huma.Schema {
	return huma.SchemaFromType(api.OpenAPI().Components.Schemas, reflect.TypeFor[T]())
}

func jsonContent[T any](api huma.API) map[string]*huma.MediaType {
	return map[string]*huma.MediaType{"application/json": {Schema: schemaFor[T](api)}}
}

func formContent[T any](api huma.API) map[string]*huma.MediaType {
	return map[string]*huma.MediaType{"application/x-www-form-urlencoded": {Schema: schemaFor[T](api)}}
}

func parsedBody[T any](ctx context.Context, key parsedBodyKey) (T, bool) {
	value, ok := ctx.Value(key).(T)
	return value, ok
}

func writeJSON[T any](context huma.Context, status int, response T) {
	encoded, err := json.Marshal(response)
	if err != nil {
		context.SetStatus(http.StatusInternalServerError)
		return
	}

	context.SetHeader("Content-Type", "application/json; charset=utf-8")
	context.SetStatus(status)
	if _, err := context.BodyWriter().Write(encoded); err != nil {
		return
	}
}

func writeOAuthError(context huma.Context, status int, code string, description string) {
	writeJSON(context, status, oauthErrorResponse{Code: code, Description: description})
}

func writeServiceError(context huma.Context, err error, consent bool) {
	var oauthError *service.OAuthError
	if !errors.As(err, &oauthError) {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	status := http.StatusBadRequest
	if consent && oauthError.Code == service.ErrorAccessDenied {
		status = http.StatusForbidden
	}
	if oauthError.Code == service.ErrorInvalidClient {
		status = http.StatusUnauthorized
	}
	writeOAuthError(context, status, string(oauthError.Code), oauthError.Description)
}

func requireContentType(expected string) func(huma.Context, func(huma.Context)) {
	return func(context huma.Context, next func(huma.Context)) {
		contentType, _, err := mime.ParseMediaType(context.Header("Content-Type"))
		if err != nil || contentType != expected {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid content type")
			return
		}
		next(context)
	}
}

func exactOrigin(origin string) func(huma.Context, func(huma.Context)) {
	return func(context huma.Context, next func(huma.Context)) {
		if context.Header("Origin") != origin {
			writeOAuthError(context, http.StatusForbidden, "access_denied", "invalid origin")
			return
		}
		next(context)
	}
}

func jsonBody[T any](key parsedBodyKey) func(huma.Context, func(huma.Context)) {
	return func(context huma.Context, next func(huma.Context)) {
		raw := humagin.Unwrap(context)
		raw.Request.Body = http.MaxBytesReader(raw.Writer, raw.Request.Body, requestBodyLimit)
		var body T
		decoder := json.NewDecoder(raw.Request.Body)
		if err := decoder.Decode(&body); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid JSON body")
			return
		}
		next(huma.WithValue(context, key, body))
	}
}

func formBody(key parsedBodyKey) func(huma.Context, func(huma.Context)) {
	return func(context huma.Context, next func(huma.Context)) {
		raw := humagin.Unwrap(context)
		raw.Request.Body = http.MaxBytesReader(raw.Writer, raw.Request.Body, requestBodyLimit)
		if err := raw.Request.ParseForm(); err != nil {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid form body")
			return
		}
		if context.Header("Authorization") != "" || raw.Request.PostForm.Has("client_secret") {
			if context.Header("Authorization") != "" {
				context.SetHeader("WWW-Authenticate", `Basic realm="oauth"`)
			}
			writeOAuthError(context, http.StatusUnauthorized, "invalid_client", "client authentication is not supported")
			return
		}
		next(huma.WithValue(context, key, raw.Request.PostForm))
	}
}

func writeRedirect(context huma.Context, location string) {
	raw := humagin.Unwrap(context)
	http.Redirect(raw.Writer, raw.Request, location, http.StatusFound)
}

func randomValue(reader io.Reader) (string, error) {
	bytes := make([]byte, 32)
	if _, err := io.ReadFull(reader, bytes); err != nil {
		return "", fmt.Errorf("generate OAuth value: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func csrfMatches(expected, actual string) bool {
	expectedHash := sha256.Sum256([]byte(expected))
	actualHash := sha256.Sum256([]byte(actual))
	return subtle.ConstantTimeCompare(expectedHash[:], actualHash[:]) == 1
}

func oauthRedirect(raw string, values url.Values) (string, error) {
	redirect, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	query := redirect.Query()
	maps.Copy(query, values)
	redirect.RawQuery = query.Encode()
	return redirect.String(), nil
}

func scopeString(scopes []string) string { return strings.Join(scopes, " ") }
