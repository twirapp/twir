package wrappers

import (
	"context"
	"net/http"
)

func WithApiKeyHeader(base http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			header := r.Header.Get("api-key")
			ctx = context.WithValue(ctx, "apiKey", header)
			r = r.WithContext(ctx)

			base.ServeHTTP(w, r)
		},
	)
}
