package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestSpotify_GetCurrentlyPlaying_returns_track_and_nil_when_empty(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantNil    bool
	}{
		{
			name:       "returns track",
			statusCode: http.StatusOK,
			body:       `{"currently_playing_type":"track","item":{"id":"track-1","uri":"spotify:track:track-1","name":"Song","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[{"url":"https://img/1.jpg"}]},"duration_ms":1234,"type":"track"},"is_playing":true,"progress_ms":42,"device":{"id":"device-1","name":"Desktop","is_active":true,"is_restricted":false,"is_private_session":false}}`,
			wantNil:    false,
		},
		{
			name:       "returns nil on 204",
			statusCode: http.StatusNoContent,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultClient := http.DefaultClient
			http.DefaultClient = &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						t.Fatalf("method = %s, want GET", req.Method)
					}
					if req.URL.Path != "/v1/me/player/currently-playing" {
						t.Fatalf("path = %s, want /v1/me/player/currently-playing", req.URL.Path)
					}
					if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
						t.Fatalf("Authorization = %q, want bearer token", auth)
					}

					return spotifyTestResponse(req, tt.statusCode, tt.body, nil), nil
				}),
			}
			t.Cleanup(func() {
				http.DefaultClient = defaultClient
			})

			got, err := NewStatic("access-token", nil).GetCurrentlyPlaying(context.Background())
			if err != nil {
				t.Fatalf("GetCurrentlyPlaying() error = %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Fatalf("GetCurrentlyPlaying() = %#v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("GetCurrentlyPlaying() = nil, want track")
			}
			if got.CurrentlyPlayingType != "track" || !got.IsPlaying || got.ProgressMs != 42 {
				t.Fatalf("currently playing = %#v, want parsed fields", got)
			}
			if got.Device.ID != "device-1" || !got.Device.IsActive || got.Device.IsRestricted || got.Device.IsPrivateSession {
				t.Fatalf("device = %#v, want parsed fields", got.Device)
			}
			if got.Item == nil || got.Item.ID != "track-1" || got.Item.ArtistName != "Artist" || got.Item.AlbumName != "Album" || got.Item.ImageURL != "https://img/1.jpg" {
				t.Fatalf("item = %#v, want parsed fields", got.Item)
			}
		})
	}
}
