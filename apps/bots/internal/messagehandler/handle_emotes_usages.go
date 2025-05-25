package messagehandler

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	channelsemotesusages "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
)

func (c *MessageHandler) handleEmotesUsagesBatched(ctx context.Context, data []handleMessage) {
	var createEmoteUsageInputs []channelsemotesusages.ChannelEmoteUsageInput

	for _, msg := range data {
		for key, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
			for i := 0; i < count; i++ {
				createEmoteUsageInputs = append(
					createEmoteUsageInputs,
					channelsemotesusages.ChannelEmoteUsageInput{
						ID:        uuid.NewString(),
						ChannelID: msg.BroadcasterUserId,
						UserID:    msg.ChatterUserId,
						Emote:     key,
						CreatedAt: time.Now().UTC(),
					},
				)
			}
		}
	}

	err := c.channelsEmotesUsagesRepository.CreateMany(ctx, createEmoteUsageInputs)
	if err != nil {
		c.logger.Error("cannot create emotes usages", slog.Any("err", err))
	}
}

func (c *MessageHandler) handleEmotesUsages(_ context.Context, msg handleMessage) error {
	c.messagesEmotesBatcher.Add(msg)
	return nil
}
