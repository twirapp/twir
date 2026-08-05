package song_request

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscorespotify "github.com/twirapp/twir/libs/bus-core/spotify"
	songrequestmode "github.com/twirapp/twir/libs/entities/song_request_mode"
	songrequestsettingsentity "github.com/twirapp/twir/libs/entities/song_requests_settings"
)

func TestProcessFromDonationSpotifyModePublishesSpotifyCreateRequest(t *testing.T) {
	channelID := uuid.New()
	publishedRequests := &songRequestQueue[buscorespotify.CreateSongRequestRequest, buscorespotify.CreateSongRequestResponse]{
		requestResponse: &buscore.QueueResponse[buscorespotify.CreateSongRequestResponse]{
			Data: buscorespotify.CreateSongRequestResponse{
				Request: buscorespotify.SongRequestData{
					Title:  "track title",
					Artist: "track artist",
				},
			},
		},
	}
	bus := buscore.NewNatsBus(nil)
	bus.Spotify.CreateSongRequest = publishedRequests

	service := &SongRequest{
		gorm:    newSongRequestTestDB(t),
		twirBus: bus,
		logger:  slog.Default(),
		songRequestsSettingsRepo: &songRequestSettingsRepositoryFake{
			settings: songrequestsettingsentity.Settings{
				Enabled:                     true,
				Mode:                        songrequestmode.ModeSpotify,
				TakeSongFromDonationMessage: true,
			},
		},
		channelService: nil,
	}

	err := service.ProcessFromDonation(
		context.Background(),
		ProcessFromDonationInput{
			Text:      "https://open.spotify.com/track/example",
			Username:  "donor-user",
			ChannelID: channelID.String(),
		},
	)
	if err != nil {
		t.Fatalf("process spotify donation song request: %v", err)
	}
	if len(publishedRequests.requested) != 1 {
		t.Fatalf("requested requests = %d, want 1", len(publishedRequests.requested))
	}

	request := publishedRequests.requested[0]
	if request.ChannelID != channelID.String() {
		t.Fatalf("channel ID = %q, want %q", request.ChannelID, channelID.String())
	}
	if request.RequesterName != "donor-user" {
		t.Fatalf("requester name = %q, want %q", request.RequesterName, "donor-user")
	}
	if request.RequesterDisplayName != "donor-user" {
		t.Fatalf("requester display name = %q, want %q", request.RequesterDisplayName, "donor-user")
	}
	if request.Source != "donation" {
		t.Fatalf("source = %q, want %q", request.Source, "donation")
	}
	if request.Query != "https://open.spotify.com/track/example" {
		t.Fatalf("query = %q, want %q", request.Query, "https://open.spotify.com/track/example")
	}
}
