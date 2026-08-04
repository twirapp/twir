package song_requests

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

// AnnounceNowPlaying sends the configured "now playing" message to the channel chat.
// It is a no-op when announcements are disabled in the song requests settings
// or when the nowPlaying translation is empty.
func (s *Service) AnnounceNowPlaying(ctx context.Context, channelID string, song model.RequestedSong) error {
	settings := model.ChannelSongRequestsSettings{}
	if err := s.gorm.WithContext(ctx).
		Where(`"channel_id" = ?`, channelID).
		First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("failed to get song requests settings: %w", err)
	}

	if !settings.AnnouncePlay {
		return nil
	}

	songLink := fmt.Sprintf("https://youtu.be/%s", song.VideoID)
	if song.SongLink.Valid && song.SongLink.String != "" {
		songLink = song.SongLink.String
	}

	orderedByDisplayName := song.OrderedByDisplayName.String
	if orderedByDisplayName == "" {
		orderedByDisplayName = song.OrderedByName
	}

	message := strings.NewReplacer(
		"{{songTitle}}", song.Title,
		"{{songLink}}", songLink,
		"{{songId}}", song.VideoID,
		"{{orderedByName}}", song.OrderedByName,
		"{{orderedByDisplayName}}", orderedByDisplayName,
	).Replace(settings.TranslationsNowPlaying)

	if strings.TrimSpace(message) == "" {
		return nil
	}

	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return fmt.Errorf("failed to parse channel id: %w", err)
	}

	if err := s.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelID:      parsedChannelID,
			Message:        message,
			SkipRateLimits: true,
		},
	); err != nil {
		return fmt.Errorf("failed to publish now playing message: %w", err)
	}

	return nil
}
