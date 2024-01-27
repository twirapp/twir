package messagehandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/shared"
)

func (c *MessageHandler) handleEmotesUsages(ctx context.Context, msg handleMessage) error {
	emotes := make(map[string]int)

	for _, f := range msg.GetMessage().GetFragments() {
		if f.GetType() != shared.FragmentType_EMOTE {
			continue
		}
		emotes[f.GetText()] += 1
	}

	channelEmotes, err := c.redis.Keys(
		ctx,
		fmt.Sprintf("emotes:channel:%s:*", msg.GetBroadcasterUserId()),
	).Result()
	if err != nil {
		return err
	}

	globalEmotes, err := c.redis.Keys(ctx, "emotes:global:*").Result()
	if err != nil {
		return err
	}

	splittedMsg := strings.Split(msg.GetMessage().GetText(), " ")

	countEmotes(
		emotes,
		channelEmotes,
		splittedMsg,
		fmt.Sprintf("emotes:channel:%s:", msg.GetBroadcasterUserId()),
	)
	countEmotes(emotes, globalEmotes, splittedMsg, "emotes:global:")

	var emotesForCreate []*model.ChannelEmoteUsage

	for key, count := range emotes {
		for i := 0; i < count; i++ {
			emotesForCreate = append(
				emotesForCreate, &model.ChannelEmoteUsage{
					ID:        uuid.NewString(),
					ChannelID: msg.GetBroadcasterUserId(),
					UserID:    msg.GetChatterUserId(),
					Emote:     key,
					CreatedAt: time.Now().UTC(),
				},
			)
		}
	}

	err = c.gorm.WithContext(ctx).CreateInBatches(
		emotesForCreate,
		100,
	).Error

	return err
}

func countEmotes(emotes map[string]int, emotesList []string, splittedMsg []string, key string) {
	for _, e := range emotesList {
		emoteSlice := strings.Split(e, key)
		emote := emoteSlice[1]

		for _, word := range splittedMsg {
			if strings.EqualFold(word, emote) {
				emotes[emote]++
			}
		}
	}
}
