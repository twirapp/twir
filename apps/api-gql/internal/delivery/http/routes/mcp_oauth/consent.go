package mcp_oauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type consentDecision struct {
	Attempt     string `json:"attempt"`
	ChannelID   string `json:"channel_id"`
	CSRFToken   string `json:"csrf_token"`
	Decision    string `json:"decision"`
	AccessLevel string `json:"access_level"`
}

func (handler *Handler) getConsent(context *gin.Context) {
	attempt, ok := handler.loadAttempt(context, context.Query("attempt"))
	if !ok {
		return
	}
	if _, _, ok := handler.browserIdentity(context); !ok {
		return
	}
	client, err := handler.service.GetClient(context.Request.Context(), attempt.ClientID)
	if err != nil {
		writeOAuthError(context, http.StatusNotFound, "invalid_client", "client not found")
		return
	}
	var metadata registeredClientMetadata
	if err := json.Unmarshal(client.Metadata, &metadata); err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	levels := []string{"read"}
	if scopesContain(attempt.RequestedScopes, "write") {
		levels = append(levels, "write")
	}
	channelID, _ := handler.selectedDashboard(context)
	clientResponse := gin.H{"id": client.ClientID, "name": metadata.ClientName}
	if metadata.ClientURI != "" {
		clientResponse["uri"] = metadata.ClientURI
	}
	context.JSON(http.StatusOK, gin.H{"client": clientResponse, "channel_id": channelID.String(), "requested_scopes": attempt.RequestedScopes, "access_levels": levels, "csrf_token": attempt.CSRFToken})
}

func (handler *Handler) postConsent(context *gin.Context) {
	if !requireContentType(context, "application/json") || context.GetHeader("Origin") != handler.origin {
		if context.Writer.Written() {
			return
		}
		writeOAuthError(context, http.StatusForbidden, "access_denied", "invalid origin")
		return
	}
	var decision consentDecision
	if !decodeJSON(context, &decision) {
		return
	}
	attempt, ok := handler.loadAttempt(context, decision.Attempt)
	if !ok {
		return
	}
	userID, channelID, ok := handler.browserIdentity(context)
	if !ok {
		return
	}
	if decision.ChannelID == "" {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "channel_id is required")
		return
	}
	if decision.ChannelID != channelID.String() {
		writeOAuthError(context, http.StatusConflict, "dashboard_changed", "selected dashboard changed")
		return
	}
	if !csrfMatches(attempt.CSRFToken, decision.CSRFToken) {
		writeOAuthError(context, http.StatusForbidden, "access_denied", "invalid CSRF token")
		return
	}
	if decision.Decision != "approve" && decision.Decision != "deny" {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid consent decision")
		return
	}
	if decision.Decision == "approve" && (decision.AccessLevel != "read" && decision.AccessLevel != "write") {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "invalid access level")
		return
	}
	if decision.AccessLevel == "write" && !scopesContain(attempt.RequestedScopes, "write") {
		writeOAuthError(context, http.StatusBadRequest, "invalid_scope", "requested scope is not permitted")
		return
	}
	if err := handler.sessions.DeleteMCPOAuthAttempt(context.Request.Context(), decision.Attempt); err != nil {
		writeOAuthError(context, http.StatusNotFound, "invalid_request", "authorization attempt not found")
		return
	}
	if decision.Decision == "deny" {
		handler.writeConsentRedirect(context, attempt, url.Values{"error": {"access_denied"}, "state": {attempt.ClientState}})
		return
	}
	scope := "read"
	if decision.AccessLevel == "write" {
		scope = "read write"
	}
	issued, err := handler.service.CreateAuthorizationCode(context.Request.Context(), service.CreateAuthorizationCodeInput{Authorize: service.AuthorizeInput{ClientID: attempt.ClientID, RedirectURI: attempt.RedirectURI, ResponseType: "code", Scope: scope, Resource: attempt.Resource, CodeChallenge: attempt.CodeChallenge, CodeChallengeMethod: "S256"}, ChannelID: channelID, ApprovingUserID: userID})
	if err != nil {
		writeServiceError(context, err, true)
		return
	}
	handler.writeConsentRedirect(context, attempt, url.Values{"code": {issued.Code}, "state": {attempt.ClientState}})
}

func (handler *Handler) loadAttempt(context *gin.Context, attemptID string) (authsessions.MCPOAuthAttempt, bool) {
	attempt, err := handler.sessions.GetMCPOAuthAttempt(context.Request.Context(), attemptID)
	if errors.Is(err, authsessions.ErrMCPOAuthAttemptExpired) {
		writeOAuthError(context, http.StatusGone, "invalid_request", "authorization attempt expired")
		return authsessions.MCPOAuthAttempt{}, false
	}
	if err != nil {
		writeOAuthError(context, http.StatusNotFound, "invalid_request", "authorization attempt not found")
		return authsessions.MCPOAuthAttempt{}, false
	}
	return attempt, true
}

func (handler *Handler) browserIdentity(context *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userID, err := handler.sessions.GetInternalUserID(context.Request.Context())
	if err != nil {
		writeOAuthError(context, http.StatusUnauthorized, "login_required", "login required")
		return uuid.Nil, uuid.Nil, false
	}
	channelID, ok := handler.selectedDashboard(context)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return userID, channelID, true
}

func (handler *Handler) selectedDashboard(context *gin.Context) (uuid.UUID, bool) {
	raw, err := handler.sessions.GetSelectedDashboard(context.Request.Context())
	channelID, parseErr := uuid.Parse(raw)
	if err != nil || parseErr != nil {
		writeOAuthError(context, http.StatusForbidden, "access_denied", "selected dashboard is unavailable")
		return uuid.Nil, false
	}
	return channelID, true
}

func (handler *Handler) writeConsentRedirect(context *gin.Context, attempt authsessions.MCPOAuthAttempt, values url.Values) {
	redirect, err := oauthRedirect(attempt.RedirectURI, values)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	context.JSON(http.StatusOK, gin.H{"redirect_to": redirect})
}
