package chat_client

import (
	"log/slog"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"

	uuid "github.com/satori/go.uuid"
)

func (c *ChatClient) updateUserStats(
	msg Message,
	badges []string,
) (
	model.Users,
	error,
) {
	user := model.Users{}

	err := c.services.DB.
		Where(`"id" = ?`, msg.User.ID).
		Preload("Stats", `"userId" = ? AND "channelId" = ?`, msg.User.ID, msg.Channel.ID).
		Find(&user).Error
	if err != nil {
		c.services.Logger.Error(
			"cannot find user",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
			slog.String("userId", msg.User.ID),
		)
		return user, err
	}

	// no user found
	if user.ID == "" {
		user.ID = msg.User.ID
		user.ApiKey = uuid.NewV4().String()
		user.IsBotAdmin = false
		user.IsTester = false
		user.Stats = createStats(msg.User.ID, msg.Channel.ID)

		if err := c.services.DB.Create(&user).Error; err != nil {
			c.services.Logger.Error(
				"cannot create user",
				slog.Any("err", err),
				slog.String("channelId", msg.Channel.ID),
				slog.String("channelId", msg.User.ID),
			)
			return user, err
		}
	} else {
		if user.Stats == nil {
			newStats := createStats(msg.User.ID, msg.Channel.ID)
			err := c.services.DB.Create(newStats).Error
			if err != nil {
				c.services.Logger.Error(
					"cannot create user stats",
					slog.Any("err", err),
					slog.String("channelId", msg.Channel.ID),
					slog.String("channelId", msg.User.ID),
				)
			}
			user.Stats = newStats
		} else {
			query := c.services.DB.
				Model(&user.Stats).
				Where(`"userId" = ? AND "channelId" = ?`, msg.DbUser.ID, msg.Channel.ID).
				Update("is_mod", lo.Contains(badges, "MODERATOR")).
				Update("is_subscriber", lo.Contains(badges, "SUBSCRIBER")).
				Update("is_vip", lo.Contains(badges, "VIP"))

			if msg.DbStream.ID != "" {
				query.Update("messages", user.Stats.Messages+1)
			}

			err := query.Error
			if err != nil {
				c.services.Logger.Error(
					"cannot update user",
					slog.Any("err", err),
					slog.String("channelId", msg.Channel.ID),
					slog.String("channelId", msg.User.ID),
				)
			}
		}
	}

	return user, nil
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
