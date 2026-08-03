package mcp_oauth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

func (handler *Handler) token(context *gin.Context) {
	context.Header("Cache-Control", "no-store")
	context.Header("Pragma", "no-cache")
	if !requireContentType(context, "application/x-www-form-urlencoded") || !parseForm(context) {
		return
	}
	form := context.Request.PostForm
	var tokens service.TokenSet
	var err error
	switch form.Get("grant_type") {
	case "authorization_code":
		if !required(form.Get("client_id"), form.Get("code"), form.Get("redirect_uri"), form.Get("code_verifier"), form.Get("resource")) {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "authorization code fields are required")
			return
		}
		tokens, err = handler.service.ExchangeAuthorizationCode(context.Request.Context(), service.ExchangeAuthorizationCodeInput{ClientID: form.Get("client_id"), Code: form.Get("code"), RedirectURI: form.Get("redirect_uri"), CodeVerifier: form.Get("code_verifier"), Resource: form.Get("resource")})
	case "refresh_token":
		if !required(form.Get("client_id"), form.Get("refresh_token"), form.Get("resource")) {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "refresh token fields are required")
			return
		}
		tokens, err = handler.service.Refresh(context.Request.Context(), service.RefreshInput{ClientID: form.Get("client_id"), RefreshToken: form.Get("refresh_token"), Scope: form.Get("scope"), Resource: form.Get("resource")})
	default:
		writeOAuthError(context, http.StatusBadRequest, "unsupported_grant_type", "unsupported grant type")
		return
	}
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	context.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "token_type": tokens.TokenType, "expires_in": int(time.Until(tokens.AccessExpiresAt).Seconds()), "refresh_token": tokens.RefreshToken, "scope": strings.Join(scopeNames(tokens), " ")})
}

func required(values ...string) bool {
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

func scopeNames(tokens service.TokenSet) []string {
	scopes := make([]string, len(tokens.Scopes))
	for index, scope := range tokens.Scopes {
		scopes[index] = string(scope)
	}
	return scopes
}
