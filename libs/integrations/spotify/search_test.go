package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestSpotify_SearchTracks_returns_tracks_and_empty_results(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantCount int
	}{
		{
			name:      "returns tracks",
			body:      `{"tracks":{"items":[{"id":"track-1","uri":"spotify:track:track-1","name":"Song","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[{"url":"https://img/1.jpg"}]},"duration_ms":1234,"type":"track"}]}}`,
			wantCount: 1,
		},
		{
			name:      "returns empty slice",
			body:      `{"tracks":{"items":[]}}`,
			wantCount: 0,
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
					if req.URL.Path != "/v1/search" {
						t.Fatalf("path = %s, want /v1/search", req.URL.Path)
					}
					if req.URL.Query().Get("q") != "lofi beats" {
						t.Fatalf("q = %q, want lofi beats", req.URL.Query().Get("q"))
					}
					if req.URL.Query().Get("type") != "track" {
						t.Fatalf("type = %q, want track", req.URL.Query().Get("type"))
					}
					if req.URL.Query().Get("limit") != "25" {
						t.Fatalf("limit = %q, want 25", req.URL.Query().Get("limit"))
					}
					if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
						t.Fatalf("Authorization = %q, want bearer token", auth)
					}

					return spotifyTestResponse(req, http.StatusOK, tt.body, nil), nil
				}),
			}
			t.Cleanup(func() {
				http.DefaultClient = defaultClient
			})

			tracks, err := NewStatic("access-token", nil).SearchTracks(context.Background(), "lofi beats", 25)
			if err != nil {
				t.Fatalf("SearchTracks() error = %v", err)
			}
			if len(tracks) != tt.wantCount {
				t.Fatalf("len(tracks) = %d, want %d", len(tracks), tt.wantCount)
			}
			if tt.wantCount == 0 {
				return
			}

			got := tracks[0]
			if got.ID != "track-1" || got.URI != "spotify:track:track-1" || got.Name != "Song" {
				t.Fatalf("track = %#v, want parsed fields", got)
			}
			if got.ArtistName != "Artist" || got.AlbumName != "Album" || got.DurationMs != 1234 || got.ImageURL != "https://img/1.jpg" || got.Type != "track" {
				t.Fatalf("track = %#v, want parsed fields", got)
			}
		})
	}
}
