package mcp_oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type consentDecision struct {
	Attempt        string   `json:"attempt"`
	ChannelID      string   `json:"channel_id"`
	CSRFToken      string   `json:"csrf_token"`
	Decision       string   `json:"decision"`
	ApprovedScopes []string `json:"approved_scopes,omitempty"`
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
	requested := make([]entity.Scope, len(attempt.RequestedScopes))
	for index, scope := range attempt.RequestedScopes {
		requested[index] = entity.Scope(scope)
	}
	normalized, err := entity.NormalizeScopes(requested)
	if err != nil {
		writeOAuthError(context, http.StatusBadRequest, "invalid_scope", "requested scope is not permitted")
		return
	}
	groups := entity.AllScopeGroups()
	requestedGroups := make([]consentScopeResponse, 0, len(groups))
	for _, group := range groups {
		if !entity.HasScope(normalized, group.Group, entity.ScopeActionRead) {
			continue
		}
		actions := []entity.ScopeAction{entity.ScopeActionRead}
		if entity.HasScope(normalized, group.Group, entity.ScopeActionEdit) {
			actions = append(actions, entity.ScopeActionEdit)
		}
		requestedGroups = append(requestedGroups, consentScopeResponse{
			Group:       group.Group,
			Name:        group.Name,
			Description: group.Description,
			Actions:     actions,
		})
	}
	channelID, _ := route.handler.selectedDashboard(context)
	response := consentClientResponse{ID: client.ClientID, Name: metadata.ClientName}
	if metadata.ClientURI != "" {
		response.URI = &metadata.ClientURI
	}
	writeJSON(context, http.StatusOK, consentResponse{
		Client:          response,
		ChannelID:       channelID.String(),
		RequestedScopes: requestedGroups,
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
	var approved []entity.Scope
	if decision.Decision == "deny" {
		if len(decision.ApprovedScopes) != 0 {
			writeOAuthError(context, http.StatusBadRequest, "invalid_request", "approved scopes are not allowed when denying consent")
			return
		}
	} else {
		var err error
		approved, err = entity.ParseScopes(strings.Join(decision.ApprovedScopes, " "))
		if err != nil {
			writeOAuthError(context, http.StatusBadRequest, "invalid_scope", "approved scope is not permitted")
			return
		}
		requested := make([]entity.Scope, len(attempt.RequestedScopes))
		for index, scope := range attempt.RequestedScopes {
			requested[index] = entity.Scope(scope)
		}
		if !entity.ScopeSubset(approved, requested) {
			writeOAuthError(context, http.StatusBadRequest, "invalid_scope", "approved scope is not permitted")
			return
		}
	}
	if err := route.handler.sessions.DeleteMCPOAuthAttempt(context.Context(), decision.Attempt); err != nil {
		writeOAuthError(context, http.StatusNotFound, "invalid_request", "authorization attempt not found")
		return
	}
	if decision.Decision == "deny" {
		route.handler.writeConsentRedirect(context, attempt, url.Values{"error": {"access_denied"}, "state": {attempt.ClientState}})
		return
	}
	issued, err := route.handler.service.CreateAuthorizationCode(context.Context(), service.CreateAuthorizationCodeInput{
		Authorize: service.AuthorizeInput{
			ClientID:            attempt.ClientID,
			RedirectURI:         attempt.RedirectURI,
			ResponseType:        "code",
			Scope:               strings.Join(entity.ScopeStrings(approved), " "),
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
