package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func (c *Handlers) handleEmotes(msg *Message) {
	emotes := make(map[string]int)

	for _, emote := range msg.Emotes {
		emotes[emote.Name] = emote.Count
	}

	channelEmotes, err := c.redis.Keys(
		context.Background(),
		fmt.Sprintf("emotes:channel:%s:*", msg.Channel.ID),
	).Result()
	if err != nil {
		return
	}

	globalEmotes, err := c.redis.Keys(context.Background(), "emotes:global:*").Result()
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

	err = c.db.CreateInBatches(
		emotesForCreate,
		100,
	).Error

	if err != nil {
		fmt.Println(err)
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
