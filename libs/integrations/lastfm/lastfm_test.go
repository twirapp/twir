package lastfm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const testGuard = 500 * time.Millisecond

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackingBody struct {
	io.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

func testResponse(statusCode int, body io.ReadCloser) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       body,
		Header:     make(http.Header),
	}
}

func newTestClient(t *testing.T, transport http.RoundTripper) *Lastfm {
	t.Helper()

	client, err := New(context.Background(), Opts{
		ApiKey:     "api-key",
		UserName:   "listener",
		HTTPClient: &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("create Last.fm client: %v", err)
	}

	return client
}

func TestNewWithUserNameNeedsOnlyAPIKeyAndMakesNoRequest(t *testing.T) {
	t.Parallel()

	requests := 0
	client, err := New(context.Background(), Opts{
		ApiKey:   "api-key",
		UserName: "listener",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			requests++
			return nil, errors.New("unexpected request")
		})},
	})
	if err != nil {
		t.Fatalf("create Last.fm client: %v", err)
	}
	if requests != 0 {
		t.Fatalf("expected no requests, got %d", requests)
	}
	if client.userName != "listener" {
		t.Fatalf("expected initialized username, got %q", client.userName)
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := New(context.Background(), Opts{UserName: "listener"})
	if err == nil {
		t.Fatal("expected missing API key error")
	}
}

func TestNewUsesDedicatedClientWithFiniteTimeout(t *testing.T) {
	t.Parallel()

	client, err := New(context.Background(), Opts{ApiKey: "api-key", UserName: "listener"})
	if err != nil {
		t.Fatalf("create Last.fm client: %v", err)
	}
	if client.httpClient == nil {
		t.Fatal("expected HTTP client")
	}
	if client.httpClient == http.DefaultClient {
		t.Fatal("must not use http.DefaultClient")
	}
	if client.httpClient.Timeout != 10*time.Second {
		t.Fatalf("expected 10s timeout, got %s", client.httpClient.Timeout)
	}
}

func TestNewResolvesLegacyUserNameWithSignedRequest(t *testing.T) {
	t.Parallel()

	requests := 0
	client, err := New(context.Background(), Opts{
		ApiKey:       "api-key",
		ClientSecret: "client-secret",
		SessionKey:   "session-key",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.URL.Scheme != "https" || req.URL.Host != "ws.audioscrobbler.com" || req.URL.Path != "/2.0/" {
				t.Errorf("unexpected endpoint: %s", req.URL.String())
			}
			if req.URL.RawQuery != "" {
				t.Errorf("legacy credentials must not appear in query: %s", req.URL.RawQuery)
			}
			if got := req.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
				t.Errorf("unexpected content type: %q", got)
			}
			if err := req.ParseForm(); err != nil {
				t.Fatalf("parse legacy request form: %v", err)
			}
			for key, expected := range map[string]string{
				"method":  "user.getinfo",
				"api_key": "api-key",
				"sk":      "session-key",
				"api_sig": "24c497d45fbcb97543d5be64b8bd9aa6",
				"format":  "json",
			} {
				if got := req.PostForm.Get(key); got != expected {
					t.Errorf("expected %s=%q, got %q", key, expected, got)
				}
			}

			return testResponse(
				http.StatusOK,
				io.NopCloser(strings.NewReader(`{"user":{"name":"legacy-listener"}}`)),
			), nil
		})},
	})
	if err != nil {
		t.Fatalf("create legacy Last.fm client: %v", err)
	}
	if requests != 1 {
		t.Fatalf("expected one username request, got %d", requests)
	}
	if client.userName != "legacy-listener" {
		t.Fatalf("expected resolved username, got %q", client.userName)
	}
}

func TestNewLegacyUserNameRequiresCredentials(t *testing.T) {
	t.Parallel()

	_, err := New(context.Background(), Opts{ApiKey: "api-key"})
	if err == nil {
		t.Fatal("expected missing legacy credentials error")
	}
}

func TestNewLegacyUserNamePropagatesCancellation(t *testing.T) {
	t.Parallel()

	requestStarted := make(chan struct{})
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		close(requestStarted)
		<-req.Context().Done()
		return nil, req.Context().Err()
	})}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	result := make(chan error, 1)
	go func() {
		_, err := New(ctx, Opts{
			ApiKey:       "api-key",
			ClientSecret: "client-secret",
			SessionKey:   "session-key",
			HTTPClient:   httpClient,
		})
		result <- err
	}()

	waitForRequest(t, requestStarted)
	cancel()
	assertCanceledPromptly(t, result)
}

func TestNewLegacyUserNameReturnsBoundedAPIErrorAndClosesBody(t *testing.T) {
	t.Parallel()

	body := &trackingBody{Reader: strings.NewReader(`{"error":9,"message":"Invalid session key"}`)}
	_, err := New(context.Background(), Opts{
		ApiKey:       "api-key",
		ClientSecret: "client-secret",
		SessionKey:   "session-key",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return testResponse(http.StatusUnauthorized, body), nil
		})},
	})
	if err == nil {
		t.Fatal("expected legacy lookup error")
	}
	if !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "Invalid session key") {
		t.Fatalf("expected status and API body, got %v", err)
	}
	if !body.closed {
		t.Fatal("expected response body to be closed")
	}
}

