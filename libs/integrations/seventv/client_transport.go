package seventv

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type headerTransport struct {
	base  http.RoundTripper
	token string
}

func (t *headerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()
	r.Header.Set("Authorization", "Bearer "+t.token)

	resp, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String(
			"7tv.x-ratelimit-global-limit",
			resp.Header.Get("X-RateLimit-Global-Limit"),
		),
		attribute.String(
			"7tv.x-ratelimit-global-remaining",
			resp.Header.Get("X-RateLimit-Global-Remaining"),
		),
		attribute.String("7tv.x-ratelimit-global-reset", resp.Header.Get("X-RateLimit-Global-Reset")),
		attribute.String("7tv.x-ratelimit-global-used", resp.Header.Get("X-RateLimit-Global-Used")),
	)

	return resp, nil
}
