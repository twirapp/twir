package mcp_oauth

import "github.com/gin-gonic/gin"

func (handler *Handler) protectedResourceMetadata(context *gin.Context) {
	handler.metadataHeaders(context)
	context.JSON(200, gin.H{"resource": handler.origin + "/api/mcp", "authorization_servers": []string{handler.origin}, "scopes_supported": []string{"read", "write"}, "bearer_methods_supported": []string{"header"}})
}

func (handler *Handler) authorizationServerMetadata(context *gin.Context) {
	handler.metadataHeaders(context)
	context.JSON(200, gin.H{"issuer": handler.origin, "authorization_endpoint": handler.origin + "/api/oauth/authorize", "token_endpoint": handler.origin + "/api/oauth/token", "registration_endpoint": handler.origin + "/api/oauth/register", "revocation_endpoint": handler.origin + "/api/oauth/revoke", "scopes_supported": []string{"read", "write"}, "response_types_supported": []string{"code"}, "grant_types_supported": []string{"authorization_code", "refresh_token"}, "token_endpoint_auth_methods_supported": []string{"none"}, "code_challenge_methods_supported": []string{"S256"}})
}

func (*Handler) metadataHeaders(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Cache-Control", "public, max-age=3600")
}
