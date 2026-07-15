package spotify

import (
	"context"
	"net/http"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetTrackReturnsNilWhenSpotifyHasNoActivePlayback(t *testing.T) {
	defaultClient := http.DefaultClient
	t.Cleanup(func() {
		http.DefaultClient = defaultClient
	})

	tests := []struct {
		name   string
		scopes []string
		path   string
	}{
		{
			name:   "player state scope",
			scopes: []string{"user-read-playback-state"},
			path:   "/v1/me/player",
		},
		{
			name:   "currently playing scope",
			scopes: []string{"user-read-currently-playing"},
			path:   "/v1/me/player/currently-playing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultClient = &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != tt.path {
						t.Errorf("request path = %q, want %q", req.URL.Path, tt.path)
					}
					if authorization := req.Header.Get("Authorization"); authorization != "Bearer access-token" {
						t.Errorf("Authorization = %q, want bearer token", authorization)
					}

					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       http.NoBody,
						Header:     make(http.Header),
						Request:    req,
					}, nil
				}),
			}

			track, err := NewStatic("access-token", tt.scopes).GetTrack(context.Background())
			if err != nil {
				t.Fatalf("GetTrack() error = %v", err)
			}
			if track != nil {
				t.Fatalf("GetTrack() track = %#v, want nil", track)
			}
		})
	}
}
