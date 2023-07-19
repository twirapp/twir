package wrappers

import (
	"net/http"

	"github.com/rs/cors"
)

func WithCors(base http.Handler) http.Handler {
	corsWrapper := cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"POST"},
			AllowedHeaders: []string{"Content-Type"},
		},
	)
	handler := corsWrapper.Handler(base)

	return handler
}
