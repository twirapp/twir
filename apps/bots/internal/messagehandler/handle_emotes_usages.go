package messagehandler

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleEmotesUsagesBatched(ctx context.Context, data []handleMessage) {
	var emotesForCreate []model.ChannelEmoteUsage

	for _, msg := range data {
		for key, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
			for i := 0; i < count; i++ {
				emotesForCreate = append(
					emotesForCreate,
					model.ChannelEmoteUsage{
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

	err := c.gorm.WithContext(ctx).CreateInBatches(
		emotesForCreate,
		100,
	).Error
	if err != nil {
		c.logger.Error("cannot create emotes usages", slog.Any("err", err))
	}
}

func (c *MessageHandler) handleEmotesUsages(ctx context.Context, msg handleMessage) error {
	c.messagesEmotesBatcher.Add(msg)

	return nil
}
