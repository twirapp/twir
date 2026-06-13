package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
)

type playEvent struct {
	ID       string `json:"id"`
	Duration int    `json:"duration"`
}

func (c *YouTube) handleMessage(session *melody.Session, msg []byte) {
	userId, ok := session.Get("userId")
	if userId == "" || !ok {
		return
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	if data.EventName == "play" {
		parsedData := &playEvent{}
		err = json.Unmarshal(bytes, parsedData)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		c.handlePlay(userId.(string), parsedData)
	}

	if data.EventName == "skip" {
		parsedData := []string{}
		err = json.Unmarshal(bytes, &parsedData)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		c.handleSkip(userId.(string), parsedData)
	}

	if data.EventName == "reorder" {
		var parsedData []model.RequestedSong
		err = json.Unmarshal(bytes, &parsedData)

		c.handleNewOrder(userId.(string), parsedData)
	}

	if data.EventName == "pause" {
		// fmt.Println("get paused")
	}

}

func (c *YouTube) handleSkip(channelId string, ids []string) {
	parsedIDs := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		parsedIDs = append(parsedIDs, parsedID)
	}

	err := c.gorm.
		Model(&model.RequestedSong{}).
		Where(`id IN ? AND "channelId" = ?::uuid`, parsedIDs, channelId).
		Update(`"deletedAt"`, time.Now()).
		Error
	if err != nil {
		c.logger.Error(err.Error())
	}
	redisKey := fmt.Sprintf("songrequests:youtube:%s:currentPlaying", channelId)
	c.redis.Del(context.Background(), redisKey)
}

func (c *YouTube) handleNewOrder(channelId string, songs []model.RequestedSong) {
	var count int64
	err := c.gorm.
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ?::uuid AND "deletedAt" IS NULL`, channelId).
		Count(&count).Error
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	for i, video := range songs {
		err = c.gorm.
			Model(&model.RequestedSong{}).
			Where(`id = ?::uuid`, video.ID).
			Update(`"queuePosition"`, i).
			Error
		if err != nil {
			c.logger.Error(err.Error())
			return
		}
	}
}

func (c *YouTube) handlePlay(userId string, data *playEvent) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("songrequests:youtube:%s:currentPlaying", userId)
	current := c.redis.Get(ctx, redisKey).Val()
	song := &model.RequestedSong{}
	err := c.gorm.Where("id = ?::uuid", data.ID).Find(song).Error
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	if song.ID == "" {
		return
	}

	channelSettings := &model.ChannelSongRequestsSettings{}
	err = c.gorm.Where(`"channel_id" = ?::uuid`, song.ChannelID).Find(channelSettings).Error
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	if channelSettings.ID == "" {
		return
	}

	var songLink string
	if song.SongLink.Valid {
		songLink = song.SongLink.String
	} else {
		songLink = fmt.Sprintf("https://youtu.be/%s", song.VideoID)
	}

	if current == "" && song.ID != "" && channelSettings.AnnouncePlay {
		message := channelSettings.TranslationsNowPlaying
		message = strings.ReplaceAll(message, "{{songTitle}}", song.Title)
		message = strings.ReplaceAll(message, "{{songLink}}", songLink)
		message = strings.ReplaceAll(message, "{{orderedByName}}", song.OrderedByName)
		message = strings.ReplaceAll(
			message,
			"{{orderedByDisplayName}}",
			song.OrderedByDisplayName.String,
		)

		c.bus.Bots.SendMessage.Publish(
			ctx,
			bots.SendMessageRequest{
				ChannelId: song.ChannelID,
				InternalChannelID: func() *uuid.UUID {
					parsedChannelID, err := uuid.Parse(song.ChannelID)
					if err != nil {
						return nil
					}

					return &parsedChannelID
				}(),
				Message:    message,
				IsAnnounce: true,
			},
		)
	}

	c.redis.Set(ctx, redisKey, data.ID, time.Duration(data.Duration)*time.Second)
}
