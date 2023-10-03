package pubsub_handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/satont/twir/libs/pubsub"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *handlers) userUpdate(data []byte) {
	userStruct := &pubsub.UserUpdateMessage{}
	if err := json.Unmarshal(data, userStruct); err != nil {
		c.logger.Error("cannot unmarshal incoming data", slog.Any("err", err))
		return
	}

	channel := model.Channels{}
	if err := c.db.Where("id = ?", userStruct.UserID).Find(&channel).Error; err != nil {
		c.logger.Error("cannot find channel", slog.String("userId", userStruct.UserID))
		return
	}

	if channel.ID == "" {
		return
	}

	bot, isBotFound := c.botsService.Instances[channel.BotID]
	if !isBotFound {
		return
	}

	if channel.IsEnabled {
		bot.Join(userStruct.UserName)
	} else {
		bot.Leave(userStruct.UserName)
	}
}
