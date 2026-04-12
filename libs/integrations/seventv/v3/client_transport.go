package v3

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type headerTransport struct {
	base  http.RoundTripper
	token string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	span := trace.SpanFromContext(req.Context())
	if span.SpanContext().IsValid() {
		span.SetAttributes(
			attribute.String("seventv.ratelimit.limit", resp.Header.Get("X-Ratelimit-Limit")),
			attribute.String("seventv.ratelimit.remaining", resp.Header.Get("X-Ratelimit-Remaining")),
			attribute.String("seventv.ratelimit.reset", resp.Header.Get("X-Ratelimit-Reset")),
		)
	}

	return resp, nil
}
