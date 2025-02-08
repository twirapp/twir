package messagehandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func (c *MessageHandler) handleEmotesUsages(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	emotes := make(map[string]int)

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE {
			continue
		}
		emotes[f.Text] += 1
	}

	splittedMsg := strings.Fields(msg.Message.Text)

	for _, part := range splittedMsg {
		// do not make redis requests if emote already present in map
		if emote, ok := emotes[part]; ok {
			emotes[part] = emote + 1
			continue
		}

		if exists, _ := c.redis.Exists(
			ctx,
			fmt.Sprintf("emotes:channel:%s:%s", msg.BroadcasterUserId, part),
		).Result(); exists == 1 {
			emotes[part] += 1
			continue
		}

		if exists, _ := c.redis.Exists(
			ctx,
			fmt.Sprintf("emotes:global:%s", part),
		).Result(); exists == 1 {
			emotes[part] += 1
			continue
		}
	}

	var emotesForCreate []model.ChannelEmoteUsage

	for key, count := range emotes {
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
