package wrappers

import (
	"context"
	"net/http"
)

const ContextHeadersKey = "request-headers"

func WithHeaders(base http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = context.WithValue(ctx, "request-headers", r.Header)
			r = r.WithContext(ctx)

			base.ServeHTTP(w, r)
		},
	)
}
