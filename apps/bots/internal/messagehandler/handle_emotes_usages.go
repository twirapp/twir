package messagehandler

import (
	"context"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleEmotesUsages(ctx context.Context, msg handleMessage) error {
	if msg.EnrichedData.ChannelStream == nil {
		return nil
	}

	var emotesForCreate []model.ChannelEmoteUsage

	for key, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
		for i := 0; i < count; i++ {
			emotesForCreate = append(
				emotesForCreate, model.ChannelEmoteUsage{
					ID:        uuid.NewString(),
					ChannelID: msg.BroadcasterUserId,
					UserID:    msg.ChatterUserId,
					Emote:     key,
					CreatedAt: time.Now().UTC(),
				},
			)
		}
	}

	err := c.gorm.WithContext(ctx).CreateInBatches(
		emotesForCreate,
		100,
	).Error

	return err
}
