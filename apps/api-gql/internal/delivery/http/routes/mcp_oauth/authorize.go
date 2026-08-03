package mcp_oauth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

const maxStateLength = 1024

func (handler *Handler) authorize(context *gin.Context) {
	state := context.Query("state")
	if state == "" || len(state) > maxStateLength {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "state is required")
		return
	}
	input := service.AuthorizeInput{ClientID: context.Query("client_id"), RedirectURI: context.Query("redirect_uri"), ResponseType: context.Query("response_type"), Scope: context.Query("scope"), Resource: context.Query("resource"), CodeChallenge: context.Query("code_challenge"), CodeChallengeMethod: context.Query("code_challenge_method")}
	client, err := handler.service.GetClient(context.Request.Context(), input.ClientID)
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	if !service.MatchesRegisteredRedirectURI(client.RedirectURIs, input.RedirectURI) {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid redirect URI")
		return
	}
	request, err := handler.service.ValidateAuthorizeInput(context.Request.Context(), input)
	if err != nil {
		handler.writeAuthorizeErrorRedirect(context, input.RedirectURI, state, err)
		return
	}
	attemptID, err := randomValue(handler.random)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	csrfToken, err := randomValue(handler.random)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	scopes := make([]string, len(request.Scopes))
	for index, scope := range request.Scopes {
		scopes[index] = string(scope)
	}
	attempt := authsessions.MCPOAuthAttempt{ClientID: request.Client.ClientID, RedirectURI: request.RedirectURI, ClientState: state, CodeChallenge: request.CodeChallenge, RequestedScopes: scopes, Resource: request.Resource, CSRFToken: csrfToken, ExpiresAt: time.Now().Add(10 * time.Minute)}
	if err := handler.sessions.SetMCPOAuthAttempt(context.Request.Context(), attemptID, attempt); err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	redirect, _ := url.Parse(handler.origin + "/dashboard/mcp/authorize")
	redirect.RawQuery = url.Values{"attempt": {attemptID}}.Encode()
	context.Redirect(http.StatusFound, redirect.String())
}

func (handler *Handler) writeAuthorizeErrorRedirect(context *gin.Context, redirectURI, state string, err error) {
	code, description := "server_error", "server error"
	var oauthError *service.OAuthError
	if errors.As(err, &oauthError) {
		code, description = string(oauthError.Code), oauthError.Description
	}
	redirect, redirectErr := oauthRedirect(redirectURI, url.Values{"error": {code}, "error_description": {description}, "state": {state}})
	if redirectErr != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	context.Redirect(http.StatusFound, redirect)
}

func scopesContain(scopes []string, wanted string) bool {
	return strings.Contains(" "+scopeString(scopes)+" ", " "+wanted+" ")
}
