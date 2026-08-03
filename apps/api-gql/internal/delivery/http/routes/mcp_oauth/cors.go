package mcp_oauth

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
)

func (handler *Handler) publicCORS(context huma.Context, next func(huma.Context)) {
	handler.publicCORSHeaders(context)
	next(context)
}

func (handler *Handler) publicPreflight(context huma.Context, next func(huma.Context)) {
	handler.publicCORSHeaders(context)
	next(context)
}

func (*Handler) publicCORSHeaders(context huma.Context) {
	humagin.Unwrap(context).Writer.Header().Del("Access-Control-Allow-Credentials")
	context.SetHeader("Access-Control-Allow-Origin", "*")
	context.SetHeader("Access-Control-Allow-Methods", "POST, OPTIONS")
	context.SetHeader("Access-Control-Allow-Headers", "Content-Type")
	context.SetHeader("Access-Control-Max-Age", "600")
}
