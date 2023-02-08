package handlers

import (
	"context"
	"fmt"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

func appendMatchedEmotes(emotes []string, channelId string, emotesKeys []string, message []string) []string {
	key := "emotes:global:"
	if channelId != "" {
		key = fmt.Sprintf("emotes:channel:%s:", channelId)
	}

	for _, emoteKey := range emotesKeys {
		emoteSlice := strings.Split(emoteKey, key)
		emote := emoteSlice[1]

		for _, word := range message {
			hasMatch := strings.Contains(word, emote)
			if !hasMatch {
				continue
			}

			emotes = append(emotes, emote)
		}
	}

	return emotes
}

func (c *Handlers) handleEmotes(msg Message) {
	emotes := make([]string, 0)

	for _, emote := range msg.Emotes {
		for i := 0; i < emote.Count; i++ {
			emotes = append(emotes, emote.Name)
		}
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

	emotes = appendMatchedEmotes(emotes, msg.Channel.ID, channelEmotes, splittedMsg)
	emotes = appendMatchedEmotes(emotes, "", globalEmotes, splittedMsg)

	c.db.Transaction(func(tx *gorm.DB) error {
		for _, emote := range emotes {
			err := tx.Create(&model.ChannelEmoteUsage{
				ID:        uuid.NewV4().String(),
				ChannelID: msg.Channel.ID,
				UserID:    msg.User.ID,
				Emote:     emote,
				CreatedAt: time.Now().UTC(),
			}).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}
