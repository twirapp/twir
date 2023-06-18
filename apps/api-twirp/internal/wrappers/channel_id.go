package wrappers

import (
	"context"
	"net/http"
)

func WithChannelId(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookies := r.Cookies()
		var channelId string
		for _, cookie := range cookies {
			if cookie.Name == "Channel-id" {
				channelId = cookie.Value
			}
		}
		ctx = context.WithValue(ctx, "channelId", channelId)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
