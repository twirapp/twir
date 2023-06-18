package wrappers

import (
	"context"
	"net/http"
)

func WithDashboardId(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookies := r.Cookies()
		var channelId string
		for _, cookie := range cookies {
			if cookie.Name == "dashboard_id" {
				channelId = cookie.Value
			}
		}

		header := r.Header.Get("dashboard-id")
		if header != "" {
			channelId = header
		}

		ctx = context.WithValue(ctx, "dashboard_id", channelId)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
