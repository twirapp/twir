package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/websockets/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"strings"
	"time"

	"github.com/satont/twir/libs/types/types/api/modules"
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

	data := &types.WebSocketMessage{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
	if data.EventName == "play" {
		parsedData := &playEvent{}
		err = json.Unmarshal(bytes, parsedData)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		c.handlePlay(userId.(string), parsedData)
	}

	if data.EventName == "skip" {
		parsedData := []string{}
		err = json.Unmarshal(bytes, &parsedData)
		if err != nil {
			c.services.Logger.Error(err)
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
		//fmt.Println("get paused")
	}

}

func (c *YouTube) handleSkip(channelId string, ids []string) {
	err := c.services.Gorm.
		Model(&model.RequestedSong{}).
		Where(`id IN (?) AND "channelId" = ?`, ids, channelId).
		Update(`"deletedAt"`, time.Now()).
		Error
	if err != nil {
		c.services.Logger.Error(err)
	}
	redisKey := fmt.Sprintf("songrequests:youtube:%s:currentPlaying", channelId)
	c.services.Redis.Del(context.Background(), redisKey)
}

func (c *YouTube) handleNewOrder(channelId string, songs []model.RequestedSong) {
	var count int64
	err := c.services.Gorm.
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
		Count(&count).Error
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	for i, video := range songs {
		err = c.services.Gorm.
			Model(&model.RequestedSong{}).
			Where(`id = ?`, video.ID).
			Update(`"queuePosition"`, i).
			Error
		if err != nil {
			c.services.Logger.Error(err)
			return
		}
	}
}

func (c *YouTube) handlePlay(userId string, data *playEvent) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("songrequests:youtube:%s:currentPlaying", userId)
	current := c.services.Redis.Get(ctx, redisKey).Val()
	song := &model.RequestedSong{}
	err := c.services.Gorm.Where("id = ?", data.ID).Find(song).Error
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
	if song.ID == "" {
		return
	}

	channelSettings := &model.ChannelModulesSettings{}
	err = c.services.Gorm.Where(
		`"channelId" = ? AND type = ?`, song.ChannelID, "youtube_song_requests",
	).Find(channelSettings).Error
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
	if channelSettings.ID == "" {
		return
	}

	youtubeSettings := &modules.YouTubeSettings{}
	err = json.Unmarshal(channelSettings.Settings, youtubeSettings)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	var songLink string
	if song.SongLink.Valid {
		songLink = song.SongLink.String
	} else {
		songLink = fmt.Sprintf("https://youtu.be/%s", song.VideoID)
	}

	if current == "" && song.ID != "" && youtubeSettings.AnnouncePlay != nil && *youtubeSettings.AnnouncePlay {
		message := youtubeSettings.Translations.NowPlaying
		message = strings.ReplaceAll(message, "{{songTitle}}", song.Title)
		message = strings.ReplaceAll(message, "{{songLink}}", songLink)
		message = strings.ReplaceAll(message, "{{orderedByName}}", song.OrderedByName)
		message = strings.ReplaceAll(message, "{{orderedByDisplayName}}", song.OrderedByDisplayName.String)

		c.services.Grpc.Bots.SendMessage(
			ctx, &bots.SendMessageRequest{
				ChannelId:   song.ChannelID,
				ChannelName: nil,
				Message:     message,
				IsAnnounce:  lo.ToPtr(true),
			},
		)
	}

	c.services.Redis.Set(ctx, redisKey, data.ID, time.Duration(data.Duration)*time.Second)
}
