package spotify

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestSpotify_GetQueue_returns_tracks(t *testing.T) {
	defaultClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %s, want GET", req.Method)
			}
			if req.URL.Path != "/v1/me/player/queue" {
				t.Fatalf("path = %s, want /v1/me/player/queue", req.URL.Path)
			}
			if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
				t.Fatalf("Authorization = %q, want bearer token", auth)
			}

			body := `{"queue":[{"id":"track-1","uri":"spotify:track:track-1","name":"Song","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[{"url":"https://img/1.jpg"}]},"duration_ms":1234,"type":"track"}]}`
			return spotifyTestResponse(req, http.StatusOK, body, nil), nil
		}),
	}
	t.Cleanup(func() {
		http.DefaultClient = defaultClient
	})

	tracks, err := NewStatic("access-token", nil).GetQueue(context.Background())
	if err != nil {
		t.Fatalf("GetQueue() error = %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("len(tracks) = %d, want 1", len(tracks))
	}
	if tracks[0].ArtistName != "Artist" || tracks[0].AlbumName != "Album" || tracks[0].ImageURL != "https://img/1.jpg" {
		t.Fatalf("track = %#v, want parsed fields", tracks[0])
	}
}

func TestSpotify_AddToQueue_handles_success_and_errors(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       int
		body             string
		headers          http.Header
		deviceID         string
		wantErr          error
		wantRateAfterSec int
		wantHasDeviceID  bool
	}{
		{
			name:            "success without device id",
			statusCode:      http.StatusNoContent,
			deviceID:        "",
			wantHasDeviceID: false,
		},
		{
			name:            "success with 200 and snapshot body",
			statusCode:      http.StatusOK,
			body:            `6VwrDGxNUFaCcdjX8YvfTJWKWZs`,
			deviceID:        "device-1",
			wantHasDeviceID: true,
		},
		{
			name:            "no active device",
			statusCode:      http.StatusNotFound,
			body:            `{"error":{"status":404,"reason":"NO_ACTIVE_DEVICE","message":"No active device found"}}`,
			deviceID:        "device-1",
			wantErr:         ErrNoActiveDevice,
			wantHasDeviceID: true,
		},
		{
			name:            "insufficient scope",
			statusCode:      http.StatusForbidden,
			body:            `{"error":{"status":403,"message":"insufficient client scope"}}`,
			deviceID:        "device-1",
			wantErr:         ErrInsufficientScope,
			wantHasDeviceID: true,
		},
		{
			name:             "rate limited",
			statusCode:       http.StatusTooManyRequests,
			body:             `{"error":{"status":429,"message":"rate limited"}}`,
			headers:          http.Header{"Retry-After": []string{"17"}},
			deviceID:         "device-1",
			wantErr:          ErrRateLimited,
			wantRateAfterSec: 17,
			wantHasDeviceID:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultClient := http.DefaultClient
			http.DefaultClient = &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodPost {
						t.Fatalf("method = %s, want POST", req.Method)
					}
					if req.URL.Path != "/v1/me/player/queue" {
						t.Fatalf("path = %s, want /v1/me/player/queue", req.URL.Path)
					}
					if req.URL.Query().Get("uri") != "spotify:track:track-1" {
						t.Fatalf("uri = %q, want track uri", req.URL.Query().Get("uri"))
					}
					_, hasDeviceID := req.URL.Query()["device_id"]
					if hasDeviceID != tt.wantHasDeviceID {
						t.Fatalf("device_id present = %v, want %v", hasDeviceID, tt.wantHasDeviceID)
					}
					if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
						t.Fatalf("Authorization = %q, want bearer token", auth)
					}

					headers := make(http.Header)
					for key, values := range tt.headers {
						headers[key] = append([]string(nil), values...)
					}
					return spotifyTestResponse(req, tt.statusCode, tt.body, headers), nil
				}),
			}
			t.Cleanup(func() {
				http.DefaultClient = defaultClient
			})

			err := NewStatic("access-token", nil).AddToQueue(context.Background(), "spotify:track:track-1", tt.deviceID)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("AddToQueue() error = %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddToQueue() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantRateAfterSec > 0 {
				var rateLimitErr *RateLimitedError
				if !errors.As(err, &rateLimitErr) {
					t.Fatalf("AddToQueue() error = %v, want RateLimitedError", err)
				}
				if rateLimitErr.RetryAfterSeconds != tt.wantRateAfterSec {
					t.Fatalf("RetryAfterSeconds = %d, want %d", rateLimitErr.RetryAfterSeconds, tt.wantRateAfterSec)
				}
			}
		})
	}
}

func TestSpotify_SkipNext_handles_success_and_rate_limit(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		body          string
		headers       http.Header
		wantErr       error
		wantRateAfter int
	}{
		{
			name:       "success",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "success with 200",
			statusCode: http.StatusOK,
		},
		{
			name:          "rate limited",
			statusCode:    http.StatusTooManyRequests,
			body:          `{"error":{"status":429,"message":"rate limited"}}`,
			headers:       http.Header{"Retry-After": []string{"21"}},
			wantErr:       ErrRateLimited,
			wantRateAfter: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultClient := http.DefaultClient
			http.DefaultClient = &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodPost {
						t.Fatalf("method = %s, want POST", req.Method)
					}
					if req.URL.Path != "/v1/me/player/next" {
						t.Fatalf("path = %s, want /v1/me/player/next", req.URL.Path)
					}
					if req.URL.Query().Get("device_id") != "device-1" {
						t.Fatalf("device_id = %q, want device-1", req.URL.Query().Get("device_id"))
					}
					if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
						t.Fatalf("Authorization = %q, want bearer token", auth)
					}

					headers := make(http.Header)
					for key, values := range tt.headers {
						headers[key] = append([]string(nil), values...)
					}
					return spotifyTestResponse(req, tt.statusCode, tt.body, headers), nil
				}),
			}
			t.Cleanup(func() {
				http.DefaultClient = defaultClient
			})

			err := NewStatic("access-token", nil).SkipNext(context.Background(), "device-1")
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("SkipNext() error = %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("SkipNext() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantRateAfter > 0 {
				var rateLimitErr *RateLimitedError
				if !errors.As(err, &rateLimitErr) {
					t.Fatalf("SkipNext() error = %v, want RateLimitedError", err)
				}
				if rateLimitErr.RetryAfterSeconds != tt.wantRateAfter {
					t.Fatalf("RetryAfterSeconds = %d, want %d", rateLimitErr.RetryAfterSeconds, tt.wantRateAfter)
				}
			}
		})
	}
}
