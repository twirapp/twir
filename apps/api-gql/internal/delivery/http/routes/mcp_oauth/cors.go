package mcp_oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) publicCORS(context *gin.Context) {
	handler.publicCORSHeaders(context)
	context.Next()
}

func (handler *Handler) publicPreflight(context *gin.Context) {
	handler.publicCORSHeaders(context)
	context.Status(http.StatusNoContent)
}

func (*Handler) publicCORSHeaders(context *gin.Context) {
	context.Writer.Header().Del("Access-Control-Allow-Credentials")
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
	context.Header("Access-Control-Allow-Headers", "Content-Type")
	context.Header("Access-Control-Max-Age", "600")
}
