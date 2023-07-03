package wrappers

import (
	"context"
	"net/http"
)

func WithDashboardId(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var channelId string

		header := r.Header.Get("dashboard-id")
		if header != "" {
			channelId = header
		}

		ctx = context.WithValue(ctx, "dashboardId", channelId)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
