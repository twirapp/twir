package twitch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type spanRoundTripper struct{}

func (t *spanRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %v", err)
		}

		span.SetAttributes(
			attribute.String("twitch.request.body", string(bodyBytes)),
		)

		// Restore the body since reading it consumes it
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("twitch.rate-limit.reset", resp.Header.Get("Ratelimit-Reset")),
		attribute.Int("http_status_code", resp.StatusCode),
	)

	parsedLimit, _ := strconv.Atoi(resp.Header.Get("Ratelimit-Limit"))
	if parsedLimit != 0 {
		span.SetAttributes(
			attribute.Int(
				"twitch.rate-limit.limit",
				parsedLimit,
			),
		)
	}

	parsedRemaining, _ := strconv.Atoi(resp.Header.Get("Ratelimit-Remaining"))
	if parsedRemaining != 0 {
		span.SetAttributes(
			attribute.Int(
				"twitch.rate-limit.remaining",
				parsedRemaining,
			),
		)
	}

	parsedReset, _ := strconv.Atoi(resp.Header.Get("Ratelimit-Reset"))
	if parsedReset != 0 {
		span.SetAttributes(
			attribute.String(
				"twitch.rate-limit.reset",
				time.Unix(int64(parsedReset), 0).String(),
			),
		)
	}

	// response body
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		span.SetAttributes(
			attribute.String("twitch.response.body", string(bodyBytes)),
		)

		// Restore the body since reading it consumes it
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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
