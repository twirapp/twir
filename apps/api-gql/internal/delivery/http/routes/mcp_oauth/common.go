package mcp_oauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

const requestBodyLimit = 16 * 1024

func writeOAuthError(context *gin.Context, status int, code string, description string) {
	context.JSON(status, gin.H{"error": code, "error_description": description})
}

func writeServiceError(context *gin.Context, err error, consent bool) {
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

func requireContentType(context *gin.Context, expected string) bool {
	contentType, _, err := mime.ParseMediaType(context.GetHeader("Content-Type"))
	if err != nil || contentType != expected {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid content type")
		return false
	}
	return true
}

func decodeJSON(context *gin.Context, destination any) bool {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, requestBodyLimit)
	decoder := json.NewDecoder(context.Request.Body)
	if err := decoder.Decode(destination); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return false
	}
	return true
}

func parseForm(context *gin.Context) bool {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, requestBodyLimit)
	if err := context.Request.ParseForm(); err != nil {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid form body")
		return false
	}
	if context.GetHeader("Authorization") != "" || context.Request.PostForm.Has("client_secret") {
		if context.GetHeader("Authorization") != "" {
			context.Header("WWW-Authenticate", `Basic realm="oauth"`)
		}
		writeOAuthError(context, http.StatusUnauthorized, "invalid_client", "client authentication is not supported")
		return false
	}
	return true
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
	for key, values := range values {
		query[key] = values
	}
	redirect.RawQuery = query.Encode()
	return redirect.String(), nil
}

func scopeString(scopes []string) string { return strings.Join(scopes, " ") }
