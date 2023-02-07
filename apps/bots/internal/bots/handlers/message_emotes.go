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

func (c *Handlers) handleEmotes(msg Message) {
	c.db.Transaction(func(tx *gorm.DB) error {
		for _, emote := range msg.Emotes {
			for i := 0; i < emote.Count; i++ {
				tx.Create(&model.ChannelEmoteUsage{
					ID:        uuid.NewV4().String(),
					ChannelID: msg.Channel.ID,
					UserID:    msg.User.ID,
					Emote:     emote.Name,
					CreatedAt: time.Now().UTC(),
				})
			}
		}

		return nil
	})

	emotesKeys, err := c.redis.Keys(
		context.Background(),
		fmt.Sprintf("channels:%s:emotes:*", msg.Channel.ID),
	).Result()

	if err != nil {
		return
	}

	splittedMsg := strings.Split(msg.Message, " ")

	c.db.Transaction(func(tx *gorm.DB) error {
		for _, emoteKey := range emotesKeys {
			emoteSlice := strings.Split(emoteKey, fmt.Sprintf("channels:%s:emotes:", msg.Channel.ID))
			emote := strings.Join(emoteSlice[1:], "")

			if err != nil {
				continue
			}

			for _, word := range splittedMsg {
				hasMatch := strings.Contains(word, emote)

				if !hasMatch {
					continue
				}

				tx.Create(&model.ChannelEmoteUsage{
					ID:        uuid.NewV4().String(),
					ChannelID: msg.Channel.ID,
					UserID:    msg.User.ID,
					Emote:     emote,
					CreatedAt: time.Now().UTC(),
				})
			}
		}

		return nil
	})
}
