package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func (c *ChatClient) handleEmotes(msg Message) {
	emotes := make(map[string]int)

	for _, emote := range msg.Emotes {
		emotes[emote.Name] = emote.Count
	}

	channelEmotes, err := c.services.Redis.Keys(
		context.Background(),
		fmt.Sprintf("emotes:channel:%s:*", msg.Channel.ID),
	).Result()
	if err != nil {
		c.services.Logger.Error(
			"cannot get emotes",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}

	globalEmotes, err := c.services.Redis.Keys(context.Background(), "emotes:global:*").Result()
	if err != nil {
		return
	}

	splittedMsg := strings.Split(msg.Message, " ")

	countEmotes(emotes, channelEmotes, splittedMsg, fmt.Sprintf("emotes:channel:%s:", msg.Channel.ID))
	countEmotes(emotes, globalEmotes, splittedMsg, "emotes:global:")

	var emotesForCreate []*model.ChannelEmoteUsage

	for key, count := range emotes {
		for i := 0; i < count; i++ {
			emotesForCreate = append(
				emotesForCreate, &model.ChannelEmoteUsage{
					ID:        uuid.NewV4().String(),
					ChannelID: msg.Channel.ID,
					UserID:    msg.User.ID,
					Emote:     key,
					CreatedAt: time.Now().UTC(),
				},
			)
		}
	}

	err = c.services.DB.CreateInBatches(
		emotesForCreate,
		100,
	).Error

	if err != nil {
		c.services.Logger.Error(
			"cannot create emotes",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
	}
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