func TestGetTrackMapsNowPlayingTrackAndClosesBody(t *testing.T) {
	t.Parallel()

	body := &trackingBody{Reader: strings.NewReader(`{
		"recenttracks": {
			"track": [{
				"name": "Track title",
				"artist": {"#text": "Track artist"},
				"image": [
					{"size": "small", "#text": "https://example.com/small.jpg"},
					{"size": "large", "#text": "https://example.com/large.jpg"}
				],
				"@attr": {"nowplaying": "true"},
				"date": {"uts": "1783872000", "#text": "12 Jul 2026, 12:00"}
			}]
		}
	}`)}
	client := newTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRecentTracksQuery(t, req, "1")
		return testResponse(http.StatusOK, body), nil
	}))

	track, err := client.GetTrack(context.Background())
	if err != nil {
		t.Fatalf("get track: %v", err)
	}
	if track == nil {
		t.Fatal("expected track, got nil")
	}
	if track.Title != "Track title" || track.Artist != "Track artist" ||
		track.Image != "https://example.com/small.jpg" || track.PlayedUTS != "12 Jul 2026, 12:00" {
		t.Fatalf("unexpected track: %#v", track)
	}
	if !body.closed {
		t.Fatal("expected response body to be closed")
	}
}

func TestGetTrackReturnsNilWithoutNowPlayingTrack(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{name: "empty recent tracks", body: `{"recenttracks":{"track":[]}}`},
		{
			name: "most recent track is not playing",
			body: `{"recenttracks":{"track":[{"name":"Stopped","@attr":{"nowplaying":"false"}}]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
				return testResponse(http.StatusOK, io.NopCloser(strings.NewReader(tt.body))), nil
			}))
			track, err := client.GetTrack(context.Background())
			if err != nil {
				t.Fatalf("get track: %v", err)
			}
			if track != nil {
				t.Fatalf("expected nil track, got %#v", track)
			}
		})
	}
}

func TestGetTrackIncludesNonSuccessBodyAndClosesIt(t *testing.T) {
	t.Parallel()

	body := &trackingBody{Reader: strings.NewReader(`{"error":11,"message":"Service Offline"}`)}
	client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(http.StatusBadGateway, body), nil
	}))

	track, err := client.GetTrack(context.Background())
	if err == nil {
		t.Fatal("expected non-success status error")
	}
	if track != nil {
		t.Fatalf("expected nil track, got %#v", track)
	}
	if !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "Service Offline") {
		t.Fatalf("expected status and API body, got %v", err)
	}
	if !body.closed {
		t.Fatal("expected response body to be closed")
	}
}

func TestGetTrackDetectsAPIErrorEnvelopeOnSuccessStatus(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(
			http.StatusOK,
			io.NopCloser(strings.NewReader(`{"error":6,"message":"User not found"}`)),
		), nil
	}))

	track, err := client.GetTrack(context.Background())
	if err == nil {
		t.Fatal("expected Last.fm API error")
	}
	if track != nil {
		t.Fatalf("expected nil track, got %#v", track)
	}
	if !strings.Contains(err.Error(), "6") || !strings.Contains(err.Error(), "User not found") {
		t.Fatalf("expected API code and message, got %v", err)
	}
}

func TestGetTrackCapsNonSuccessErrorBody(t *testing.T) {
	t.Parallel()

	const omittedTail = "must-not-appear-in-error"
	apiBody := "api-error-prefix:" + strings.Repeat("x", 8<<10) + omittedTail
	client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(
			http.StatusBadGateway,
			io.NopCloser(strings.NewReader(apiBody)),
		), nil
	}))

	_, err := client.GetTrack(context.Background())
	if err == nil {
		t.Fatal("expected non-success status error")
	}
	if !strings.Contains(err.Error(), "api-error-prefix") || !strings.Contains(err.Error(), "truncated") {
		t.Fatalf("expected bounded API body and truncation marker, got %v", err)
	}
	if strings.Contains(err.Error(), omittedTail) {
		t.Fatal("error leaked body content beyond the error-body cap")
	}
	if len(err.Error()) > (4<<10)+256 {
		t.Fatalf("non-success error is too large: %d bytes", len(err.Error()))
	}
}

func TestGetTrackRejectsOversizedResponseWithoutReadingPastBound(t *testing.T) {
	t.Parallel()

	reader := strings.NewReader(strings.Repeat("x", (1<<20)+128))
	body := &trackingBody{Reader: reader}
	client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(http.StatusOK, body), nil
	}))

	_, err := client.GetTrack(context.Background())
	if err == nil || !strings.Contains(err.Error(), "exceeds 1048576 bytes") {
		t.Fatalf("expected bounded response error, got %v", err)
	}
	if reader.Len() == 0 {
		t.Fatal("response reader was consumed beyond the configured bound")
	}
	if !body.closed {
		t.Fatal("expected oversized response body to be closed")
	}
}

func TestGetTrackInvalidJSONReturnsBoundedError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(http.StatusOK, io.NopCloser(strings.NewReader(`{"recenttracks":`))), nil
	}))

	_, err := client.GetTrack(context.Background())
	if err == nil {
		t.Fatal("expected decode error")
	}
	if len(err.Error()) > 512 {
		t.Fatalf("decode error unexpectedly includes response body: %d bytes", len(err.Error()))
	}
}

func TestGetTrackPropagatesCancellation(t *testing.T) {
	t.Parallel()

	requestStarted := make(chan struct{})
	client := newTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		close(requestStarted)
		<-req.Context().Done()
		return nil, req.Context().Err()
	}))
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	result := make(chan error, 1)
	go func() {
		_, err := client.GetTrack(ctx)
		result <- err
	}()

	waitForRequest(t, requestStarted)
	cancel()
	assertCanceledPromptly(t, result)
}

func TestGetTrackHonorsHTTPClientTimeout(t *testing.T) {
	t.Parallel()

	client, err := New(context.Background(), Opts{
		ApiKey:   "api-key",
		UserName: "listener",
		HTTPClient: &http.Client{
			Timeout: 25 * time.Millisecond,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				<-req.Context().Done()
				return nil, req.Context().Err()
			}),
		},
	})
	if err != nil {
		t.Fatalf("create Last.fm client: %v", err)
	}

	startedAt := time.Now()
	_, err = client.GetTrack(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	if elapsed := time.Since(startedAt); elapsed >= testGuard {
		t.Fatalf("client timeout was not prompt: %s", elapsed)
	}
}

func TestGetRecentTracksMapsHistoryAndSkipsNowPlaying(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRecentTracksQuery(t, req, "2")
		return testResponse(http.StatusOK, io.NopCloser(strings.NewReader(`{
			"recenttracks":{"track":[
				{"name":"Current","artist":{"#text":"Current artist"},"@attr":{"nowplaying":"true"}},
				{"name":"Previous","artist":{"#text":"Previous artist"},"image":[{"#text":"cover.jpg"}],"date":{"uts":"123","#text":"date"}},
				{"name":"Older","artist":{"#text":"Older artist"},"date":{"uts":"122","#text":"date"}}
			]}
		}`))), nil
	}))

	tracks, err := client.GetRecentTracks(context.Background(), 2)
	if err != nil {
		t.Fatalf("get recent tracks: %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("expected two historical tracks, got %#v", tracks)
	}
	if tracks[0] != (Track{Title: "Previous", Artist: "Previous artist", Image: "cover.jpg", PlayedUTS: "123"}) {
		t.Fatalf("unexpected first historical track: %#v", tracks[0])
	}
	if tracks[1] != (Track{Title: "Older", Artist: "Older artist", PlayedUTS: "122"}) {
		t.Fatalf("unexpected second historical track: %#v", tracks[1])
	}
}

func TestGetRecentTracksNormalizesLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		limit    int
		expected string
	}{
		{name: "zero", limit: 0, expected: "10"},
		{name: "negative", limit: -1, expected: "10"},
		{name: "too large", limit: 51, expected: "10"},
		{name: "valid", limit: 7, expected: "7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
				assertRecentTracksQuery(t, req, tt.expected)
				return testResponse(
					http.StatusOK,
					io.NopCloser(strings.NewReader(`{"recenttracks":{"track":[]}}`)),
				), nil
			}))
			if _, err := client.GetRecentTracks(context.Background(), tt.limit); err != nil {
				t.Fatalf("get recent tracks: %v", err)
			}
		})
	}
}

func TestGetRecentTracksPropagatesCancellation(t *testing.T) {
	t.Parallel()

	requestStarted := make(chan struct{})
	client := newTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		close(requestStarted)
		<-req.Context().Done()
		return nil, req.Context().Err()
	}))
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	result := make(chan error, 1)
	go func() {
		_, err := client.GetRecentTracks(ctx, 10)
		result <- err
	}()

	waitForRequest(t, requestStarted)
	cancel()
	assertCanceledPromptly(t, result)
}

func assertRecentTracksQuery(t *testing.T, req *http.Request, expectedLimit string) {
	t.Helper()

	if req.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", req.Method)
	}
	if req.URL.Scheme != "https" || req.URL.Host != "ws.audioscrobbler.com" || req.URL.Path != "/2.0/" {
		t.Errorf("unexpected endpoint: %s", req.URL.String())
	}
	query := req.URL.Query()
	for key, expected := range map[string]string{
		"method":  "user.getrecenttracks",
		"api_key": "api-key",
		"user":    "listener",
		"limit":   expectedLimit,
		"format":  "json",
	} {
		if got := query.Get(key); got != expected {
			t.Errorf("expected %s=%q, got %q", key, expected, got)
		}
	}
}

func waitForRequest(t *testing.T, started <-chan struct{}) {
	t.Helper()

	select {
	case <-started:
	case <-time.After(testGuard):
		t.Fatal("timed out waiting for Last.fm request")
	}
}

func assertCanceledPromptly(t *testing.T, result <-chan error) {
	t.Helper()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(testGuard):
		t.Fatal("Last.fm request did not return promptly after cancellation")
	}
}
