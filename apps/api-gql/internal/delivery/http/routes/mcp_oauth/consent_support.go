package mcp_oauth

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

func (handler *Handler) loadAttempt(context huma.Context, attemptID string) (authsessions.MCPOAuthAttempt, bool) {
	attempt, err := handler.sessions.GetMCPOAuthAttempt(context.Context(), attemptID)
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

func (handler *Handler) browserIdentity(context huma.Context) (uuid.UUID, uuid.UUID, bool) {
	userID, err := handler.sessions.GetInternalUserID(context.Context())
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

func (handler *Handler) selectedDashboard(context huma.Context) (uuid.UUID, bool) {
	raw, err := handler.sessions.GetSelectedDashboard(context.Context())
	channelID, parseErr := uuid.Parse(raw)
	if err != nil || parseErr != nil {
		writeOAuthError(context, http.StatusForbidden, "access_denied", "selected dashboard is unavailable")
		return uuid.Nil, false
	}
	return channelID, true
}

func (handler *Handler) writeConsentRedirect(context huma.Context, attempt authsessions.MCPOAuthAttempt, values url.Values) {
	redirect, err := oauthRedirect(attempt.RedirectURI, values)
	if err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	writeJSON(context, http.StatusOK, consentRedirectResponse{RedirectTo: redirect})
}

type consentClientResponse struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	URI  *string `json:"uri,omitempty"`
}

type consentResponse struct {
	Client          consentClientResponse  `json:"client"`
	ChannelID       string                 `json:"channel_id"`
	RequestedScopes []consentScopeResponse `json:"requested_scopes"`
	CSRFToken       string                 `json:"csrf_token"`
}

type consentScopeResponse struct {
	Group       entity.ScopeGroup    `json:"group"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Actions     []entity.ScopeAction `json:"actions"`
}

type consentRedirectResponse struct {
	RedirectTo string `json:"redirect_to"`
}
