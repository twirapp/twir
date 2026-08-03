package mcp_oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type consentDecision struct {
	Attempt     string `json:"attempt"`
	ChannelID   string `json:"channel_id"`
	CSRFToken   string `json:"csrf_token"`
	Decision    string `json:"decision"`
	AccessLevel string `json:"access_level"`
}

type getConsentRoute struct {
	handler *Handler
}

type getConsentInput struct {
	Attempt string `query:"attempt"`
}

var _ httpbase.Route[*getConsentInput, *huma.StreamResponse] = (*getConsentRoute)(nil)

func newGetConsent(handler *Handler) *getConsentRoute {
	return &getConsentRoute{handler: handler}
}

func (*getConsentRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-consent",
		Method:        http.MethodGet,
		Path:          "/oauth/consent",
		Tags:          []string{"MCP OAuth"},
		Summary:       "Get OAuth consent details",
		DefaultStatus: http.StatusOK,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Consent details",
			},
		},
	}
}

func (route *getConsentRoute) Handler(_ context.Context, input *getConsentInput) (*huma.StreamResponse, error) {
	return rawStream(func(context huma.Context) { route.getConsent(context, input.Attempt) }), nil
}

func (route *getConsentRoute) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Responses = map[string]*huma.Response{
		"200": {Description: "Consent details", Content: jsonContent[consentResponse](api)},
		"401": {Description: "Unauthorized", Content: jsonContent[oauthErrorResponse](api)},
		"403": {Description: "Forbidden", Content: jsonContent[oauthErrorResponse](api)},
		"404": {Description: "Not found", Content: jsonContent[oauthErrorResponse](api)},
		"410": {Description: "Gone", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	huma.Register(api, meta, route.Handler)
}

func (route *getConsentRoute) getConsent(context huma.Context, attemptID string) {
	attempt, ok := route.handler.loadAttempt(context, attemptID)
	if !ok {
		return
	}
	if _, _, ok := route.handler.browserIdentity(context); !ok {
		return
	}
	client, err := route.handler.service.GetClient(context.Context(), attempt.ClientID)
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
	channelID, _ := route.handler.selectedDashboard(context)
	response := consentClientResponse{ID: client.ClientID, Name: metadata.ClientName}
	if metadata.ClientURI != "" {
		response.URI = &metadata.ClientURI
	}
	writeJSON(context, http.StatusOK, consentResponse{
		Client:          response,
		ChannelID:       channelID.String(),
		RequestedScopes: attempt.RequestedScopes,
		AccessLevels:    levels,
		CSRFToken:       attempt.CSRFToken,
	})
}

type postConsentRoute struct {
	handler *Handler
}

var _ httpbase.Route[*emptyInput, *huma.StreamResponse] = (*postConsentRoute)(nil)

func newPostConsent(handler *Handler) *postConsentRoute {
	return &postConsentRoute{handler: handler}
}

func (*postConsentRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-consent-submit",
		Method:        http.MethodPost,
		Path:          "/oauth/consent",
		Tags:          []string{"MCP OAuth"},
		Summary:       "Submit OAuth consent",
		DefaultStatus: http.StatusOK,
	}
}

func (route *postConsentRoute) Handler(ctx context.Context, _ *emptyInput) (*huma.StreamResponse, error) {
	decision, ok := parsedBody[consentDecision](ctx, consentBodyKey)
	if !ok {
		return rawStream(func(context huma.Context) {
			writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		}), nil
	}
	return rawStream(func(context huma.Context) { route.postConsent(context, decision) }), nil
}

func (route *postConsentRoute) Register(api huma.API) {
	meta := route.GetMeta()
	meta.RequestBody = &huma.RequestBody{Content: jsonContent[consentDecision](api)}
	meta.Responses = map[string]*huma.Response{
		"200": {Description: "Authorization callback redirect", Content: jsonContent[consentRedirectResponse](api)},
		"400": {Description: "Invalid consent", Content: jsonContent[oauthErrorResponse](api)},
		"401": {Description: "Unauthorized", Content: jsonContent[oauthErrorResponse](api)},
		"403": {Description: "Forbidden", Content: jsonContent[oauthErrorResponse](api)},
		"404": {Description: "Not found", Content: jsonContent[oauthErrorResponse](api)},
		"409": {Description: "Conflict", Content: jsonContent[oauthErrorResponse](api)},
		"410": {Description: "Gone", Content: jsonContent[oauthErrorResponse](api)},
		"500": {Description: "Server error", Content: jsonContent[oauthErrorResponse](api)},
	}
	meta.Middlewares = huma.Middlewares{requireContentType("application/json"), exactOrigin(route.handler.origin), jsonBody[consentDecision](consentBodyKey)}
	huma.Register(api, meta, route.Handler)
}

func (route *postConsentRoute) postConsent(context huma.Context, decision consentDecision) {
	attempt, ok := route.handler.loadAttempt(context, decision.Attempt)
	if !ok {
		return
	}
	userID, channelID, ok := route.handler.browserIdentity(context)
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
	if err := route.handler.sessions.DeleteMCPOAuthAttempt(context.Context(), decision.Attempt); err != nil {
		writeOAuthError(context, http.StatusNotFound, "invalid_request", "authorization attempt not found")
		return
	}
	if decision.Decision == "deny" {
		route.handler.writeConsentRedirect(context, attempt, url.Values{"error": {"access_denied"}, "state": {attempt.ClientState}})
		return
	}
	scope := "read"
	if decision.AccessLevel == "write" {
		scope = "read write"
	}
	issued, err := route.handler.service.CreateAuthorizationCode(context.Context(), service.CreateAuthorizationCodeInput{
		Authorize: service.AuthorizeInput{
			ClientID:            attempt.ClientID,
			RedirectURI:         attempt.RedirectURI,
			ResponseType:        "code",
			Scope:               scope,
			Resource:            attempt.Resource,
			CodeChallenge:       attempt.CodeChallenge,
			CodeChallengeMethod: "S256",
		},
		ChannelID:       channelID,
		ApprovingUserID: userID,
	})
	if err != nil {
		writeServiceError(context, err, true)
		return
	}
	route.handler.writeConsentRedirect(context, attempt, url.Values{"code": {issued.Code}, "state": {attempt.ClientState}})
}

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
	Client          consentClientResponse `json:"client"`
	ChannelID       string                `json:"channel_id"`
	RequestedScopes []string              `json:"requested_scopes"`
	AccessLevels    []string              `json:"access_levels"`
	CSRFToken       string                `json:"csrf_token"`
}

type consentRedirectResponse struct {
	RedirectTo string `json:"redirect_to"`
}
