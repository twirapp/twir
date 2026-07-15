package spotify

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetTrackIncludesTiming(t *testing.T) {
	tests := []struct {
		name         string
		scopes       []string
		expectedPath string
	}{
		{
			name:         "currently playing endpoint",
			expectedPath: "/v1/me/player/currently-playing",
		},
		{
			name:         "player state endpoint",
			scopes:       []string{"user-read-playback-state"},
			expectedPath: "/v1/me/player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalClient := http.DefaultClient
			t.Cleanup(func() {
				http.DefaultClient = originalClient
			})

			http.DefaultClient = &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != tt.expectedPath {
						t.Fatalf("expected request path %q, got %q", tt.expectedPath, req.URL.Path)
					}

					body := `{
						"progress_ms": 70000,
						"item": {
							"duration_ms": 238000,
							"name": "Heat Waves",
							"artists": [{"name": "Glass Animals"}],
							"album": {"images": [{"url": "https://example.com/heat-waves.jpg"}]}
						},
						"is_playing": true
					}`

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				}),
			}

			track, err := NewStatic("token", tt.scopes).GetTrack(context.Background())
			if err != nil {
				t.Fatalf("GetTrack() error = %v", err)
			}
			if track.ProgressMs != 70000 {
				t.Errorf("ProgressMs = %d, want 70000", track.ProgressMs)
			}
			if track.DurationMs != 238000 {
				t.Errorf("DurationMs = %d, want 238000", track.DurationMs)
			}
		})
	}
}
