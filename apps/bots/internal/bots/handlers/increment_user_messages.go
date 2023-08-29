package handlers

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"

	uuid "github.com/satori/go.uuid"
)

func (c *Handlers) incrementUserMessages(userId, channelId string) {
	stream := model.ChannelsStreams{}
	if err := c.db.Where(`"userId" = ?`, channelId).Find(&stream).Error; err != nil {
		c.logger.Error(
			"cannot get channel stream",
			slog.Any("err", err),
			slog.String("channelId", channelId),
		)
		return
	}

	user := model.Users{}
	err := c.db.
		Where(`"id" = ?`, userId).
		Preload("Stats", `"userId" = ? AND "channelId" = ?`, userId, channelId).
		Find(&user).Error
	if err != nil {
		c.logger.Error(
			"cannot find user",
			slog.Any("err", err),
			slog.String("channelId", channelId),
			slog.String("userId", userId),
		)
		return
	}

	// no user found
	if user.ID == "" {
		user.ID = userId
		user.ApiKey = uuid.NewV4().String()
		user.IsBotAdmin = false
		user.IsTester = false
		user.Stats = createStats(userId, channelId)

		if err := c.db.Create(&user).Error; err != nil {
			c.logger.Error(
				"cannot create user",
				slog.Any("err", err),
				slog.String("channelId", channelId),
				slog.String("channelId", userId),
			)
			return
		}
	} else {
		if user.Stats == nil {
			newStats := createStats(userId, channelId)
			err := c.db.Create(newStats).Error
			if err != nil {
				c.logger.Error(
					"cannot create user stats",
					slog.Any("err", err),
					slog.String("channelId", channelId),
					slog.String("channelId", userId),
				)
			}
		} else if stream.ID != "" {
			err := c.db.
				Model(&user.Stats).
				Where(`"userId" = ? AND "channelId" = ?`, userId, channelId).
				Update("messages", user.Stats.Messages+1).
				Error
			if err != nil {
				c.logger.Error(
					"cannot update user",
					slog.Any("err", err),
					slog.String("channelId", channelId),
					slog.String("channelId", userId),
				)
			}
		}
	}
}

func createStats(userId, channelId string) *model.UsersStats {
	stats := &model.UsersStats{
		ID:                uuid.NewV4().String(),
		UserID:            userId,
		ChannelID:         channelId,
		Messages:          1,
		Watched:           0,
		UsedChannelPoints: 0,
	}
	return stats
}
