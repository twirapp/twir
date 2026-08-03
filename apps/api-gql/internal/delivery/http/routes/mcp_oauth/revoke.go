package mcp_oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
)

func (handler *Handler) revoke(context *gin.Context) {
	context.Header("Cache-Control", "no-store")
	context.Header("Pragma", "no-cache")
	if !requireContentType(context, "application/x-www-form-urlencoded") || !parseForm(context) {
		return
	}
	clientID := context.Request.PostForm.Get("client_id")
	token := context.Request.PostForm.Get("token")
	if clientID == "" || token == "" {
		writeOAuthError(context, http.StatusBadRequest, "invalid_request", "client_id and token are required")
		return
	}
	if _, err := handler.service.GetClient(context.Request.Context(), clientID); err != nil {
		writeServiceError(context, err, false)
		return
	}
	if err := handler.service.Revoke(context.Request.Context(), service.RevokeInput{ClientID: clientID, Token: token}); err != nil {
		writeServiceError(context, err, false)
		return
	}
	context.Status(http.StatusOK)
}
