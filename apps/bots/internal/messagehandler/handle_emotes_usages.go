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
	// if msg.DbStream == nil {
	// 	return nil
	// }

	emotes := make(map[string]int)

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE {
			continue
		}
		emotes[f.Text] += 1
	}

	channelEmotesIter := c.redis.Scan(
		ctx,
		0,
		fmt.Sprintf("emotes:channel:%s:*", msg.BroadcasterUserId),
		0,
	).Iterator()
	var channelEmotes []string
	for channelEmotesIter.Next(ctx) {
		channelEmotes = append(channelEmotes, channelEmotesIter.Val())
	}

	if err := channelEmotesIter.Err(); err != nil {
		return err
	}

	globalEmotesIter := c.redis.Scan(ctx, 0, "emotes:global:*", 0).Iterator()
	var globalEmotes []string
	for globalEmotesIter.Next(ctx) {
		globalEmotes = append(globalEmotes, globalEmotesIter.Val())
	}

	if err := globalEmotesIter.Err(); err != nil {
		return err
	}

	splittedMsg := strings.Split(msg.Message.Text, " ")

	countEmotes(
		emotes,
		channelEmotes,
		splittedMsg,
		fmt.Sprintf("emotes:channel:%s:", msg.BroadcasterUserId),
	)
	countEmotes(emotes, globalEmotes, splittedMsg, "emotes:global:")

	var emotesForCreate []*model.ChannelEmoteUsage

	for key, count := range emotes {
		for i := 0; i < count; i++ {
			emotesForCreate = append(
				emotesForCreate, &model.ChannelEmoteUsage{
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
