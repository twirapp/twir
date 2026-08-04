package song_requests

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/requested_song"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
)

// AnnounceNowPlaying sends the configured "now playing" message to the channel chat.
// It is a no-op when announcements are disabled in the song requests settings
// or when the nowPlaying translation is empty.
func (s *Service) AnnounceNowPlaying(
	ctx context.Context,
	channelID string,
	song requested_song.RequestedSong,
) error {
	settings, err := s.settingsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("failed to get song requests settings: %w", err)
	}

	if !settings.AnnouncePlay {
		return nil
	}

	message := strings.NewReplacer(
		"{{songTitle}}", song.Title,
		"{{songLink}}", song.Link(),
		"{{songId}}", song.VideoID,
		"{{orderedByName}}", song.OrderedByName,
		"{{orderedByDisplayName}}", song.RequesterDisplayName(),
	).Replace(settings.TranslationsNowPlaying)

	if strings.TrimSpace(message) == "" {
		return nil
	}

	if err := s.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelID:      settings.ChannelID,
			Message:        message,
			SkipRateLimits: true,
		},
	); err != nil {
		return fmt.Errorf("failed to publish now playing message: %w", err)
	}

	return nil
}
