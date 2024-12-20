package twitch

import (
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type spanRoundTripper struct{}

func (t *spanRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	span := trace.SpanFromContext(r.Context())
	defer span.End()
	span.SetAttributes(
		attribute.String("twitch.rate-limit.limit", resp.Header.Get("Ratelimit-Limit")),
		attribute.String(
			"twitch.rate-limit.remaining",
			resp.Header.Get("Ratelimit-Remaining"),
		),
		attribute.String("twitch.rate-limit.reset", resp.Header.Get("Ratelimit-Reset")),
	)

	parsedReset, _ := strconv.Atoi(resp.Header.Get("Ratelimit-Reset"))
	if parsedReset != 0 {
		span.SetAttributes(
			attribute.String(
				"twitch.rate-limit.reset",
				time.Unix(int64(parsedReset), 0).String(),
			),
		)
	}

	if userID, ok := r.Context().Value(userIDCtxKey{}).(string); ok {
		span.SetAttributes(attribute.String("twitch.user.id", userID))
	}

	return resp, nil
}

func createHttpClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(
			&spanRoundTripper{},
		),
	}
}
