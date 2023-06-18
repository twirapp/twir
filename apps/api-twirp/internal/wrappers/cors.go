package wrappers

import (
	"github.com/rs/cors"
	"net/http"
)

func WithCors(base http.Handler) http.Handler {
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{"Content-Type"},
	})
	handler := corsWrapper.Handler(base)

	return handler
}
