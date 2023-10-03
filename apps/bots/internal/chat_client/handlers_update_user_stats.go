package chat_client

import (
	"log/slog"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"

	uuid "github.com/satori/go.uuid"
)

func (c *ChatClient) updateUserStats(userId, channelId string, badges []string) {
	stream := model.ChannelsStreams{}
	if err := c.services.DB.Where(`"userId" = ?`, channelId).Find(&stream).Error; err != nil {
		c.services.Logger.Error(
			"cannot get channel stream",
			slog.Any("err", err),
			slog.String("channelId", channelId),
		)
		return
	}

	user := model.Users{}
	err := c.services.DB.
		Where(`"id" = ?`, userId).
		Preload("Stats", `"userId" = ? AND "channelId" = ?`, userId, channelId).
		Find(&user).Error
	if err != nil {
		c.services.Logger.Error(
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

		if err := c.services.DB.Create(&user).Error; err != nil {
			c.services.Logger.Error(
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
			err := c.services.DB.Create(newStats).Error
			if err != nil {
				c.services.Logger.Error(
					"cannot create user stats",
					slog.Any("err", err),
					slog.String("channelId", channelId),
					slog.String("channelId", userId),
				)
			}
		} else {
			query := c.services.DB.
				Model(&user.Stats).
				Where(`"userId" = ? AND "channelId" = ?`, userId, channelId).
				Update("is_mod", lo.Contains(badges, "MODERATOR")).
				Update("is_subscriber", lo.Contains(badges, "SUBSCRIBER")).
				Update("is_vip", lo.Contains(badges, "VIP"))

			if stream.ID != "" {
				query.Update("messages", user.Stats.Messages+1)
			}

			err := query.Error
			if err != nil {
				c.services.Logger.Error(
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
