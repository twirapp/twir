package mcp_oauth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

type registeredClientMetadata struct {
	ClientName              string   `json:"client_name"`
	ClientURI               string   `json:"client_uri,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
}

func (handler *Handler) register(context *gin.Context) {
	if !requireContentType(context, "application/json") {
		return
	}
	var metadata json.RawMessage
	if !decodeJSON(context, &metadata) {
		return
	}
	var secret struct {
		ClientSecret json.RawMessage `json:"client_secret"`
	}
	if json.Unmarshal(metadata, &secret) != nil || secret.ClientSecret != nil {
		writeOAuthError(context, http.StatusBadRequest, "invalid_client_metadata", "client secrets are not supported")
		return
	}
	client, err := handler.service.RegisterClient(context.Request.Context(), service.RegisterClientInput{Metadata: metadata})
	if err != nil {
		writeServiceError(context, err, false)
		return
	}
	var response registeredClientMetadata
	if err := json.Unmarshal(client.Metadata, &response); err != nil {
		writeOAuthError(context, http.StatusInternalServerError, "server_error", "server error")
		return
	}
	context.JSON(http.StatusCreated, struct {
		ClientID       string `json:"client_id"`
		ClientIssuedAt int64  `json:"client_id_issued_at"`
		registeredClientMetadata
	}{ClientID: client.ClientID, ClientIssuedAt: client.CreatedAt.Unix(), registeredClientMetadata: response})
}
